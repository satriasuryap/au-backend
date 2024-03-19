package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AppRoutesRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/app", controller.GetApps())
	incomingRoutes.GET("/app/:app_id", controller.GetApp())
	incomingRoutes.POST("/app", controller.CreateApp())
	incomingRoutes.PATCH("/app/:app_id", controller.UpdateApp())
}
