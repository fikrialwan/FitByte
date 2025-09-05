package controller

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/pkg/handler"
	"github.com/fikrialwan/FitByte/pkg/validator"
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
	if handler.BindAndValidate(ctx, &request) {
		return
	}

	response, err := c.userService.Verify(request.Email, request.Password)
	if errors.Is(err, dto.ErrUserNotFound) {
		handler.ResponseError(ctx, http.StatusNotFound, "Invalid email or password")
		return
	} else if err != nil {
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, response)
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
	if handler.BindAndValidate(ctx, &request) {
		return
	}

	response, err := c.userService.Register(request.Email, request.Password)
	if errors.Is(err, dto.ErrUserEmailExist) {
		handler.ResponseError(ctx, http.StatusConflict, "Email already exists")
		return
	} else if err != nil {
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusCreated, response)
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
		handler.ResponseError(ctx, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, response)
}

// Register godoc
// @Summary Update user profile
// @Description Update user detail profile by id
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UserRequest true "profile data"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} utils.FailedResponse
// @Failure 500 {object} utils.FailedResponse
// @Router /user [patch]
func (c UserController) UpdateProfile(ctx *gin.Context) {
	// Check content type first
	contentType := ctx.GetHeader("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "application/json") {
		handler.ResponseError(ctx, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}

	userId := ctx.GetString("user_id")

	// Validate JSON payload using the improved validator
	body, err := ctx.GetRawData()
	if err != nil {
		handler.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Use the JSON validator for comprehensive validation
	schema := validator.GetUserValidationSchema()
	if err := validator.ValidateJSON(body, schema); err != nil {
		handler.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Reset the body for normal binding
	ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	var request dto.UserRequest
	if handler.BindAndValidate(ctx, &request) {
		return
	}

	response, err := c.userService.UpdateProfile(userId, request)
	if err != nil {
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, response)
}
