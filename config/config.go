package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Global variables to store configuration values
var (
	IsGinReleaseMode     bool   // Indicates if the GIN framework is running in release mode (production)
	APIKeyEnabled        bool   // Flag to enable or disable API key authentication
	APIKeyHeader         string // Header name for the API key in HTTP requests
	APIKey               string // API key used for authentication of external/internal API calls
	ServerPort           string // Port on which the server listens for incoming HTTP requests
	FilesDir             string // Directory path where uploaded or static files are stored
	StaticFilesServePath string // URL path from which static files will be served
	IsServeStaticFiles   bool   // Flag to enable or disable serving of static files
)

// LoadConfig initializes the configuration by loading environment variables
// from a .env file. If the .env file is missing, it uses system environment variables.
func LoadConfig() error {
	// Attempt to load environment variables from a .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
		// Proceed with system environment variables if .env is not available.
	}

	// Retrieve and assign environment variables to configuration variables.
	IsGinReleaseMode = parseBool(getEnv("IS_GIN_RELEASE_MODE", "true"))
	APIKeyEnabled = parseBool(getEnv("API_KEY_ENABLED", "true"))
	APIKeyHeader = getEnv("API_KEY_HEADER", "X-API-Key")
	APIKey = getEnv("API_KEY", "123456890")
	ServerPort = getEnv("SERVER_PORT", "9935")
	FilesDir = getEnv("FILES_DIR", "./files")
	StaticFilesServePath = getEnv("STATIC_FILES_SERVE_PATH", "/files")
	IsServeStaticFiles = parseBool(getEnv("IS_SERVE_STATIC_FILES", "true"))

	// Log the loaded configuration values for easier debugging and tracking.
	log.Println("Is GIN Release Mode set to:", IsGinReleaseMode)
	log.Println("Server API key enabled set to:", APIKeyEnabled)
	log.Println("Server API key header set to:", APIKeyHeader)
	log.Println("Server API key set to:", APIKey)
	log.Println("Server Port set to:", ServerPort)
	log.Println("Files Directory set to:", FilesDir)
	log.Println("Static Files Serve Path set to:", StaticFilesServePath)
	log.Println("Is Serve Static Files enabled:", IsServeStaticFiles)

	return nil
}

// getEnv retrieves the value of the environment variable with the given key.
// If the environment variable is not set, it returns the provided default value.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue // Use default if no environment variable is set.
	}
	return value
}

// parseBool converts a string to a boolean.
// Returns true if the input is "true" (case-insensitive), otherwise false.
func parseBool(value string) bool {
	return strings.ToLower(value) == "true"
}
