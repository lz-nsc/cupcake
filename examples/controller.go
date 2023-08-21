package main

import (
	"github.com/lz-nsc/cupcake"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `json:"name" cupcakeorm:"PRIMARY KEY"`
	Age  int    `json:"age"`
}
type UserController struct {
	*cupcake.BaseController
}

func main() {
	cc := cupcake.New()

	controller := UserController{cupcake.NewBaseController(&User{})}
	cc.Route("/users", controller)
	cc.Run(":8080")
}
