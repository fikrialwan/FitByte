package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(server *gin.Engine, fileController controller.FileController, jwtService service.JwtService) {
	routes := server.Group("/v1")
	routes.Use(middlewares.Authenticate(jwtService))
	routes.POST("/file", fileController.UploadFile)
}
