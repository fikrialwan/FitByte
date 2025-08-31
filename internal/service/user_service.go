package service

import (
	"time"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/repository"
	"github.com/fikrialwan/FitByte/pkg/helpers"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository repository.UserRepository
	jwtService     JwtService
}

func NewUserService(userRepository repository.UserRepository, jwtService JwtService) UserService {
	return UserService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s UserService) Verify(email, password string) (dto.LoginRegisterResponse, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return dto.LoginRegisterResponse{}, err
	}

	validPass, err := helpers.CheckPassword(user.Password, []byte(password))
	if err != nil || !validPass {
		return dto.LoginRegisterResponse{}, dto.ErrUserNotFound
	}

	token := s.jwtService.GenerateAccessToken(user.ID.String())

	return dto.LoginRegisterResponse{
		Email: email,
		Token: token,
	}, nil
}

func (s UserService) Register(email, password string) (dto.LoginRegisterResponse, error) {
	user, _ := s.userRepository.GetByEmail(email)
	if user.ID != uuid.Nil {
		return dto.LoginRegisterResponse{}, dto.ErrUserEmailExist
	}

	user.Email = email
	user.Password = password
	user.CreatedAt = time.Now()

	err := s.userRepository.CreateUser(user)
	if err != nil {
		return dto.LoginRegisterResponse{}, err
	}

	token := s.jwtService.GenerateAccessToken(user.ID.String())

	return dto.LoginRegisterResponse{
		Email: email,
		Token: token,
	}, nil
}
