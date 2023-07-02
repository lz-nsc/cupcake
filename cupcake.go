package cupcake

import (
	"net/http"
)

// cupcake request handler
type HandlerFunc func(*Response, *Request)

type Cupcake struct {
	router *router
}

// Construct a new cupcake server
func New() *Cupcake {
	return &Cupcake{router: newRouter()}
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
	cc.router.handle(NewResponse(w), NewRequest(r))
}

func (cc *Cupcake) run(address string) {
	http.ListenAndServe(address, cc)
}

func (cc *Cupcake) GET(pattern string, handler HandlerFunc) {
	cc.router.addRouter(GET, pattern, handler)
}

func (cc *Cupcake) POST(pattern string, handler HandlerFunc) {
	cc.router.addRouter(POST, pattern, handler)
}
func (cc *Cupcake) PUT(pattern string, handler HandlerFunc) {
	cc.router.addRouter(PUT, pattern, handler)
}
func (cc *Cupcake) DELETE(pattern string, handler HandlerFunc) {
	cc.router.addRouter(DELETE, pattern, handler)
}
