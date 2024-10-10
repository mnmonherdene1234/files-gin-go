package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	APIKeyHeader = "X-API-Key"
)

var (
	APIKey     string
	ServerPort string
	FilesDir   string
)

func LoadConfig() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
		return err
	}

	APIKey = os.Getenv("API_KEY")
	ServerPort = os.Getenv("SERVER_PORT")
	FilesDir = os.Getenv("FILES_DIR")

	// Set default values if not provided
	if APIKey == "" {
		APIKey = "123456890"
		log.Println("Server API key set to " + APIKey)
	}
	if ServerPort == "" {
		ServerPort = "9935"
		log.Println("Server Port set to " + ServerPort)
	}
	if FilesDir == "" {
		FilesDir = "./files"
		log.Println("Files Directory set to " + FilesDir)
	}

	return nil
}
