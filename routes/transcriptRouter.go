package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func TranscriptRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/transcript/:user_id", controller.GetTranscript())
	incomingRoutes.POST("/transcript/:user_id/:course_id", controller.CreateTranscript())
	incomingRoutes.PATCH("/transcript/:transcript_id", controller.UpdateTranscript())
	incomingRoutes.DELETE("/transcript/:transcript_id", controller.DeleteTranscript())
}
