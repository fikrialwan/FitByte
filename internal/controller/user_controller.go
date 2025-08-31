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
