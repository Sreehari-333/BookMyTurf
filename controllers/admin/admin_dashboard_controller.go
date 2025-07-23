package admin

import (
	"BookMyTurf/db"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var turfCount int64
	var bookingCount int64
	var userCount int64
	var blockedCount int64
	var totalRevenue float64

	db.DB.Model(&models.Turf{}).Count(&turfCount)

	db.DB.Model(&models.Booking{}).Count(&bookingCount)

	db.DB.Model(&models.User{}).Count(&userCount)

	db.DB.Model(&models.BlockedSlot{}).Count(&blockedCount)

	type Result struct {
		Total float64
	}

	var result Result

	db.DB.
		Table("bookings").
		Select("COALESCE(SUM(turfs.price), 0) AS total").
		Joins("JOIN turfs ON turfs.id = bookings.turf_id").
		Scan(&result)

	totalRevenue = result.Total

	c.JSON(http.StatusOK, gin.H{
		"total_turfs":    turfCount,
		"total_bookings": bookingCount,
		"total_users":    userCount,
		"blocked_slots":  blockedCount,
		"revenue":        totalRevenue,
	})
}
