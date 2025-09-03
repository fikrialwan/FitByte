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

	ctx.JSON(http.StatusCreated, response)
}
