package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func TranscriptRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/transcript", controller.GetTranscripts())
	incomingRoutes.GET("/transcript/transcript_id", controller.GetTranscript())
	incomingRoutes.POST("/transcript", controller.CreateTranscript())
	//incomingRoutes.PATCH("/transcript/transcript_id", controller.UpdateTranscript())
}
