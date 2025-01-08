package handler

import (
	"net/http"

	"github.com/SkandarEverest/refresh-golang/internal/dto"
	"github.com/SkandarEverest/refresh-golang/internal/exception"
	"github.com/SkandarEverest/refresh-golang/internal/usecase"
	"github.com/SkandarEverest/refresh-golang/pkg/helper"
	"github.com/SkandarEverest/refresh-golang/pkg/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Log      *logrus.Logger
	UseCase  *usecase.UserUseCase
	Validate *validator.Validate
}

func NewUserHandler(useCase *usecase.UserUseCase, logger *logrus.Logger, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *UserHandler) CreateUser(ctx echo.Context) error {
	var request = new(dto.AuthRequest)

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("request doesn’t pass validation"))
	}

	if err := c.Validate.Struct(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("request doesn’t pass validation"))
	}

	var token *string
	var err error

	if request.Action == "create" {
		token, err = c.UseCase.Create(ctx.Request().Context(), request)
	}

	if request.Action == "login" {
		token, err = c.UseCase.Login(ctx.Request().Context(), request)
	}

	if err != nil {
		ex, ok := err.(*exception.CustomError)
		if ok {
			c.Log.Warnf("Failed to authenticate: %+v", err)
			return ctx.JSON(ex.StatusCode, ex)
		}
		panic(err)
	}

	return ctx.JSON(http.StatusOK, &dto.AuthResponse{
		AccessToken: *token,
	})
}

func (c *UserHandler) GetUser(ctx echo.Context) error {

	userData := ctx.Get("user").(*jwt.JwtClaim)
	user, err := c.UseCase.GetUser(ctx.Request().Context(), userData.Id)
	if err != nil {
		ex, ok := err.(*exception.CustomError)
		if ok {
			c.Log.Warnf("Failed to login user : %+v", err)
			return ctx.JSON(ex.StatusCode, ex)
		}
		panic(err)
	}

	return ctx.JSON(http.StatusOK, &dto.UserResponse{
		Email:           user.Email,
		Username:        helper.DerefString(user.Username, ""),
		UserImageUri:    helper.DerefString(user.UserImageUri, ""),
		CompanyName:     helper.DerefString(user.CompanyName, ""),
		CompanyImageUri: helper.DerefString(user.CompanyImageUri, ""),
	})
}

func (c *UserHandler) UpdateUser(ctx echo.Context) error {
	var request = new(dto.UpdateUserRequest)

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("request doesn’t pass validation"))
	}

	if err := c.Validate.Struct(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("request doesn’t pass validation"))
	}

	userData := ctx.Get("user").(*jwt.JwtClaim)
	user, err := c.UseCase.UpdateUser(ctx.Request().Context(), request, userData.Id)
	if err != nil {
		ex, ok := err.(*exception.CustomError)
		if ok {
			c.Log.Warnf("Failed to login user : %+v", err)
			return ctx.JSON(ex.StatusCode, ex)
		}
		panic(err)
	}

	return ctx.JSON(http.StatusOK, &dto.UserResponse{
		Email:           user.Email,
		Username:        helper.DerefString(user.Username, ""),
		UserImageUri:    helper.DerefString(user.UserImageUri, ""),
		CompanyName:     helper.DerefString(user.CompanyName, ""),
		CompanyImageUri: helper.DerefString(user.CompanyImageUri, ""),
	})
}
