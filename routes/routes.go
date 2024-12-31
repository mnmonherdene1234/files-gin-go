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
	serveStaticFiles(router)

	// Define protected routes with API key authentication middleware.
	registerProtectedRoutes(router)
}

// serveStaticFiles configures the router to serve static files
// from the specified directory if the feature is enabled.
func serveStaticFiles(router *gin.Engine) {
	if config.IsServeStaticFiles {
		router.Static(config.StaticFilesServePath, config.FilesDir)
		log.Printf("✅ Static file serving enabled: Path = '%s', Directory = '%s'",
			config.StaticFilesServePath, config.FilesDir)
	} else {
		log.Println("⚠️ Static file serving is currently disabled. To enable, set 'IS_SERVE_STATIC_FILES=true' in the .env file and restart the server.")
	}
}

// registerProtectedRoutes creates a group of routes protected by API key authentication.
func registerProtectedRoutes(router *gin.Engine) {
	// Apply API key authentication middleware to the route group.
	var protected *gin.RouterGroup

	if config.APIKeyEnabled {
		protected = router.Group("/", middlewares.APIKeyAuthMiddleware(config.APIKey))
	} else {
		protected = router.Group("/")
	}

	// Define protected endpoints for file management.
	protected.POST("/upload", handlers.UploadFileHandler)   // Endpoint for file uploads.
	protected.DELETE("/delete", handlers.DeleteFileHandler) // Endpoint for deleting files.
	protected.GET("/size", handlers.SizeHandler)
}
