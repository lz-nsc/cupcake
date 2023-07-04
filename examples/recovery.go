package main

import (
	"fmt"

	"github.com/lz-nsc/cupcake"
	"github.com/lz-nsc/cupcake/middlewares"
)

func main() {
	cc := cupcake.New()

	cc.MiddlerWare(middlewares.Recovery)

	cc.GET("/cupcake", func(resp *cupcake.Response, req *cupcake.Request) {
		panic("For no reason. LOL")
	})

	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
