package controller

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/pkg/handler"
	"github.com/fikrialwan/FitByte/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityController struct {
	activityService service.ActivityService
}

func NewActivityController(activityService service.ActivityService) ActivityController {
	return ActivityController{activityService}
}

// GetActivities godoc
// @Summary      Get all activities
// @Description  Retrieve activities with optional filters (limit, offset, activity type, date range, calories burned range).
// @Tags         activities
// @Accept       json
// @Produce      json
// @Param        limit              query     int     false  "Limit (default: 5)"         minimum(1) maximum(100)
// @Param        offset             query     int     false  "Offset (default: 0)"        minimum(0)
// @Param        activityType       query     string  false  "Activity Type" Enums(Walking,Running,Yoga,Stretching,Cycling,Swimming,Dancing,Hiking,HIIT,JumpRope)
// @Param        doneAtFrom         query     string  false  "Filter from date (ISO8601)" format(date-time)
// @Param        doneAtTo           query     string  false  "Filter to date (ISO8601)" format(date-time)
// @Param        caloriesBurnedMin  query     int     false  "Minimum calories burned"
// @Param        caloriesBurnedMax  query     int     false  "Maximum calories burned"
// @Success 200 {array} dto.ActivityResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /activity [get]
func (c ActivityController) GetActivity(ctx *gin.Context) {
	var filter dto.ActivityFilter

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		handler.ResponseError(ctx, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Validate filter parameters
	if filter.Limit < 0 || filter.Limit > 100 {
		handler.ResponseError(ctx, http.StatusBadRequest, "Limit must be between 0 and 100")
		return
	}
	if filter.Offset < 0 {
		handler.ResponseError(ctx, http.StatusBadRequest, "Offset must be non-negative")
		return
	}
	if filter.CaloriesBurnedMin < 0 || filter.CaloriesBurnedMax < 0 {
		handler.ResponseError(ctx, http.StatusBadRequest, "Calories burned values must be non-negative")
		return
	}
	if filter.CaloriesBurnedMin > 0 && filter.CaloriesBurnedMax > 0 && filter.CaloriesBurnedMin > filter.CaloriesBurnedMax {
		handler.ResponseError(ctx, http.StatusBadRequest, "Minimum calories burned cannot be greater than maximum")
		return
	}

	log.Printf("DEBUG FILTER: %+v\n", filter)

	userID := ctx.GetString("user_id")
	res, err := c.activityService.GetActivity(filter, userID)
	if err != nil {
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, res)
}

// CreateActivity godoc
// @Summary Create activity
// @Description Create a new activity with automatic calorie calculation
// @Tags activities
// @Accept json
// @Produce json
// @Param request body dto.ActivityRequest true "Activity data"
// @Success 201 {object} dto.CreateActivityResponse "Activity created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request - Invalid input format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /activity [post]
func (c ActivityController) CreateActivity(ctx *gin.Context) {
	var request dto.ActivityRequest
	if handler.BindAndValidate(ctx, &request) {
		return
	}

	response, err := c.activityService.CreateActivity(request, ctx.GetString("user_id"))
	if err != nil {
		// If it's an activity type validation error, return 400
		if strings.Contains(err.Error(), "invalid activity type") {
			handler.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		// For other errors, return 500
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusCreated, response)
}

// UpdateActivity godoc
// @Summary Update activity
// @Description Update an existing activity with automatic calorie recalculation
// @Tags activities
// @Accept json
// @Produce json
// @Param activityId path string true "Activity ID"
// @Param request body dto.ActivityUpdateRequest true "Activity update data"
// @Success 200 {object} dto.ActivityResponse "Activity updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request - Invalid input format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Activity not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /activity/{activityId} [patch]
func (c ActivityController) UpdateActivity(ctx *gin.Context) {
	// Check content type first
	contentType := ctx.GetHeader("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "application/json") {
		handler.ResponseError(ctx, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}

	activityID := ctx.Param("activityId")
	if activityID == "" {
		handler.ResponseError(ctx, http.StatusBadRequest, "Activity ID is required")
		return
	}

	// Validate UUID format - return 404 for invalid format as it means "not found"
	if _, err := uuid.Parse(activityID); err != nil {
		handler.ResponseError(ctx, http.StatusNotFound, "Activity not found")
		return
	}

	// Validate JSON payload using the improved validator
	body, err := ctx.GetRawData()
	if err != nil {
		handler.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Use the JSON validator for comprehensive validation
	schema := validator.GetActivityValidationSchema()
	if err := validator.ValidateJSON(body, schema); err != nil {
		handler.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Reset the body for normal binding
	ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	var request dto.ActivityUpdateRequest
	if handler.BindAndValidate(ctx, &request) {
		return
	}

	userID := ctx.GetString("user_id")
	response, err := c.activityService.UpdateActivity(activityID, userID, request)
	if err != nil {
		// If it's a GORM record not found error, return 404
		if strings.Contains(err.Error(), "record not found") {
			handler.ResponseError(ctx, http.StatusNotFound, "Activity not found")
			return
		}
		// If it's an activity type validation error, return 400
		if strings.Contains(err.Error(), "invalid activity type") {
			handler.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		// For other errors, return 500
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, response)
}

// DeleteActivity godoc
// @Summary Delete activity
// @Description Delete an existing activity by ID
// @Tags activities
// @Accept json
// @Produce json
// @Param activityId path string true "Activity ID"
// @Success 200 {object} map[string]interface{} "Activity deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request - Invalid activity ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Activity not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /activity/{activityId} [delete]
func (c ActivityController) DeleteActivity(ctx *gin.Context) {
	activityID := ctx.Param("activityId")
	if activityID == "" {
		handler.ResponseError(ctx, http.StatusBadRequest, "Activity ID is required")
		return
	}

	// Validate UUID format - return 404 for invalid format as it means "not found"
	if _, err := uuid.Parse(activityID); err != nil {
		handler.ResponseError(ctx, http.StatusNotFound, "Activity not found")
		return
	}

	userID := ctx.GetString("user_id")
	err := c.activityService.DeleteActivity(activityID, userID)
	if err != nil {
		// If it's a GORM record not found error, return 404
		if strings.Contains(err.Error(), "record not found") {
			handler.ResponseError(ctx, http.StatusNotFound, "Activity not found")
			return
		}
		// For other errors, return 500
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}
