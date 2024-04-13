package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PrefRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/pref", controller.GetPrefs())
	incomingRoutes.GET("/pref/:pref_id", controller.GetPref())
	incomingRoutes.POST("/pref", controller.CreatePref())
}
