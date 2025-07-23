package user

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchTurfs(c *gin.Context) {
	location := c.Query("location")
	date := c.Query("date")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if location == "" || date == "" || startTime == "" || endTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "location, date, start_time, and end_time are required",
		})
		return
	}

	var turfs []models.Turf

	// Find turfs that are NOT booked during the given time range
	subQuery := db.DB.Table("bookings").
		Select("turf_id").
		Where("date = ?", date).
		Where("NOT (end_time <= ? OR start_time >= ?)", startTime, endTime)

	err := db.DB.
		Where("location LIKE ?", "%"+location+"%").
		Where("id NOT IN (?)", subQuery).
		Find(&turfs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search turfs",
		})
		return
	}

	c.JSON(http.StatusOK, turfs)
}
