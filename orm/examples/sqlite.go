package main

import (
	"fmt"

	"github.com/lz-nsc/cupcake/orm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `cupcakeorm:"PRIMARY KEY"`
	Age  int
}
type Users []User

func main() {
	oe, err := orm.NewORMEngine("sqlite3", "cupcake.db")
	if err != nil {
		fmt.Printf("failed to create orm engine, err: %s\n", err.Error())
		return
	}

	defer oe.Close()

	session := oe.NewSession()

	session.Model(User{})
	_ = session.DropTable()

	_ = session.CreateTable()
	// Insert
	count, err := session.Insert(
		&User{
			Name: "John",
			Age:  20,
		},
		&User{
			Name: "Jack",
			Age:  30,
		},
	)
	if err != nil {
		fmt.Printf("failed to insert record, err: %s\n", err.Error())
		return
	}
	fmt.Printf("Successfully insert %d row(s)\n", count)

	// List
	users := &Users{}
	err = session.FindAll(users)
	if err != nil {
		fmt.Printf("failed to list records, err: %s\n", err.Error())
		return
	}

	// First and OrderBy
	user := &User{}
	err = session.OrderBy("age desc").FindOne(user)
	if err != nil {
		fmt.Printf("failed to get first record, err: %s\n", err.Error())
		return
	}

	fmt.Printf("Successfully get first record, name: %s\n", user.Name)

	// Update
	fmt.Printf("successfully list %d row(s)\n", len(*users))
	count, err = session.Where("Name = ?", "John").Update("Age", 25)
	if err != nil {
		fmt.Printf("failed to update record, err: %s\n", err.Error())
		return
	}

	fmt.Printf("Successfully update %d row(s)\n", count)
	// Delete
	count, err = session.Where("Name = ?", "Jack").Delete()
	if err != nil {
		fmt.Printf("failed to delete record, err: %s\n", err.Error())
		return
	}
	fmt.Printf("Successfully delete %d row(s)\n", count)
}
