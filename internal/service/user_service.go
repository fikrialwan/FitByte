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

	err := s.userRepository.CreateUser(&user)
	if err != nil {
		return dto.LoginRegisterResponse{}, err
	}

	token := s.jwtService.GenerateAccessToken(user.ID.String())

	return dto.LoginRegisterResponse{
		Email: email,
		Token: token,
	}, nil
}

func (s UserService) GetProfile(userId string) (dto.UserResponse, error) {
	user, err := s.userRepository.GetById(userId)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Email:      user.Email,
		Name:       user.Name,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		ImageUri:   user.ImageUri,
	}, nil
}

func (s UserService) UpdateProfile(userId string, request dto.UserRequest) (dto.UserResponse, error) {
	user, err := request.ToUserEntity(userId)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if err = s.userRepository.Update(&user); err != nil {
		return dto.UserResponse{}, err
	}

	return dto.NewUserResponseFromEntity(user), nil
}
