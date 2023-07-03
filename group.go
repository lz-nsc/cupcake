package cupcake

import (
	"fmt"
	"net/http"
)

type RouteGroup struct {
	prefix      string
	middlewares []MiddlerWare
	engine      *Cupcake
}

type MiddlerWare func(HandlerFunc) HandlerFunc

func NewGroup(prefix string, engine *Cupcake) *RouteGroup {
	// Remove '/' from tail
	if prefix != "" && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	return &RouteGroup{
		prefix: prefix,
		engine: engine,
	}
}
func (group *RouteGroup) Group(prefix string) *RouteGroup {
	engine := group.engine
	child := &RouteGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}

	engine.groups = append(engine.groups, child)
	return child
}

func (group *RouteGroup) MiddlerWare(m MiddlerWare) {
	group.middlewares = append(group.middlewares, m)
}

func (group *RouteGroup) addRouter(method methodType, path string, handler HandlerFunc) {
	if path[0] != '/' {
		path = "/" + path
	}
	pattern := group.prefix + path
	handler = group.wrapMiddlewares(handler)
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouteGroup) handle(resp *Response, req *Request) {
	handler, err := group.engine.router.handler(resp, req)
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

	fmt.Println(req.String())
	handler(resp, req)
}

func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter(GET, pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter(POST, pattern, handler)
}
func (group *RouteGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRouter(PUT, pattern, handler)
}
func (group *RouteGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRouter(DELETE, pattern, handler)
}
func (group *RouteGroup) wrapMiddlewares(handler HandlerFunc) HandlerFunc {
	for _, middlerWare := range group.middlewares {
		handler = middlerWare(handler)
	}
	return handler
}
