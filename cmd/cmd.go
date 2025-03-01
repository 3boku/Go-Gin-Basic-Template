package cmd

import (
	"Go-Gin-Basic-Template/database"
	"Go-Gin-Basic-Template/router"
)

type Cmd struct {
	router *router.Router
}

func NewCmd() {
	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}

	c := &Cmd{
		router: router.NewRouter(db),
	}

	c.router.SetupRoutes()
	err = c.router.ServerStart()
	if err != nil {
		panic(err)
	}
}
