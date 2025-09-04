package controller

import (
	"net/http"
	"strings"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/pkg/handler"
	"github.com/gin-gonic/gin"
)

type ActivityController struct {
	activityService service.ActivityService
}

func NewActivityController(activityService service.ActivityService) ActivityController {
	return ActivityController{activityService}
}

// CreateActivity godoc
// @Summary Create activity
// @Description Create a new activity with automatic calorie calculation
// @Tags activities
// @Accept json
// @Produce json
// @Param request body dto.ActivityRequest true "Activity data"
// @Success 201 {object} dto.ActivityResponse "Activity created successfully"
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
	activityID := ctx.Param("activityId")
	if activityID == "" {
		handler.ResponseError(ctx, http.StatusBadRequest, "Activity ID is required")
		return
	}

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
