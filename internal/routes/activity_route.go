package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterActivityRoutes(router gin.IRouter, activityController controller.ActivityController, jwtService service.JwtService) {
	router.Use(middlewares.Authenticate(jwtService))
	router.GET("/activity", activityController.GetActivity)
	router.POST("/activity", activityController.CreateActivity)
	router.PATCH("/activity/:activityId", activityController.UpdateActivity)
}
