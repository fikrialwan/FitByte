package validator

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateJSONMiddleware creates a middleware that validates JSON payloads
func ValidateJSONMiddleware(schema ValidationSchema) gin.HandlerFunc {
	validator := NewJSONValidator(schema)
	
	return func(ctx *gin.Context) {
		// Only validate for requests with JSON content
		contentType := ctx.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			ctx.Next()
			return
		}

		// Read the raw body
		body, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			ctx.Abort()
			return
		}

		// Skip validation for empty body
		if len(body) == 0 {
			ctx.Next()
			return
		}

		// Validate the JSON
		if err := validator.ValidateJSON(body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			ctx.Abort()
			return
		}

		// Reset the body for the next handler
		ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))
		ctx.Next()
	}
}

// ValidateJSON is a helper function for manual validation in controllers
func ValidateJSON(jsonData []byte, schema ValidationSchema) error {
	validator := NewJSONValidator(schema)
	return validator.ValidateJSON(jsonData)
}
