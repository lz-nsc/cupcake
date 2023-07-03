package cupcake

import (
	"log"
	"time"
)

func Logger(handler HandlerFunc) HandlerFunc {
	return HandlerFunc(func(resp *Response, req *Request) {
		t := time.Now()
		handler(resp, req)
		log.Printf("%s %s response with %d in %v", req.method, req.path, resp.statusCode, time.Since(t))
	})
}
