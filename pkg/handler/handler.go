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
	err := ctx.ShouldBind(data)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, err.Error())
		return true
	}

	isError := validator.Check(data)
	if isError {
		ResponseError(ctx, http.StatusBadRequest, "Request form error")
		return true
	}

	return false
}
