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
