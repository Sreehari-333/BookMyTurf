package public

import (
	"BookMyTurf/config"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Booking turf

func BookTurf(c *gin.Context) {

	var bookingInput models.Booking

	if err := c.ShouldBindJSON(&bookingInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid input",
			"details": err.Error(),
		})
		return
	}

	userId, exist := c.Get("user_id")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	bookingInput.UserId = userId.(uint)
	bookingInput.Status = "booked"

	// Checking if time slot already booked for the turf

	var existingBooking models.Booking

	err := config.DB.Where("turf_id = ? AND start_time < ? AND end_time > ? ", bookingInput.TurfId, bookingInput.EndTime, bookingInput.StartTime).First(&existingBooking).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "time slot already booked",
		})
		return
	}

	// Creating booking

	if err := config.DB.Create(&bookingInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create booking",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking successfull",
		"booking": bookingInput,
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

	if err := config.DB.Preload("Turf").Where("user_id = ?", userId).Order("start_time asc").Find(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch  booking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Bookings": booking,
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

	if err := config.DB.First(&booking, "id = ? AND user_id = ?", bookingId, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "booking not found",
		})
		return
	}

	// changing status

	booking.Status = "cancelled"
	if err := config.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	// Delete

	if err := config.DB.Delete(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cancel booking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
	})

}
