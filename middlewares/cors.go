package middlewares

import (
	"github.com/fikrialwan/FitByte/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware with configurable settings
func CORS(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.GetCORSAllowedOrigins(),
		AllowMethods:     cfg.GetCORSAllowedMethods(),
		AllowHeaders:     cfg.GetCORSAllowedHeaders(),
		ExposeHeaders:    cfg.GetCORSExposeHeaders(),
		AllowCredentials: cfg.CORSAllowCredentials,
		MaxAge:           cfg.GetCORSMaxAge(),
	})
}
