package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	"github.com/mnmonherdene1234/files-gin-go/utils"
)

// FilesListHandler handles the request to list all files in the configured directory.
// FilesListHandler handles the request to list all files in a specified directory.
// @Summary List files
// @Description Retrieves a list of all files in the configured directory.
// @Tags files
// @Produce json
// @Param   	 X-API-Key header string true "API Key"
// @Success 200 {array} string "List of files"
// @Failure 500 {object} map[string]string "Failed to list files"
// @Router /list-files [get]
func FilesListHandler(c *gin.Context) {
	// Retrieve the directory path from the configuration
	directoryPath := config.FilesDir

	// List all files in the directory
	files, err := utils.ListFiles(directoryPath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list files", err)
		return
	}

	// Respond with the list of files
	c.JSON(http.StatusOK, files)
}
