package main

import (
	"golang-au-backend/middleware"
	"golang-au-backend/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.AppRoutes(router)
	routes.NewsRoutes(router)
	routes.PrefRoutes(router)

	router.Run(":" + port)
}
