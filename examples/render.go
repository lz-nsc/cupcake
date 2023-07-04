package main

import (
	"fmt"
	"net/http"

	"github.com/lz-nsc/cupcake"
)

func main() {
	cc := cupcake.New()

	cc.LoadTemplates("templates/*")
	cc.GET("/cupcake/{name:[a-z]*}", func(resp *cupcake.Response, req *cupcake.Request) {
		resp.Render(http.StatusOK, "welcome.tmpl", req.Params())
	})

	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
