package main

import (
	"Go-Gin-Basic-Template/cmd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./secret/.env")
	if err != nil {
		panic(err)
	}

	cmd.NewCmd()
}
