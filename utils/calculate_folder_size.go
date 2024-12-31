package utils

import (
	"log"
	"os"
	"path/filepath"
)

// CalculateFolderSize calculates the total size of a directory recursively.
// It traverses the directory using filepath.Walk and sums up the sizes of all files.
func CalculateFolderSize(dir string) (int64, error) {
	var totalSize int64

	// Walk through the directory and process each file or folder.
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return err // Return the error to stop further processing
		}
		// Add file size to totalSize if it's not a directory
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil // Continue processing
	})
	return totalSize, err
}
