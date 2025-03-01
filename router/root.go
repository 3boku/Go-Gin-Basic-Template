package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter(db *gorm.DB) *Router {
	r := &Router{
		Engine: gin.Default(),
	}

	return r
}

func (r *Router) ServerStart() error {
	return r.Engine.Run(os.Getenv("PORT"))
}

func (r *Router) SetupRoutes() {
	// TODO: Implement routes

}
