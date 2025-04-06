package model

type Product struct {
	ID    int64   `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
	Stock int     `json:"stock" db:"stock"`
}

var ProductFields = []string{"id", "name", "price", "stock"}
