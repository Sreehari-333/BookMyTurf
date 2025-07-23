package main

import (
	"BookMyTurf/db"
	"BookMyTurf/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	db.ConnectDB() // Connect to DB

	router := gin.Default() // Create router

	router.Static("/static", "./template")

	//  routes

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	routes.PublicRoutes(router)

	router.Run(":8080") //  Starting server

}
