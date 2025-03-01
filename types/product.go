package types

import "github.com/google/uuid"

type Product struct {
	BasicModel
	ID    uuid.UUID `gorm:"id"`
	Name  string    `gorm:"name"`
	Price float64   `gorm:"price"`
}
