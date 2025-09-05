package handler

import (
	"net/http"

	"github.com/fikrialwan/FitByte/pkg/validator"
	"github.com/gin-gonic/gin"
)

// Handle response error
func ResponseError(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"error": message})
}

// Handle response success
func ResponseSuccess(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}

// Parse request data & validate struct
func BindAndValidate(ctx *gin.Context, data interface{}) bool {
	// Check for empty body first
	if ctx.Request.ContentLength == 0 {
		ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return true
	}

	// First try to bind - this will catch JSON parsing errors and basic type mismatches
	err := ctx.ShouldBind(data)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return true
	}

	// Then validate the struct with our validation rules
	isError := validator.Check(data)
	if isError {
		ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return true
	}

	return false
}
