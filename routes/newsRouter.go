package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func NewsRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/news", controller.GetNews())
	incomingRoutes.GET("/news/:news_id", controller.GetNewsOne())
	incomingRoutes.POST("/news", controller.CreateNews())
	incomingRoutes.PATCH("/news/:news_id", controller.UpdateNews())
	incomingRoutes.DELETE("/news/:news_id", controller.DeleteNews())

}
