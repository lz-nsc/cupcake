package cupcake

import (
	"net/http"
)

// cupcake request handler
type HandlerFunc func(*Response, *Request)

type Cupcake struct {
	*RouteGroup
	router *router
	groups []*RouteGroup
}

// Construct a new cupcake server
func New() *Cupcake {
	engine := &Cupcake{router: newRouter()}
	// Make the engine itself a group with empty prefix
	engine.RouteGroup = NewGroup("", engine)
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
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
