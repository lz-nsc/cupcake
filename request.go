package cupcake

import (
	"fmt"
	"net/http"
)

type Request struct {
	req    *http.Request
	path   string
	method string
	params map[string]string
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		req:    r,
		path:   r.URL.Path,
		method: r.Method,
		params: make(map[string]string),
	}
}
func (r *Request) PostForm(key string) string {
	return r.req.FormValue(key)
}

func (r *Request) Query(key string) string {
	return r.req.URL.Query().Get(key)
}

func (r Request) String() string {
	return fmt.Sprintf("Request: %s %s", r.method, r.path)
}

func (r Request) Path() string {
	return r.path
}

func (r Request) Method() string {
	return r.method
}
func (r *Request) SetParam(params map[string]string) {
	r.params = params
}
func (r Request) Param(key string) string {
	return r.params[key]
}
