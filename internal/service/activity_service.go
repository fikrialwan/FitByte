package service

import (
	"errors"

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

func (s ActivityService) CreateActivity(activityReq dto.ActivityRequest, userId string) (dto.ActivityResponse, error) {
	CaloriesPerMinute := activityReq.ActivityType.CaloriesPerMinute()
	if CaloriesPerMinute == 0 {
		return dto.ActivityResponse{}, errors.New("invalid activity type")
	}

	activity := entity.Acticity{
		ActivityType:      activityReq.ActivityType,
		DoneAt:            activityReq.DoneAt,
		DurationInMinutes: activityReq.DurationInMinutes,
		CaloriesBurned:    CaloriesPerMinute * activityReq.DurationInMinutes,
		UserID:            uuid.MustParse(userId),
	}

	createdActivity, err := s.activityRepository.CreateActivity(activity)
	if err != nil {
		return dto.ActivityResponse{}, err
	}

	return dto.ActivityResponse{
		ID:                createdActivity.ID,
		ActivityType:      createdActivity.ActivityType,
		DoneAt:            createdActivity.DoneAt,
		DurationInMinutes: createdActivity.DurationInMinutes,
		CaloriesBurned:    createdActivity.CaloriesBurned,
		CreatedAt:         createdActivity.CreatedAt,
		UpdatedAt:         createdActivity.UpdatedAt,
	}, nil
}
