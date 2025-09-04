package repository

import (
	"github.com/fikrialwan/FitByte/internal/entity"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return ActivityRepository{db}
}

func (r ActivityRepository) CreateActivity(activity entity.Activity) (entity.Activity, error) {
	result := r.db.Create(&activity)

	if result.Error != nil {
		return entity.Activity{}, result.Error
	}

	return activity, nil
}

func (r ActivityRepository) GetActivityByID(activityID, userID string) (entity.Activity, error) {
	var activity entity.Activity
	result := r.db.Where("id = ? AND user_id = ?", activityID, userID).First(&activity)
	
	if result.Error != nil {
		return entity.Activity{}, result.Error
	}
	
	return activity, nil
}

func (r ActivityRepository) UpdateActivity(activity entity.Activity) (entity.Activity, error) {
	result := r.db.Save(&activity)
	
	if result.Error != nil {
		return entity.Activity{}, result.Error
	}
	
	return activity, nil
}
