package cupcake

import (
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	req      *http.Request
	path     string
	method   string
	params   map[string]string
	data     map[string]string
	fullData map[string]string
	wild     string
}

func NewRequest(r *http.Request) *Request {
	r.ParseForm()
	req := &Request{
		req:      r,
		path:     r.URL.Path,
		method:   r.Method,
		params:   make(map[string]string),
		data:     make(map[string]string),
		fullData: make(map[string]string),
	}
	for key, value := range r.PostForm {
		if len(value) == 0 {
			req.data[key] = ""
		} else {
			req.data[key] = value[0]
		}
	}
	for key, value := range r.Form {
		if len(value) == 0 {
			req.fullData[key] = ""
		} else {
			req.fullData[key] = value[0]
		}
	}
	return req
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
func (r *Request) SetParams(params map[string]string) {
	r.params = params
}
func (r *Request) SetWild(wild string) {
	r.wild = wild
}

func (r Request) Params() map[string]string {
	return r.params
}
func (r Request) Param(key string) string {
	return r.params[key]
}

func (r Request) Wild() string {
	return r.wild
}

func (r Request) Body() io.ReadCloser {
	return r.req.Body
}

func (r Request) Data() map[string]string {
	return r.data
}
