package dto

import "errors"

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
		Name       string `json:"name" binding:"min=2,max=60"`
		ImageUri   string `json:"imageUri" binding:"url"`
	}

	UserResponse struct {
		Preference string `json:"preference"`
		WeightUnit string `json:"weightUnit"`
		HeightUnit string `json:"heightUnit"`
		Weight     int    `json:"weight"`
		Height     int    `json:"height"`
		Name       string `json:"name"`
		ImageUri   string `json:"imageUri"`
	}
)
