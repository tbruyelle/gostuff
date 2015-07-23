package main

import "fmt"
import "github.com/jmoiron/sqlx"
import "github.com/tbruyelle/gostuff/dalgen/example/model"
import "os"

func main() {
	// DB connect
	dsn := "postgres://tbruyelle:toto@localhost:5432/daogen?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting db", err)
		os.Exit(2)
	}
	// Create the table
	err = model.CreateUserTable(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error create table:", err)
	}

	// Add a user
	u := &model.User{Name: "Nathan"}
	id, err := u.Insert(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error insert", err)
		os.Exit(2)
	}

	// Find that user
	u, err = model.FindUserByID(db, id)
	if err != nil {
		fmt.Println(os.Stderr, "Error requesting", err)
		os.Exit(2)
	}
	if u == nil {
		fmt.Println("User not found")
	} else {
		fmt.Println("User found", *u)
	}
}
