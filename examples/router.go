package main

import (
	"fmt"
	"net/http"

	"github.com/lz-nsc/cupcake"
)

func main() {
	cc := cupcake.New()
	cc.GET("/cupcake/{id}", func(resp *cupcake.Response, req *cupcake.Request) {
		id := req.Param("id")
		resp.String(http.StatusOK, fmt.Sprintf("Welcome to cupcake id[%s]!", id))
	})
	cc.GET("/cupcake/{name:[a-z]+}", func(resp *cupcake.Response, req *cupcake.Request) {
		name := req.Param("name")
		resp.String(http.StatusOK, fmt.Sprintf("Welcome to cupcake name[%s]!", name))
	})

	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
