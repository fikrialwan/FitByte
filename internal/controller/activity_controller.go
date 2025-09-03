package controller

import (
	"log"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("DEBUG FILTER: %+v\n", filter)

	res, err := c.activityService.GetActivity(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
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

	ctx.JSON(http.StatusCreated, response)
}
