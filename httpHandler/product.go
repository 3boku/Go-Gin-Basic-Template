package httpHandler

import (
	"Go-Gin-Basic-Template/controller"
	"Go-Gin-Basic-Template/types"
	"Go-Gin-Basic-Template/types/requestTypes"
	"Go-Gin-Basic-Template/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductHandler struct {
	ProductController *controller.ProductController
}

func (h *ProductHandler) Insert(c *gin.Context) {
	var product requestTypes.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	statusCode, message, err := h.ProductController.Insert(&product)
	if err != nil {
		utils.RespondWithError(c, statusCode, message, err)
		return
	}

	utils.RespondWithSuccess(c, statusCode, message)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var product requestTypes.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	statusCode, message, err := h.ProductController.Update(id, &product)
	if err != nil {
		utils.RespondWithError(c, statusCode, message, err)
		return
	}

	utils.RespondWithSuccess(c, statusCode, message)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	statusCode, message, err := h.ProductController.Delete(id)
	if err != nil {
		utils.RespondWithError(c, statusCode, message, err)
		return
	}

	utils.RespondWithSuccess(c, statusCode, message)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	statusCode, product, err := h.ProductController.GetAll()
	if err != nil {
		utils.RespondWithError(c, statusCode, "SELECT 오류", err)
		return
	}

	utils.RespondWithGet(c, statusCode, *product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	statusCode, product, err := h.ProductController.Get(id)
	if err != nil {
		utils.RespondWithError(c, statusCode, "SELECT 오류", err)
		return
	}
	var response []types.Product
	response = append(response, *product)

	utils.RespondWithGet(c, statusCode, response)
}
