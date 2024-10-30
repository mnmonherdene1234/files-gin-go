package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
)

func APIKeyAuthMiddleware(expectedAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(config.APIKeyHeader)

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
