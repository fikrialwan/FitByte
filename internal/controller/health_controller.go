package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	db           *gorm.DB
	cacheService service.CacheService
	fileService  service.FileService
}

func NewHealthController(db *gorm.DB, cacheService service.CacheService, fileService service.FileService) HealthController {
	return HealthController{
		db:           db,
		cacheService: cacheService,
		fileService:  fileService,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h HealthController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "FitByte API",
		"version":   "1.0.0",
	})
}

// ReadinessCheck godoc
// @Summary Readiness check endpoint
// @Description Returns the readiness status of the application with actual health checks
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ready [get]
func (h HealthController) ReadinessCheck(ctx *gin.Context) {
	checks := make(map[string]interface{})
	allHealthy := true

	// Database connectivity check
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			checks["database"] = map[string]interface{}{
				"status": "error",
				"error":  "failed to get database instance: " + err.Error(),
			}
			allHealthy = false
		} else if err := sqlDB.Ping(); err != nil {
			checks["database"] = map[string]interface{}{
				"status": "error",
				"error":  "database ping failed: " + err.Error(),
			}
			allHealthy = false
		} else {
			checks["database"] = map[string]interface{}{
				"status": "healthy",
			}
		}
	} else {
		checks["database"] = map[string]interface{}{
			"status": "error",
			"error":  "database not initialized",
		}
		allHealthy = false
	}

	// Redis connectivity check
	if h.cacheService != nil {
		// Test Redis with a simple ping operation
		testKey := "health_check_" + time.Now().Format("20060102150405")
		testValue := "ping"

		// Try to set and get a test value
		err := h.cacheService.Set(testKey, testValue, 10*time.Second)
		if err != nil {
			checks["redis"] = map[string]interface{}{
				"status": "error",
				"error":  "redis set operation failed: " + err.Error(),
			}
			allHealthy = false
		} else {
			// Try to get the value back
			_, err := h.cacheService.Get(testKey)
			if err != nil {
				checks["redis"] = map[string]interface{}{
					"status": "error",
					"error":  "redis get operation failed: " + err.Error(),
				}
				allHealthy = false
			} else {
				checks["redis"] = map[string]interface{}{
					"status": "healthy",
				}
				// Clean up test key
				h.cacheService.Delete(testKey)
			}
		}
	} else {
		checks["redis"] = map[string]interface{}{
			"status": "error",
			"error":  "cache service not initialized",
		}
		allHealthy = false
	}

	// MinIO connectivity check
	if h.fileService != nil {
		// Test MinIO connectivity using the interface method
		err := h.fileService.CheckConnectivity(context.Background())
		if err != nil {
			checks["minio"] = map[string]interface{}{
				"status": "error",
				"error":  "MinIO connectivity check failed: " + err.Error(),
			}
			allHealthy = false
		} else {
			checks["minio"] = map[string]interface{}{
				"status": "healthy",
			}
		}
	} else {
		checks["minio"] = map[string]interface{}{
			"status": "error",
			"error":  "file service not initialized",
		}
		allHealthy = false
	}

	// Determine overall status and HTTP status code
	status := "ready"
	httpStatus := http.StatusOK
	if !allHealthy {
		status = "not ready"
		httpStatus = http.StatusServiceUnavailable
	}

	ctx.JSON(httpStatus, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	})
}
