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
			// TODO: add error message
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			// TODO: add error message
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			// TODO: add error message
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			// TODO: add error message
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			// TODO: add error message
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("token", authHeader)
		ctx.Set("user_id", userId)
		ctx.Next()
	}
}
