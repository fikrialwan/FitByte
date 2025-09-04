package service

import (
	"time"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/entity"
	"github.com/fikrialwan/FitByte/internal/repository"
	"github.com/fikrialwan/FitByte/pkg/helpers"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository repository.UserRepository
	jwtService     JwtService
	cacheService   CacheService
}

func NewUserService(userRepository repository.UserRepository, jwtService JwtService, cacheService CacheService) UserService {
	return UserService{
		userRepository: userRepository,
		jwtService:     jwtService,
		cacheService:   cacheService,
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
	existingUser, _ := s.userRepository.GetByEmail(email)
	if existingUser.ID != uuid.Nil {
		return dto.LoginRegisterResponse{}, dto.ErrUserEmailExist
	}

	// Generate UUID upfront for immediate use in JWT token
	userID := uuid.New()

	// Create new entity.User struct instead of modifying fetched one
	newUser := entity.User{
		ID:        userID,
		Email:     email,
		Password:  password,
		Timestamp: entity.Timestamp{CreatedAt: time.Now()},
	}

	err := s.userRepository.CreateUser(&newUser)
	if err != nil {
		return dto.LoginRegisterResponse{}, err
	}

	token := s.jwtService.GenerateAccessToken(userID.String())

	return dto.LoginRegisterResponse{
		Email: email,
		Token: token,
	}, nil
}

func (s UserService) GetProfile(userId string) (dto.UserResponse, error) {
	if profile, err := s.cacheService.GetUserProfile(userId); err == nil {
		return profile, nil
	}

	user, err := s.userRepository.GetById(userId)
	if err != nil {
		return dto.UserResponse{}, err
	}

	response := dto.UserResponse{
		Email:      user.Email,
		Name:       user.Name,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		ImageUri:   user.ImageUri,
	}

	// Cache the result for 5 minutes
	s.cacheService.SetUserProfile(userId, response, 5*time.Minute)

	return response, nil
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
