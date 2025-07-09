package routes

import (
	authcontroller "BookMyTurf/controllers/auth"

	"github.com/gin-gonic/gin"
)

// Auth routes

func AuthRoutes(router *gin.Engine) {

	auth := router.Group("/auth")
	{
		auth.POST("/register", authcontroller.Register)         // route for register
		auth.POST("/login", authcontroller.Login)               // route for login
		auth.POST("/refreshtoken", authcontroller.RefreshToken) // Generating new access token
		auth.POST("/logout", authcontroller.Logout)
	}
}
