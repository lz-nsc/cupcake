package main

import (
	"fmt"

	"github.com/lz-nsc/cupcake"
)

func main() {
	cc := cupcake.New()
	cc.Static("assets", "statics")
	fmt.Println("Start cupcake server")
	cc.Run(":8080")
}
