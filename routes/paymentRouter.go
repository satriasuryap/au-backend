package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/payment/:user_id", controller.GetPayment())
	incomingRoutes.POST("/payment/:user_id", controller.CreatePayment())
	incomingRoutes.PATCH("/payment/:payment_id", controller.UpdatePayment())
	incomingRoutes.DELETE("/payment/:payment_id", controller.DeletePayment())
}
