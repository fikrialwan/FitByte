package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router gin.IRouter, userController controller.UserController, jwtService service.JwtService) {

	router.POST("/login", userController.Login)
	router.POST("/register", userController.Register)

	userRoutes := router.Group("/user")
	userRoutes.Use(middlewares.Authenticate(jwtService))
	userRoutes.GET("/", userController.GetProfile)
}
