package model

type User struct {
	ID        string  `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Email     string  `json:"email" db:"email"`
	Age       int     `json:"age" db:"age"`
	CreatedAt *string `json:"created_at" db:"created_at"`
	UpdatedAt *string `json:"updated_at" db:"updated_at"`
}

var UserFields = []string{"id", "name", "email", "age", "created_at", "updated_at"}
