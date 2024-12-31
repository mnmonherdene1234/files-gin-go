package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mnmonherdene1234/files-gin-go/config"
	"github.com/mnmonherdene1234/files-gin-go/utils"
)

// SizeHandler responds with the total size of the configured folder.
//
// @Summary      Get folder size
// @Description  Calculates and returns the total size of the folder specified in the configuration.
// @Tags         Files
// @Produce      json
// @Param   	 X-API-Key header string true "API Key"
// @Success      200 {object} map[string]interface{} "Folder size in bytes"
// @Failure      500 {object} map[string]interface{} "Error message with details"
// @Router       /size [get]
func SizeHandler(c *gin.Context) {
	// Retrieve the folder path from the configuration
	folderPath := config.FilesDir
	if folderPath == "" {
		c.JSON(500, gin.H{
			"error": "Folder path not configured",
		})
		return
	}

	// Calculate the size of the folder
	size, err := utils.CalculateFolderSize(folderPath)
	if err != nil {
		log.Printf("Error calculating folder size: %v", err)
		c.JSON(500, gin.H{
			"error":  "Failed to calculate folder size",
			"detail": err.Error(),
		})
		return
	}

	// Respond with the folder size in bytes
	c.JSON(200, gin.H{
		"size": size, // Size in bytes
	})
}
