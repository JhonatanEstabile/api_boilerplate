package model

type Store struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

var StoreFields = []string{"id", "name", "description"}
