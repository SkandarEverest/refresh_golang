package handler

import (
	"net/http"

	"github.com/SkandarEverest/refresh-golang/internal/dto"
	"github.com/SkandarEverest/refresh-golang/internal/exception"
	"github.com/SkandarEverest/refresh-golang/internal/usecase"

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
	var request = new(dto.UserRequest)

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("request doesn’t pass validation"))
	}

	if err := c.Validate.Struct(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, exception.BadRequest("request doesn’t pass validation"))
	}

	_, err := c.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		ex, ok := err.(*exception.CustomError)
		if ok {
			c.Log.Warnf("Failed to login user : %+v", err)
			return ctx.JSON(ex.StatusCode, ex)
		}
		panic(err)
	}

	return ctx.JSON(http.StatusOK, &dto.UserResponse{
		AccessToken: "token",
	})
}
