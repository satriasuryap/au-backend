package routes

import (
	controller "golang-au-backend/controllers"

	"github.com/gin-gonic/gin"
)

func CoursesRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/course", controller.GetCourses())
	incomingRoutes.GET("/course/:course_id", controller.GetCourse())
	incomingRoutes.POST("/course", controller.CreateCourses())
	incomingRoutes.PATCH("/course/:course_id", controller.UpdateCourses())
	incomingRoutes.DELETE("/course/:course_id", controller.DeleteCourse())
}
