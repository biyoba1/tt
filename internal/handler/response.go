package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, statusCode int, message string) error {
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
	return errors.New(message)
}
