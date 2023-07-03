package main

import (
	"fmt"
	"net/http"

	"github.com/lz-nsc/cupcake"
)

func main() {
	cc := cupcake.New()
	cc.GET("/cupcake", func(resp *cupcake.Response, req *cupcake.Request) {
		resp.String(http.StatusOK, "Welcome to cupcake!")
	})

	group := cc.Group("/v1")
	group.GET("/cupcake/{name:[a-z]+}", func(resp *cupcake.Response, req *cupcake.Request) {
		name := req.Param("name")
		resp.String(http.StatusOK, fmt.Sprintf("Welcome to cupcake name[%s]!", name))
	})

	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
