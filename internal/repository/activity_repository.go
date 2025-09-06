package repository

import (
	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/entity"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return ActivityRepository{db}
}

func (r ActivityRepository) GetActivity(filter dto.ActivityFilter, userID string) ([]entity.Activity, error) {
	var activities []entity.Activity
	query := r.db.Model(&entity.Activity{}).Where("user_id = ?", userID)

	if filter.ActivityType != "" {
		query = query.Where("activity_type = ?", filter.ActivityType)
	}
	if !filter.DoneAtFrom.IsZero() {
		query = query.Where("done_at >= ?", filter.DoneAtFrom)
	}
	if !filter.DoneAtTo.IsZero() {
		query = query.Where("done_at <= ?", filter.DoneAtTo)
	}
	if filter.CaloriesBurnedMin > 0 {
		query = query.Where("calories_burned >= ?", filter.CaloriesBurnedMin)
	}
	if filter.CaloriesBurnedMax > 0 {
		query = query.Where("calories_burned <= ?", filter.CaloriesBurnedMax)
	}

	// default pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 5
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	if err := query.Limit(limit).Offset(offset).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
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

func (r ActivityRepository) DeleteActivity(activityID, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", activityID, userID).Delete(&entity.Activity{})

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected (activity existed and was deleted)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
