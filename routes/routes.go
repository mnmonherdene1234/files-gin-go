package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	"github.com/mnmonherdene1234/files-gin-go/handlers"
	"github.com/mnmonherdene1234/files-gin-go/middlewares"
	"github.com/mnmonherdene1234/files-gin-go/utils"
)

// Setup initializes all routes and middleware for the Gin router.
func Setup(router *gin.Engine) {
	// Configure CORS settings for cross-origin requests.
	utils.ConfigureCORS(router)

	// Serve static files if enabled in the configuration.
	configureStaticFileServing(router)

	// Define protected routes with API key authentication middleware.
	defineProtectedRoutes(router)
}

// configureStaticFileServing configures the router to serve static files
// from the specified directory if the feature is enabled.
func configureStaticFileServing(router *gin.Engine) {
	if config.IsServeStaticFiles {
		router.Static(config.StaticFilesServePath, config.FilesDir)
		log.Printf("✅ Static file serving enabled: Path = '%s', Directory = '%s'",
			config.StaticFilesServePath, config.FilesDir)
	} else {
		log.Println("⚠️ Static file serving is currently disabled. To enable, set 'IS_SERVE_STATIC_FILES=true' in the .env file and restart the server.")
	}
}

// defineProtectedRoutes creates a group of routes protected by API key authentication.
func defineProtectedRoutes(router *gin.Engine) {
	// Apply API key authentication middleware to the route group.
	var protectedRoutes *gin.RouterGroup

	if config.APIKeyEnabled {
		protectedRoutes = router.Group("/", middlewares.APIKeyAuthMiddleware(config.APIKey))
	} else {
		protectedRoutes = router.Group("/")
	}

	// Define protected endpoints for file management.
	protectedRoutes.POST("/upload", handlers.UploadFileHandler)   // Endpoint for file uploads.
	protectedRoutes.DELETE("/delete", handlers.DeleteFileHandler) // Endpoint for deleting files.
	protectedRoutes.GET("/size", handlers.SizeHandler)            // Endpoint for getting file size.
	protectedRoutes.GET("/list", handlers.FilesListHandler)       // Endpoint for listing files.
}
