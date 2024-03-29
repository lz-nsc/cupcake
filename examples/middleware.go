package main

import (
	"fmt"
	"net/http"

	"github.com/lz-nsc/cupcake"
	"github.com/lz-nsc/cupcake/middlewares"
)

func main() {
	cc := cupcake.New()
	cc.MiddlerWare(middlewares.Logger)
	cc.GET("/cupcake", func(resp *cupcake.Response, req *cupcake.Request) {
		resp.String(http.StatusOK, "Welcome to cupcake!")
	})

	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
