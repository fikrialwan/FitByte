package middlewares

import (
	"net/http"
	"strings"

	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
			})
			ctx.Abort()
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid",
			})
			ctx.Abort()
			return
		}

		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to extract user ID from token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("token", authHeader)
		ctx.Set("user_id", userId)
		ctx.Next()
	}
}
