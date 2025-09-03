package repository

import (
	"errors"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/entity"
	"gorm.io/gorm"
)

type (
	UserRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (r UserRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	result := r.db.Where("email=?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, dto.ErrUserNotFound
	} else if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (r UserRepository) CreateUser(user entity.User) error {
	result := r.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r UserRepository) GetById(userId string) (entity.User, error) {
	var user entity.User
	result := r.db.Where("id=?", userId).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, dto.ErrUserNotFound
	} else if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
