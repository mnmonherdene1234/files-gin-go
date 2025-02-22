package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is a helper function for handling errors in a consistent way.
func ErrorResponse(c *gin.Context, status int, message string, err error) {
	log.Printf("%s: %v", message, err)
	c.JSON(status, gin.H{"error": message})
}
