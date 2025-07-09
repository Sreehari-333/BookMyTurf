package routes

import (
	publicController "BookMyTurf/controllers/public"
	"github.com/gin-gonic/gin"
)

// Listing all turfs

func PublicRoutes(router *gin.Engine) {
	router.GET("/turfs", publicController.GetAllTurfs)
	router.GET("/turfs/search", publicController.SearchTurfs)
}
