package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	_ "github.com/mnmonherdene1234/files-gin-go/docs"
	"github.com/mnmonherdene1234/files-gin-go/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func Start() {
	// Load configuration
	settings, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s", settings.ServerPort)

	// Initialize the engine
	engine := gin.New()
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.Use(gin.Logger(), gin.Recovery())

	// Set up middleware and routes
	routes.Setup(engine, settings)

	// Start the server
	if err := engine.Run(":" + settings.ServerPort); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
