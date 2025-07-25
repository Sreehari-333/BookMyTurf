package routes

import (
	"BookMyTurf/controllers/payment"
	publicControllers "BookMyTurf/controllers/user"
	"BookMyTurf/middleware"
	"BookMyTurf/services"

	"github.com/gin-gonic/gin"
)

// User routes

func UserRoutes(router *gin.Engine) {

	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware())
	user.POST("/bookings", services.BookTurf)
	user.GET("/bookings", services.GetMyBookings)
	user.DELETE("/bookings/:id", services.CancelBooking)
	user.GET("/profile", publicControllers.GetProfile)
	user.PUT("/profile", publicControllers.UpdateProfile)
	user.GET("/checkslots", services.GetAvailableSlots)

	// Payment

	user.POST("/payment", payment.CreatePayment)
	user.GET("/payment/:id", payment.GetPaymentByID)
}
