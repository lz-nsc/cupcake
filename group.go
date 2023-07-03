package cupcake

type RouteGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Cupcake
}

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

func (group *RouteGroup) addRouter(method methodType, path string, handler HandlerFunc) {
	if path[0] != '/' {
		path = "/" + path
	}
	pattern := group.prefix + path
	group.engine.router.addRouter(method, pattern, handler)
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
