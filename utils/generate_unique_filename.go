package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// GenerateUniqueFilename creates a unique filename by appending a timestamp to the original filename.
// It preserves the original file extension.
// Example: "document_2024-10-04T15-30-25.123456789.pdf"
func GenerateUniqueFilename(originalFilename string) string {
	if originalFilename == "" {
		originalFilename = "file"
	}

	// Extract the file extension
	ext := filepath.Ext(originalFilename)

	// Get the base name without the extension
	baseName := strings.TrimSuffix(filepath.Base(originalFilename), ext)

	// Generate a timestamp with nanosecond precision in a readable format
	timestamp := time.Now().Format("2006-01-02T15-04-05.000000000")

	// Combine the base name, timestamp, and extension to form the unique filename
	uniqueFilename := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)

	return uniqueFilename
}
