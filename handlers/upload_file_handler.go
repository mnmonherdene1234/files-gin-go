package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	"github.com/mnmonherdene1234/files-gin-go/utils"
)

// UploadFileHandler handles file uploads without size limitation.
// @Summary Upload a file
// @Description Upload a large file to the server
// @Tags files
// @Accept  multipart/form-data
// @Produce json
// @Param   X-API-Key header string true "API Key"
// @Param   file formData file true "File to upload"
// @Success 200 {object} map[string]string "File uploaded successfully"
// @Failure 400 {object} map[string]string "No file received"
// @Failure 500 {object} map[string]string "Failed to create upload directory or Failed to save the file"
// @Router /upload [post]
func UploadFileHandler(c *gin.Context) {
	// Allow for very large files (set a very large limit)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<63-1) // No limit

	// Retrieve the uploaded file from the form
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No file received", err)
		return
	}

	// Ensure the directory for saving files exists
	if err := createUploadDir(config.FilesDir); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create upload directory", err)
		return
	}

	// Generate a unique filename to prevent conflicts
	uniqueFilename := utils.GenerateUniqueFilename(uploadedFile.Filename)
	uploadFilePath := filepath.Join(config.FilesDir, uniqueFilename)

	// Save the uploaded file
	if err := c.SaveUploadedFile(uploadedFile, uploadFilePath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save the file", err)
		return
	}

	// Respond with success and file information
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": uniqueFilename,
	})
}

// createUploadDir ensures the upload directory exists, creating it if necessary.
func createUploadDir(directoryPath string) error {
	return os.MkdirAll(directoryPath, os.ModePerm)
}
