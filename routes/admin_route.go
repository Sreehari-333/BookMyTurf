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

	admin.POST("/turf", adminController.CreateTurf)       // add turf
	admin.PUT("/turf/:id", adminController.UpdateTurf)    // update turf
	admin.DELETE("/turf/:id", adminController.DeleteTurf) // Delete turf

}
