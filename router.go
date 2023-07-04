package cupcake

type router struct {
	node *radixNode
}

func newRouter() *router {
	return &router{node: NewNode("")}
}

func (r *router) addRouter(method methodType, path string, handler HandlerFunc) {
	r.node.InsertNode(path, method, handler)
}

func (r *router) handler(resp *Response, req *Request) (HandlerFunc, error) {
	method := parseMethod(req.Method())
	handler, params, wild, err := r.node.Route(req.Path(), method)
	if err != nil {
		return nil, err
	}
	req.SetParams(params)
	req.SetWild(wild)

	return handler, nil
}
