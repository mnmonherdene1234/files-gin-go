package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	"github.com/mnmonherdene1234/files-gin-go/handlers"
	"github.com/mnmonherdene1234/files-gin-go/middlewares"
	"github.com/mnmonherdene1234/files-gin-go/utils"
)

func Setup(router *gin.Engine) {
	// Configure CORS
	utils.ConfigureCORS(router)

	// Serve static files from the files directory
	router.Static("/files", config.FilesDir)

	// Apply API key middleware
	protected := router.Group("/", middlewares.APIKeyAuthMiddleware(config.APIKey))

	// Protected routes
	protected.POST("/upload", handlers.UploadFileHandler)
	protected.DELETE("/delete", handlers.DeleteFileHandler)
}
