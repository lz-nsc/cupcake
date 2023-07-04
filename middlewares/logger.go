package middlewares

import (
	"log"
	"time"

	"github.com/lz-nsc/cupcake"
)

func Logger(handler cupcake.HandlerFunc) cupcake.HandlerFunc {
	return cupcake.HandlerFunc(func(resp *cupcake.Response, req *cupcake.Request) {
		t := time.Now()
		handler(resp, req)
		log.Printf("%s %s response with %d in %v", req.Method(), req.Path(), resp.StatusCode(), time.Since(t))
	})
}
