package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:generate daogen -type User
type User struct {
	ID   int64
	Name string
}

func (u *User) Insert(db *sqlx.DB) (int64, error) {
	err := db.QueryRow(`INSERT INTO users (name) VALUES ( $1 ) RETURNING id`, u.Name).Scan(&u.ID)
	return u.ID, err
}
