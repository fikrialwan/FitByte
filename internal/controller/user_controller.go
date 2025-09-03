package controller

import (
	"errors"
	"net/http"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRegisterRequest true "Login credentials"
// @Success 200 {object} dto.LoginRegisterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func (c UserController) Login(ctx *gin.Context) {
	var request dto.LoginRegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := c.userService.Verify(request.Email, request.Password)
	if errors.Is(err, dto.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary User registration
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRegisterRequest true "Registration credentials"
// @Success 201 {object} dto.LoginRegisterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func (c UserController) Register(ctx *gin.Context) {
	var request dto.LoginRegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := c.userService.Register(request.Email, request.Password)
	if errors.Is(err, dto.ErrUserEmailExist) {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user [get]
func (c UserController) GetProfile(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	response, err := c.userService.GetProfile(userId)
	if errors.Is(err, dto.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
