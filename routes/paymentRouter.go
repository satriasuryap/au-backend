package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/payment", controller.GetPayments())
	incomingRoutes.GET("/payment/:payment_id", controller.GetPayment())
	incomingRoutes.POST("/payment", controller.CreatePayment())
	incomingRoutes.PATCH("/payment/:payment_id", controller.UpdatePayment())
}
