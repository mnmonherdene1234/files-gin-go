package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func APIKeyAuthMiddleware(expectedAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			log.Println("Missing API key")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			return
		}
		if apiKey != expectedAPIKey {
			log.Println("Invalid API key")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	}
}
