package types

type Product struct {
	BasicModel
	Name  string  `gorm:"name"`
	Price float64 `gorm:"price"`
}
