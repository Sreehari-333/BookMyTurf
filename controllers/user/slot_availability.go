package user

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAvailableSlots(c *gin.Context) {
	turfIDStr := c.Query("turf_id")
	dateStr := c.Query("date")

	// Validation
	if turfIDStr == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
		return
	}

	turfID, err := strconv.Atoi(turfIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid turf_id"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (yyyy-mm-dd)"})
		return
	}

	// Slot range config (6 AM to 11 PM, 1-hour slots)
	startHour := 6
	endHour := 23

	var availableSlots []map[string]string

	for hour := startHour; hour < endHour; hour++ {
		startTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)
		endTime := startTime.Add(time.Hour)

		var count int64
		err := db.DB.Model(&models.Booking{}).
			Where("turf_id = ? AND start_time < ? AND end_time > ?", turfID, endTime, startTime).
			Count(&count).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if count == 0 {
			availableSlots = append(availableSlots, map[string]string{
				"start_time": startTime.Format("15:04"),
				"end_time":   endTime.Format("15:04"),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"available_slots": availableSlots})
}
