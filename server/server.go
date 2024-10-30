package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	_ "github.com/mnmonherdene1234/files-gin-go/docs" // Swagger documentation
	"github.com/mnmonherdene1234/files-gin-go/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Start initializes and starts the Gin server with the necessary configuration and routes.
func Start() {
	// Load configuration from .env or system environment variables
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	// Set Gin mode based on the configuration (release or debug mode)
	setGinMode(config.IsGinReleaseMode)

	// Create a new Gin engine instance
	engine := gin.New()

	// Add middleware for logging, recovery, and Swagger UI
	registerMiddlewares(engine)

	// Setup application routes
	routes.Setup(engine)

	// Start the HTTP server on the configured port
	if err := engine.Run(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}

// setGinMode configures Gin to run in either release or debug mode.
func setGinMode(isReleaseMode bool) {
	if isReleaseMode {
		gin.SetMode(gin.ReleaseMode) // Production mode: disables debug output
	} else {
		gin.SetMode(gin.DebugMode) // Development mode: enables debug output
	}
}

// registerMiddlewares adds necessary middlewares such as logging, recovery, and Swagger.
func registerMiddlewares(engine *gin.Engine) {
	// Route to serve Swagger UI for API documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Use built-in Gin middleware: Logger and Recovery
	engine.Use(gin.Logger(), gin.Recovery())
}
