package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterActivityRoutes(server *gin.Engine, activityController controller.ActivityController, jwtService service.JwtService) {
	routes := server.Group("/v1")
	routes.Use(middlewares.Authenticate(jwtService))
	routes.POST("/activity", activityController.CreateActivity)
}
