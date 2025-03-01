package utils

import "github.com/gin-gonic/gin"

func RespondWithError(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{
		"status": status,
		"error":  message,
		"reason": err.Error(),
	})
}
