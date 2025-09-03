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

func (r ActivityRepository) CreateActivity(activity entity.Acticity) (entity.Acticity, error) {
	result := r.db.Create(&activity)

	if result.Error != nil {
		return entity.Acticity{}, result.Error
	}

	return activity, nil
}
