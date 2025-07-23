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

		admin.POST("/turf", adminController.CreateTurf)
		admin.PUT("/turf/:id", adminController.UpdateTurf)
		admin.DELETE("/turf/:id", adminController.DeleteTurf)

		// Slot management

		admin.POST("/blockslot", adminController.BlockSlot)
		admin.GET("/blockslot", adminController.GetBlockedSlots)
		admin.DELETE("/blockslot/:id", adminController.UnBlockSlots)

		// Admin dashboard stats

		admin.GET("/dashboardstats", adminController.GetDashboardStats)

		// User management

		admin.GET("/users", adminController.GetAllUsers)
		admin.PUT("/users/block/:id", adminController.BlockUser)
		admin.PUT("/users/unblock/:id", adminController.UnblockUser)

		// Booking management

		admin.GET("/allbookings", adminController.GetAllBookings)
	}
}
