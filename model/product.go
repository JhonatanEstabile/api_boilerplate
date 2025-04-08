package model

type Product struct {
	ID        string  `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Price     float64 `json:"price" db:"price"`
	Stock     int     `json:"stock" db:"stock"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}

var ProductFields = []string{"id", "name", "price", "stock", "created_at", "updated_at"}
