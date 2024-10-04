package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	APIKeyHeader = "X-API-Key"
)

func LoadConfig() (*SettingsModel, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
	}

	configModel := &SettingsModel{
		APIKey:     os.Getenv("API_KEY"),
		ServerPort: os.Getenv("SERVER_PORT"),
		FilesDir:   os.Getenv("FILES_DIR"),
	}

	// Set default values if not provided
	if configModel.APIKey == "" {
		configModel.APIKey = "123456890"
		log.Println("Server API key set to " + configModel.APIKey)
	}
	if configModel.ServerPort == "" {
		configModel.ServerPort = "9935"
		log.Println("Server Port set to " + configModel.ServerPort)
	}
	if configModel.FilesDir == "" {
		configModel.FilesDir = "./files"
		log.Println("Files Directory set to " + configModel.FilesDir)
	}

	return configModel, nil
}
