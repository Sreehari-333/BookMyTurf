package public

import (
	"BookMyTurf/config"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Listing all turfs

func GetAllTurfs(c *gin.Context) {

	var turfs []models.Turf

	if err := config.DB.Find(&turfs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch turfs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"turfs": turfs,
	})
}
