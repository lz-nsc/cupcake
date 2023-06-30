package cupcake

import (
	"fmt"
	"net/http"
)

// cupcake request handler
type HandlerFunc func(*Response, *Request)

type handlers map[string]HandlerFunc

type Cupcake struct {
	router map[string]handlers
}

// Construct a new cupcake server
func New() *Cupcake {
	return &Cupcake{router: make(map[string]handlers)}
}

// Run a cupcake server
// cupcake.Run()
// cupcake.Run(":80")
// cupcake.Run("localhost")
// cupcake.Run("127.0.0.1:8080")
func (cc *Cupcake) Run(params ...string) {
	if len(params) > 0 {
		// Only accept first param and discard the rest
		cc.run(params[0])
	}
	cc.run("")
}

func (cc *Cupcake) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := NewRequest(r)
	response := NewResponse(w)

	if handlers, ok := cc.router[request.Path()]; ok {
		if handler, ok := handlers[request.Method()]; ok {
			fmt.Println(request.String())

			handler(response, request)
		} else {
			fmt.Printf("%s : 405 Method Not Allowed\n", request.String())
			response.Error(http.StatusMethodNotAllowed, "405 Method Not Allowed")
		}
	} else {
		fmt.Printf("%s : 404 NOT FOUND\n", request.String())
		response.Error(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", request.Path()))
	}
}

func (cc *Cupcake) run(address string) {
	http.ListenAndServe(address, cc)
}

func (cc *Cupcake) addRouter(method string, pattern string, handler HandlerFunc) {
	if cc.router[pattern] == nil {
		cc.router[pattern] = handlers{}
	}
	cc.router[pattern][method] = handler
}

func (cc *Cupcake) GET(pattern string, handler HandlerFunc) {
	cc.addRouter(http.MethodGet, pattern, handler)
}

func (cc *Cupcake) POST(pattern string, handler HandlerFunc) {
	cc.addRouter(http.MethodPost, pattern, handler)
}
func (cc *Cupcake) PUT(pattern string, handler HandlerFunc) {
	cc.addRouter(http.MethodPut, pattern, handler)
}
func (cc *Cupcake) DELETE(pattern string, handler HandlerFunc) {
	cc.addRouter(http.MethodDelete, pattern, handler)
}
