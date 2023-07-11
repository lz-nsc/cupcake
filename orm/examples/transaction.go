package main

import (
	"fmt"

	"github.com/lz-nsc/cupcake/orm"
	"github.com/lz-nsc/cupcake/orm/session"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `cupcakeorm:"PRIMARY KEY"`
	Age  int
}
type Users []User

func main() {
	oe, err := orm.NewORMEngine("sqlite3", "gee.db")
	if err != nil {
		fmt.Printf("failed to create orm engine, err: %s\n", err.Error())
		return
	}

	defer oe.Close()
	s := oe.NewSession()

	s.Model(User{})
	_ = s.DropTable()
	_ = s.CreateTable()

	_, err = oe.Transaction(func(s *session.Session) (result interface{}, err error) {
		s.Model(&User{})
		_, err = s.Insert(&User{"Tom", 18})
		_, err = s.Insert(&User{"Jack", 28})
		return
	})

	users := &Users{}
	_ = s.FindAll(users)
	for idx, user := range *users {
		fmt.Printf("User#%d - Name:%s, Age:%d\n", idx, user.Name, user.Age)
	}
}
