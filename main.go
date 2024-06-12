package main

import (
	"golang-au-backend/database"
	"golang-au-backend/middleware"
	"golang-au-backend/routes"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var appCollection *mongo.Collection = database.OpenCollection(database.Client, "app")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	routes.UserRoutes(router)
	routes.NewsRoutes(router)
	routes.CoursesRoutes(router)
	routes.PaymentRoutes(router)
	routes.AppRoutes(router)
	routes.PrefRoutes(router)
	routes.TranscriptRoutes(router)
	router.Use(middleware.Authentication())

	router.Run(":" + port)
}
