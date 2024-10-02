package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	APIKey     string
	ServerPort string
	FilesDir   string
}

func main() {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s", config.ServerPort)

	// Initialize the router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Set up middleware and routes
	setupRoutes(router, config)

	// Start the server
	if err := router.Run(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
	}

	config := &Config{
		APIKey:     os.Getenv("API_KEY"),
		ServerPort: os.Getenv("SERVER_PORT"),
		FilesDir:   os.Getenv("FILES_DIR"),
	}

	// Set default values if not provided
	if config.APIKey == "" {
		config.APIKey = "123456890"
		log.Println("Server API key set to " + config.APIKey)
	}
	if config.ServerPort == "" {
		config.ServerPort = "9935"
		log.Println("Server Port set to " + config.ServerPort)
	}
	if config.FilesDir == "" {
		config.FilesDir = "./files"
		log.Println("Files Directory set to " + config.FilesDir)
	}

	return config, nil
}

func setupRoutes(router *gin.Engine, config *Config) {
	// Configure CORS
	configureCORS(router)

	// Serve static files from the files directory
	router.Static("/files", config.FilesDir)

	// Apply API key middleware
	protected := router.Group("/", APIKeyAuthMiddleware(config.APIKey))

	// Protected routes
	protected.POST("/upload", UploadFileHandler(config.FilesDir))
}

func configureCORS(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
}

func UploadFileHandler(filesDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the file from form data
		file, err := c.FormFile("file")
		if err != nil {
			log.Printf("File retrieval error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
			return
		}

		// Ensure the files directory exists
		if err := os.MkdirAll(filesDir, os.ModePerm); err != nil {
			log.Printf("Upload directory creation error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Generate a unique filename
		filename := generateUniqueFilename(file.Filename)
		uploadFilePath := filepath.Join(filesDir, filename)

		// Save the uploaded file
		if err := c.SaveUploadedFile(file, uploadFilePath); err != nil {
			log.Printf("File saving error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
			return
		}

		// Return success response with the file URL
		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"filename": filename,
		})
	}
}

func generateUniqueFilename(originalFilename string) string {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(originalFilename)
	name := filepath.Base(originalFilename[:len(originalFilename)-len(ext)])
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
