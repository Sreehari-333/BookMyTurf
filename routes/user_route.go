package routes

import (
	publicControllers "BookMyTurf/controllers/user"
	"BookMyTurf/middleware"
	"github.com/gin-gonic/gin"
)

// User routes

func UserRoutes(router *gin.Engine) {

	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware())
	user.POST("/bookings", publicControllers.BookTurf)
	user.GET("/bookings", publicControllers.GetMyBookings)
	user.DELETE("/bookings/:id", publicControllers.CancelBooking)
	user.GET("/profile", publicControllers.GetProfile)
	user.PUT("/profile", publicControllers.UpdateProfile)
	user.GET("/checkslots", publicControllers.GetAvailableSlots)
}
