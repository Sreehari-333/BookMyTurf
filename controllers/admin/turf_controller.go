package admin

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Creating turf

func CreateTurf(c *gin.Context) {

	var turf models.Turf

	err := c.ShouldBindBodyWithJSON(&turf)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.DB.Create(&turf).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "turf created successfully",
		"turf":    turf,
	})
}

// Update turf

func UpdateTurf(c *gin.Context) {

	turfId := c.Param("id")

	var existingTirf models.Turf

	if err := db.DB.First(&existingTirf, turfId).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "turf not found",
		})
		return
	}

	var updatedData models.Turf

	if err := c.ShouldBindBodyWithJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := db.DB.Model(&existingTirf).Updates(updatedData).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update turf",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "turf updated successfuly",
	})
}

// Delete Turf

func DeleteTurf(c *gin.Context) {

	turfId := c.Param("id")

	var turf models.Turf

	if err := db.DB.First(&turf, turfId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "turf not found",
		})
		return
	}

	if err := db.DB.Delete(&turf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete turf",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "turf deleted successfuly",
	})

}
