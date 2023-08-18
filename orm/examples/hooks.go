package main

import (
	"fmt"

	"github.com/lz-nsc/cupcake/orm"
	"github.com/lz-nsc/cupcake/orm/session"

	_ "github.com/mattn/go-sqlite3"
)

type HookUser struct {
	ID   int `cupcakeorm:"PRIMARY KEY"`
	Name string
	Age  int
}

func (user *HookUser) BeforeInsert(session *session.Session) error {
	fmt.Println("Triggered BeforeInsert hook")
	user.ID += 100
	return nil
}

func main() {
	oe, err := orm.NewORMEngine("sqlite3", "cupcake.db")
	if err != nil {
		fmt.Printf("failed to create orm engine, err: %s\n", err.Error())
		return
	}

	defer oe.Close()

	session := oe.NewSession()

	session.Model(HookUser{})
	_ = session.DropTable()

	_ = session.CreateTable()
	// Insert
	count, err := session.Insert(
		&HookUser{
			ID:   1,
			Name: "John",
			Age:  20,
		},
	)

	if err != nil {
		fmt.Printf("failed to insert record, err: %s\n", err.Error())
		return
	}
	fmt.Printf("Successfully insert %d row(s)\n", count)
	user := &HookUser{}
	err = session.FindOne(user)
	if err != nil {
		fmt.Printf("failed to get first record, err: %s\n", err.Error())
		return
	}
	fmt.Printf("Successfully get first record, ID: %d\n", user.ID)

}
