package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
		filename := utils.GenerateUniqueFilename(file.Filename)
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
