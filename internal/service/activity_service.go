package service

import (
	"fmt"
	"strings"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/entity"
	"github.com/fikrialwan/FitByte/internal/repository"
	"github.com/google/uuid"
)

type ActivityService struct {
	activityRepository repository.ActivityRepository
}

func NewActivityService(activityRepository repository.ActivityRepository) ActivityService {
	return ActivityService{activityRepository}
}

func (s ActivityService) GetActivity(filter dto.ActivityFilter) ([]dto.ActivityResponse, error) {
	activities, err := s.activityRepository.GetActivity(filter)
	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return []dto.ActivityResponse{}, nil
	}

	var responses []dto.ActivityResponse
	for _, activity := range activities {
		responses = append(responses, dto.ActivityResponse{
			ID:                activity.ID,
			ActivityType:      activity.ActivityType,
			DoneAt:            activity.DoneAt,
			DurationInMinutes: activity.DurationInMinutes,
			CaloriesBurned:    activity.CaloriesBurned,
			CreatedAt:         activity.CreatedAt,
		})
	}

	return responses, nil
}

func (s ActivityService) CreateActivity(activityReq dto.ActivityRequest, userId string) (dto.CreateActivityResponse, error) {
	if !activityReq.ActivityType.IsValid() {
		validTypes := entity.GetValidActivityTypeStrings()
		return dto.CreateActivityResponse{}, fmt.Errorf("invalid activity type '%s'. valid types: %s",
			activityReq.ActivityType, strings.Join(validTypes, ", "))
	}

	caloriesBurned := activityReq.ActivityType.CalculateBurnedCalories(activityReq.DurationInMinutes)

	activity := entity.Activity{
		ActivityType:      activityReq.ActivityType,
		DoneAt:            activityReq.DoneAt,
		DurationInMinutes: activityReq.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		UserID:            uuid.MustParse(userId),
	}

	createdActivity, err := s.activityRepository.CreateActivity(activity)
	if err != nil {
		return dto.CreateActivityResponse{}, err
	}

	return dto.CreateActivityResponse{
		ID:                createdActivity.ID,
		ActivityType:      createdActivity.ActivityType,
		DoneAt:            createdActivity.DoneAt,
		DurationInMinutes: createdActivity.DurationInMinutes,
		CaloriesBurned:    createdActivity.CaloriesBurned,
		CreatedAt:         createdActivity.CreatedAt,
		UpdatedAt:         createdActivity.UpdatedAt,
	}, nil
}

func (s ActivityService) UpdateActivity(activityID, userID string, updateReq dto.ActivityUpdateRequest) (dto.ActivityResponse, error) {
	// Get existing activity
	activity, err := s.activityRepository.GetActivityByID(activityID, userID)
	if err != nil {
		return dto.ActivityResponse{}, err
	}

	// Update fields if provided
	if updateReq.ActivityType != nil {
		if !updateReq.ActivityType.IsValid() {
			validTypes := entity.GetValidActivityTypeStrings()
			return dto.ActivityResponse{}, fmt.Errorf("invalid activity type '%s'. valid types: %s",
				*updateReq.ActivityType, strings.Join(validTypes, ", "))
		}
		activity.ActivityType = *updateReq.ActivityType
	}

	if updateReq.DoneAt != nil {
		activity.DoneAt = *updateReq.DoneAt
	}

	if updateReq.DurationInMinutes != nil {
		activity.DurationInMinutes = *updateReq.DurationInMinutes
	}

	// Recalculate calories based on current activity type and duration
	activity.CaloriesBurned = activity.ActivityType.CalculateBurnedCalories(activity.DurationInMinutes)

	// Update in database
	updatedActivity, err := s.activityRepository.UpdateActivity(activity)
	if err != nil {
		return dto.ActivityResponse{}, err
	}

	return dto.ActivityResponse{
		ID:                updatedActivity.ID,
		ActivityType:      updatedActivity.ActivityType,
		DoneAt:            updatedActivity.DoneAt,
		DurationInMinutes: updatedActivity.DurationInMinutes,
		CaloriesBurned:    updatedActivity.CaloriesBurned,
		CreatedAt:         updatedActivity.CreatedAt,
		UpdatedAt:         updatedActivity.UpdatedAt,
	}, nil
}
