package router

import (
	"Go-Gin-Basic-Template/controller"
	"Go-Gin-Basic-Template/httpHandler"
	"Go-Gin-Basic-Template/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

type Router struct {
	Engine *gin.Engine

	ProductHandler *httpHandler.ProductHandler
}

func NewRouter(db *gorm.DB) *Router {
	productRepository := &repository.ProductRepository{DB: db}
	productController := &controller.ProductController{ProductRepository: productRepository}
	productHandler := &httpHandler.ProductHandler{ProductController: productController}

	r := &Router{
		Engine:         gin.Default(),
		ProductHandler: productHandler,
	}

	return r
}

func (r *Router) ServerStart() error {
	return r.Engine.Run(os.Getenv("PORT"))
}

func (r *Router) SetupRoutes() {
	product := r.Engine.Group("/product")
	{
		product.POST("", r.ProductHandler.Insert)
		product.PATCH("/:id", r.ProductHandler.Update)
		product.DELETE("/:id", r.ProductHandler.Delete)
		product.GET("", r.ProductHandler.GetAll)
		product.GET("/:id", r.ProductHandler.GetByID)
	}
}
