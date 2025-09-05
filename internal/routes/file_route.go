package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router gin.IRouter, fileController controller.FileController, jwtService service.JwtService) {
	router.Use(middlewares.Authenticate(jwtService))
	router.POST("/file", fileController.UploadFile)
}
