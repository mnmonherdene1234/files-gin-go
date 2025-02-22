package utils

import "os"

// ListFiles returns a list of file names in the specified directory.
// It skips directories and only includes regular files.
func ListFiles(directoryPath string) ([]string, error) {
	var fileNames []string

	// Read the directory entries
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// Extract the file names
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	return fileNames, nil
}
