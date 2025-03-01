package utils

import (
	"Go-Gin-Basic-Template/types"
	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, status int, message string, err error) {
	response := &ErrorResponse{
		Status: status,
		Error:  message,
		Reason: err.Error(),
	}

	c.JSON(status, response)
}

func RespondWithSuccess(c *gin.Context, status int, message string) {
	response := &Response{
		Status:  status,
		Message: message,
	}

	c.JSON(status, response)
}

func RespondWithGet(c *gin.Context, status int, product []types.Product) {
	response := &GetResponse{
		Status: status,
		Data:   product,
	}
	c.JSON(status, response)
}
