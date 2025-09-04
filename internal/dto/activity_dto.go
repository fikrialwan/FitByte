package dto

import (
	"time"

	"github.com/fikrialwan/FitByte/internal/entity"
	"github.com/google/uuid"
)

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
		DoneAt            *time.Time           `json:"doneAt,omitempty" example:"2024-01-15T07:30:00Z"`
		DurationInMinutes *int                 `json:"durationInMinutes,omitempty" binding:"omitempty,numeric,min=1,max=1440" example:"30" minimum:"1" maximum:"1440"`
	}

	// ActivityResponse represents the response payload for activity operations
	ActivityResponse struct {
		ID                uuid.UUID           `json:"activityId" example:"123e4567-e89b-12d3-a456-426614174000"`
		ActivityType      entity.ActivityType `json:"activityType" example:"Running"`
		DoneAt            time.Time           `json:"doneAt" example:"2024-01-15T07:30:00Z"`
		DurationInMinutes int                 `json:"durationInMinutes" example:"30"`
		CaloriesBurned    int                 `json:"caloriesBurned" example:"300"`
		CreatedAt         time.Time           `json:"createdAt" example:"2024-01-15T10:30:00Z"`
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
