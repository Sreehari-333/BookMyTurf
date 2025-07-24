package routes

import (
	adminController "BookMyTurf/controllers/admin"
	"BookMyTurf/middleware"
	"github.com/gin-gonic/gin"
)

// Admin Routes

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")

	admin.Use(middleware.AuthMiddleware(), middleware.IsAdmin())
	{

		// Turf management

		admin.POST("/turf", adminController.CreateTurf)       // Creating turf
		admin.PUT("/turf/:id", adminController.UpdateTurf)    // Updating turf
		admin.DELETE("/turf/:id", adminController.DeleteTurf) // Deleting turf

		// Slot management

		admin.POST("/blockslot", adminController.BlockSlot)          // Blocking slots
		admin.GET("/blockslot", adminController.GetBlockedSlots)     // Get blocked slots
		admin.DELETE("/blockslot/:id", adminController.UnBlockSlots) // Unblock slots

		// Admin dashboard stats

		admin.GET("/dashboardstats", adminController.GetDashboardStats) // Get dashboard stats

		// User management

		admin.GET("/users", adminController.GetAllUsers)             // Get all users
		admin.PUT("/users/block/:id", adminController.BlockUser)     // Block users
		admin.PUT("/users/unblock/:id", adminController.UnblockUser) // Unblock users
		admin.PUT("/users/update", adminController.AdminUpdateUser)

		// Booking management

		admin.GET("/allbookings", adminController.GetAllBookings) // Get all bookings
	}
}
