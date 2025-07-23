package user

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Booking turf

func BookTurf(c *gin.Context) {
	var input models.Booking

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid input",
			"details": err.Error(),
		})
		return
	}

	// Get user ID from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check for existing bookings that overlap the given time
	var existing models.Booking
	err := db.DB.Where(
		"turf_id = ? AND start_time < ? AND end_time > ?",
		input.TurfId, input.EndTime, input.StartTime,
	).First(&existing).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "time slot already booked"})
		return
	}

	// Check for blocked slots
	var blocked []models.BlockedSlot
	err = db.DB.Where(
		"turf_id = ? AND date = ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))",
		input.TurfId, input.Date,
		input.EndTime, input.EndTime,
		input.StartTime, input.StartTime,
	).Find(&blocked).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check blocked slots"})
		return
	}

	if len(blocked) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "time slot is blocked by admin"})
		return
	}

	// Create new booking
	booking := models.Booking{
		TurfId:    input.TurfId,
		UserId:    userId.(uint),
		Date:      input.Date,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    "booked",
	}

	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking successful",
		"booking": booking,
	})
}

// Get user booking history

func GetMyBookings(c *gin.Context) {

	userId, exist := c.Get("user_id")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var booking []models.Booking

	if err := db.DB.Preload("Turf").Where("user_id = ?", userId).Order("start_time asc").Find(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch  booking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Bookings": booking,
		"UserId":   userId,
	})
}

// Cancel Booking

func CancelBooking(c *gin.Context) {

	bookingId := c.Param("id") // get booking id from url

	userIdVal, exist := c.Get("user_id")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userId := userIdVal.(uint)

	var booking models.Booking

	if err := db.DB.First(&booking, "id = ? AND user_id = ?", bookingId, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "booking not found",
		})
		return
	}

	// changing status

	booking.Status = "cancelled"
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	// Delete

	if err := db.DB.Delete(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cancel booking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
	})

}
