package cupcake

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/lz-nsc/cupcake/orm"
	"github.com/lz-nsc/cupcake/orm/session"
)

// cupcake request handler
type HandlerFunc func(*Response, *Request)

var defaultDBEngine *orm.ORMEngine

type Cupcake struct {
	*RouteGroup
	router *router
	groups []*RouteGroup
	render *template.Template // for html render
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
	cc.handle(NewResponse(w, cc.render), NewRequest(r))
}

func (cc *Cupcake) run(address string) {
	http.ListenAndServe(address, cc)
}

func (cc *Cupcake) LoadTemplates(path string) {
	cc.render = template.Must(template.New("").ParseGlob(path))
}

func newDBSession() *session.Session {
	if defaultDBEngine == nil {
		oe, err := orm.NewORMEngine("sqlite3", "cupcake.db")
		if err != nil {
			panic(fmt.Sprintf("Failed to create db engine, err: %s", err))
		}
		defaultDBEngine = oe
	}
	return defaultDBEngine.NewSession()
}
