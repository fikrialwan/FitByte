package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(server *gin.Engine, userController controller.UserController, jwtService service.JwtService) {
	routes := server.Group("/v1")

	routes.POST("/login", userController.Login)
	routes.POST("/register", userController.Register)

	userRoutes := routes.Group("/user")
	userRoutes.Use(middlewares.Authenticate(jwtService))
	userRoutes.GET("/", userController.GetProfile)
}
