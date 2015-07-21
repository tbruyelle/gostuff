package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

type User struct {
	ID   int64
	Name string
}

func (u *User) Insert(db *sqlx.DB) (int64, error) {
	err := db.QueryRow(`INSERT INTO users (name) VALUES ( $1 ) RETURNING id`, u.Name).Scan(&u.ID)
	return u.ID, err
}

func CreateUserTable(db *sqlx.DB) error {
	sql := `create table users (
	id bigserial primary key,
	name text not null)`
	_, err := db.Exec(sql)
	return err
}

func FindUserById(db *sqlx.DB, ID int64) (*User, error) {
	u := &User{}
	err := db.Get(u, "select * from users where id=$1", ID)
	return u, err
}

func main() {
	// DB connect
	dsn := "postgres://tbruyelle:toto@localhost:5432/daogen?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting db", err)
		os.Exit(2)
	}
	// Create the table
	err = CreateUserTable(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error create table:", err)
	}

	// Add a user
	u := &User{Name: "Nathan"}
	ID, err := u.Insert(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error insert", err)
		os.Exit(2)
	}

	// Find that user
	u, err = FindUserById(db, ID)
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
