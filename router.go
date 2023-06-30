package cupcake

import (
	"fmt"
	"net/http"
)

type handlers map[string]HandlerFunc

type router struct {
	pathes map[string]handlers
}

func newRouter() *router {
	return &router{pathes: make(map[string]handlers)}
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	if r.pathes[pattern] == nil {
		r.pathes[pattern] = handlers{}
	}
	r.pathes[pattern][method] = handler
}
func (r *router) handle(resp *Response, req *Request) {
	if handlers, ok := r.pathes[req.Path()]; ok {
		if handler, ok := handlers[req.Method()]; ok {
			fmt.Println(req.String())
			handler(resp, req)
		} else {
			fmt.Printf("%s : 405 Method Not Allowed\n", req.String())
			resp.Error(http.StatusMethodNotAllowed, "405 Method Not Allowed")
		}
	} else {
		fmt.Printf("%s : 404 NOT FOUND\n", req.String())
		resp.Error(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", req.Path()))
	}
}
