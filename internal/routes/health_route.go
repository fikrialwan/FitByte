package routes

import (
	"github.com/fikrialwan/FitByte/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(router gin.IRouter, healthController controller.HealthController) {
	router.GET("/health", healthController.HealthCheck)
	router.GET("/ready", healthController.ReadinessCheck)
}
