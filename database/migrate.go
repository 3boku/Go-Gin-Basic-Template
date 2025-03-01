package database

import (
	"Go-Gin-Basic-Template/types"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(
		&types.Product{},
	)
}
