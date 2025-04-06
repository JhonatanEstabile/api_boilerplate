package model

type User struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Age   int    `json:"age" db:"age"`
}

var UserFields = []string{"id", "name", "email", "age"}
