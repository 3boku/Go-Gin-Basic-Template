package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDatabase() (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASS"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PORT"))), &gorm.Config{})
	return db, err
}
