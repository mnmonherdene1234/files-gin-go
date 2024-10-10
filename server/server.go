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
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s", config.ServerPort)

	// Initialize the engine
	engine := gin.New()
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.Use(gin.Logger(), gin.Recovery())

	// Set up middleware and routes
	routes.Setup(engine)

	// Start the server
	if err := engine.Run(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
