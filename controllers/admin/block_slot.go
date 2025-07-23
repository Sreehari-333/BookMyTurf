package admin

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Blocking slots

func BlockSlot(c *gin.Context) {

	var input models.BlockedSlot

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	// Blocking slot

	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not block slot",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "slot block successfully",
	})

}

// Unblocking blocked slots

func UnBlockSlots(c *gin.Context) {

	blockedSlotId := c.Param("id")

	var blockedSlot models.BlockedSlot

	// checking if its exists

	if err := db.DB.First(&blockedSlot, blockedSlotId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "blocked slot not found",
		})
		return
	}

	// Unblocking the slot

	if err := db.DB.Delete(&blockedSlot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete blocked slot",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "slot unblocked successfully",
	})
}

// Listing all blocked slots

func GetBlockedSlots(c *gin.Context) {

	var blockedSlots []models.BlockedSlot

	if err := db.DB.Preload("Turf").Order("date asc").Find(&blockedSlots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch blocked turf slots",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Blocked turfs": blockedSlots,
	})
}
