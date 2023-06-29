package cupcake

import (
	"fmt"
	"net/http"
)

// cupcake request handler
type HandlerFunc func(http.ResponseWriter, *http.Request)

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
	pattern := r.URL.Path
	method := r.Method

	if handlers, ok := cc.router[pattern]; ok {
		if handler, ok := handlers[method]; ok {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, "405 Method Not Allowed")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
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
