package main

import (
	"cupcake"
	"fmt"
	"net/http"

	"github.com/lz-nsc/cupcake"
)

func main() {
	cc := cupcake.New()
	cc.GET("/cupcake", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to cupcake!")
	})
	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
