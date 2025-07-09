package main

import (
	"BookMyTurf/config"
	"BookMyTurf/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB() // Connecting DB

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	routes.PublicRoutes(router)

	router.Run(":8080") // Starting server on default port
}
