package cupcake

import (
	"fmt"
	"net/http"
)

type router struct {
	node *radixNode
}

func newRouter() *router {
	return &router{node: NewNode("")}
}

func (r *router) addRouter(method methodType, path string, handler HandlerFunc) {
	r.node.InsertNode(path, method, handler)
}

func (r *router) handle(resp *Response, req *Request) {
	method := parseMethod(req.Method())
	handler, params, err := r.node.Route(req.Path(), method)
	if err != nil {
		switch err {
		case ErrNotAllow:
			fmt.Printf("%s : 405 Method Not Allowed\n", req.String())
			resp.Error(http.StatusMethodNotAllowed, "405 Method Not Allowed")
		case ErrNotFound:
			fmt.Printf("%s : 404 NOT FOUND\n", req.String())
			resp.Error(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", req.Path()))
		default:
			fmt.Printf("Unknown error: %s\n", err.Error())
			resp.Error(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s\n", err.Error()))
		}
		return
	}

	req.SetParam(params)
	fmt.Println(req.String())
	handler(resp, req)
}
