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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	response, err := c.userService.Verify(request.Email, request.Password)
	if errors.Is(err, dto.ErrUserNotFound) {
		// TODO: add error message
		ctx.AbortWithStatusJSON(http.StatusNotFound, nil)
	} else if err != nil {
		// TODO: add error message
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
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
// @Success 200 {object} dto.LoginRegisterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func (c UserController) Register(ctx *gin.Context) {
	var request dto.LoginRegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		// TODO: add error message
		ctx.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	response, err := c.userService.Register(request.Email, request.Password)
	if errors.Is(err, dto.ErrUserEmailExist) {
		// TODO: add error message
		ctx.AbortWithStatusJSON(http.StatusConflict, nil)
	} else if err != nil {
		// TODO: add error message
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
	}

	ctx.JSON(http.StatusOK, response)
}
