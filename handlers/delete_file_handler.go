package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	"github.com/mnmonherdene1234/files-gin-go/utils"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteFileRequest represents the expected JSON structure for the delete file request.
type DeleteFileRequest struct {
	Filename string `json:"filename" binding:"required"`
}

// DeleteFileHandler handles file deletions by filename from JSON request body.
// @Summary Delete a file by filename
// @Description Delete a file from the server using the filename provided in the JSON body
// @Tags files
// @Accept json
// @Produce json
// @Param   X-API-Key header string true "API Key"
// @Param   request body DeleteFileRequest true "Delete file request body"
// @Success 200 {object} map[string]string "File deleted successfully"
// @Failure 400 {object} map[string]string "Invalid request body or Filename not provided"
// @Failure 404 {object} map[string]string "File not found"
// @Failure 500 {object} map[string]string "Failed to delete the file"
// @Router /delete [delete]
func DeleteFileHandler(c *gin.Context) {
	// Parse the JSON body to get the filename
	var req DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body or Filename not provided", err)
		return
	}

	// Construct the full path to the file
	filePath := filepath.Join(config.FilesDir, req.Filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found", err)
		return
	}

	// Attempt to delete the file
	if err := os.Remove(filePath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the file", err)
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message":  "File deleted successfully",
		"filename": req.Filename,
	})
}
