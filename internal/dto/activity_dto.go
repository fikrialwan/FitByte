package dto

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/fikrialwan/FitByte/internal/entity"
	"github.com/google/uuid"
)

// PreciseTime wraps time.Time to preserve exact millisecond formatting
type PreciseTime struct {
	time.Time
	originalFormat string
}

// UnmarshalJSON preserves the original time format
func (pt *PreciseTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	timeStr := strings.Trim(string(data), `"`)
	pt.originalFormat = timeStr
	
	// Parse the time
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}
	pt.Time = t
	return nil
}

// MarshalJSON returns the original format if available
func (pt PreciseTime) MarshalJSON() ([]byte, error) {
	if pt.originalFormat != "" {
		return json.Marshal(pt.originalFormat)
	}
	return json.Marshal(pt.Time.Format(time.RFC3339))
}

type (
	ActivityFilter struct {
		Limit             int       `form:"limit"`
		Offset            int       `form:"offset"`
		ActivityType      string    `form:"activityType"`
		DoneAtFrom        time.Time `form:"doneAtFrom"`
		DoneAtTo          time.Time `form:"doneAtTo"`
		CaloriesBurnedMin int       `form:"caloriesBurnedMin"`
		CaloriesBurnedMax int       `form:"caloriesBurnedMax"`
	}

	// ActivityRequest represents the request payload for creating an activity
	ActivityRequest struct {
		ActivityType      entity.ActivityType `json:"activityType" binding:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope" example:"Running" enums:"Walking,Yoga,Stretching,Cycling,Swimming,Dancing,Hiking,Running,HIIT,JumpRope"`
		DoneAt            time.Time           `json:"doneAt" binding:"required" example:"2024-01-15T07:30:00Z"`
		DurationInMinutes int                 `json:"durationInMinutes" binding:"required,numeric,min=1,max=1440" example:"30" minimum:"1" maximum:"1440"`
	}

	// ActivityUpdateRequest represents the request payload for updating an activity
	ActivityUpdateRequest struct {
		ActivityType      *entity.ActivityType `json:"activityType,omitempty" binding:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope" example:"Running" enums:"Walking,Yoga,Stretching,Cycling,Swimming,Dancing,Hiking,Running,HIIT,JumpRope"`
		DoneAt            *PreciseTime         `json:"doneAt,omitempty" binding:"omitempty" swaggertype:"string" format:"date-time" example:"2024-01-15T07:30:00Z"`
		DurationInMinutes *int                 `json:"durationInMinutes,omitempty" binding:"omitempty,min=1,max=1440" example:"30" minimum:"1" maximum:"1440"`
	}

	// ActivityResponse represents the response payload for activity operations
	ActivityResponse struct {
		ID                uuid.UUID           `json:"activityId" example:"123e4567-e89b-12d3-a456-426614174000"`
		ActivityType      entity.ActivityType `json:"activityType" example:"Running"`
		DoneAt            PreciseTime         `json:"doneAt" swaggertype:"string" format:"date-time" example:"2024-01-15T07:30:00Z"`
		DurationInMinutes int                 `json:"durationInMinutes" example:"30"`
		CaloriesBurned    int                 `json:"caloriesBurned" example:"300"`
		CreatedAt         time.Time           `json:"createdAt" example:"2024-01-15T10:30:00Z"`
		UpdatedAt         time.Time           `json:"updatedAt" example:"2024-01-15T10:30:00Z"`
	}

	// CreateActivityResponse represents the response payload for activity operations
	CreateActivityResponse struct {
		ID                uuid.UUID           `json:"activityId" example:"123e4567-e89b-12d3-a456-426614174000"`
		ActivityType      entity.ActivityType `json:"activityType" example:"Running"`
		DoneAt            time.Time           `json:"doneAt" example:"2024-01-15T07:30:00Z"`
		DurationInMinutes int                 `json:"durationInMinutes" example:"30"`
		CaloriesBurned    int                 `json:"caloriesBurned" example:"300"`
		CreatedAt         time.Time           `json:"createdAt" example:"2024-01-15T10:30:00Z"`
		UpdatedAt         time.Time           `json:"updatedAt" example:"2024-01-15T10:30:00Z"`
	}
)
