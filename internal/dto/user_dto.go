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
)
