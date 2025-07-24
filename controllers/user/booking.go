package user

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"
	"time"

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

	// Parse input.Date to time.Time
	bookingDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. Expected YYYY-MM-DD"})
		return
	}

	// Combine date and start time into full UTC timestamp
	fullStart := time.Date(
		bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		input.StartTime.Hour(), input.StartTime.Minute(), 0, 0, time.UTC,
	)

	// Ensure booking is at least 1 hour in the future
	now := time.Now().UTC()
	if fullStart.Before(now.Add(1 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot book past or near-time slots. Bookings must be at least 1 hour ahead."})
		return
	}

	// Check for existing overlapping bookings
	var existing models.Booking
	err = db.DB.Where(
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

	// Create booking
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

// func CancelBooking(c *gin.Context) {
// 	bookingId := c.Param("id")

// 	userIdVal, exist := c.Get("user_id")
// 	if !exist {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	userId := userIdVal.(uint)

// 	// Fetch the booking with related User and Turf
// 	var booking models.Booking
// 	if err := db.DB.Preload("User").Preload("Turf").
// 		First(&booking, "id = ? AND user_id = ?", bookingId, userId).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
// 		return
// 	}

// 	// Check if already cancelled
// 	if booking.Status == "cancelled" {
// 		c.JSON(http.StatusConflict, gin.H{"error": "booking is already cancelled"})
// 		return
// 	}

// 	// Prevent cancelling past or ongoing bookings
// 	now := time.Now()
// 	bookingDate, err := time.Parse("2006-01-02", booking.Date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid booking date format"})
// 		return
// 	}
// 	if bookingDate.Before(now) || (bookingDate.Equal(now) && booking.EndTime.Before(now)) {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "cannot cancel past or ongoing bookings"})
// 		return
// 	}

// 	// Update booking status to "cancelled"
// 	booking.Status = "cancelled"
// 	if err := db.DB.Save(&booking).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking status"})
// 		return
// 	}

// 	// Return full booking details after cancellation
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "booking cancelled successfully",
// 		"booking": booking,
// 	})
// }

func CancelBooking(c *gin.Context) {
	// Parse booking ID from URL
	bookingId := c.Param("id")

	// Get user ID from context (set by JWT middleware)
	userIdVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userId := userIdVal.(uint)

	// Fetch the booking for that user
	var booking models.Booking
	if err := db.DB.Preload("User").Preload("Turf").
		First(&booking, "id = ? AND user_id = ?", bookingId, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found or does not belong to the user"})
		return
	}

	// Check if already cancelled
	if booking.Status == "cancelled" {
		c.JSON(http.StatusConflict, gin.H{"error": "booking already cancelled"})
		return
	}

	// Parse booking.Date (YYYY-MM-DD)
	parsedDate, err := time.Parse("2006-01-02", booking.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid booking date format"})
		return
	}

	// Combine date with EndTime to form complete booking end time
	bookingEnd := time.Date(
		parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		booking.EndTime.Hour(), booking.EndTime.Minute(), booking.EndTime.Second(), 0, time.UTC,
	)

	// Check if booking has already ended
	if bookingEnd.Before(time.Now().UTC()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot cancel past or ongoing bookings"})
		return
	}

	// Mark as cancelled
	booking.Status = "cancelled"
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "booking cancelled successfully",
		"booking": booking,
	})
}
