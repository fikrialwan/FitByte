package dto

import (
	"errors"

	"github.com/fikrialwan/FitByte/internal/entity"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserEmailExist = errors.New("user email exists")
)

type (
	LoginRegisterRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=32"`
	}

	LoginRegisterResponse struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	UserRequest struct {
		Preference string `json:"preference" binding:"required,oneof=CARDIO WEIGHT"`
		WeightUnit string `json:"weightUnit" binding:"required,oneof=KG LBS"`
		HeightUnit string `json:"heightUnit" binding:"required,oneof=CM INCH"`
		Weight     int    `json:"weight" binding:"required,min=10,max=1000"`
		Height     int    `json:"height" binding:"required,min=3,max=250"`
		Name       string `json:"name" binding:"omitempty,min=2,max=60"`
		ImageUri   string `json:"imageUri" binding:"omitempty,url"`
	}

	UserResponse struct {
		Email      string `json:"email"`
		Preference string `json:"preference"`
		WeightUnit string `json:"weightUnit"`
		HeightUnit string `json:"heightUnit"`
		Weight     int    `json:"weight"`
		Height     int    `json:"height"`
		Name       string `json:"name"`
		ImageUri   string `json:"imageUri"`
	}
)

func (req UserRequest) ToUserEntity(userIdStr string) (entity.User, error) {
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{
		ID:         userId,
		Name:       req.Name,
		Preference: req.Preference,
		WeightUnit: req.WeightUnit,
		HeightUnit: req.HeightUnit,
		Weight:     req.Weight,
		Height:     req.Height,
		ImageUri:   req.ImageUri,
	}, nil
}

func NewUserResponseFromEntity(user entity.User) UserResponse {
	return UserResponse{
		Email:      user.Email,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		Name:       user.Name,
		ImageUri:   user.ImageUri,
	}
}
