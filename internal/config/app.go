package config

import (
	"github.com/SkandarEverest/refresh-golang/internal/handler"
	"github.com/SkandarEverest/refresh-golang/internal/usecase"

	db "github.com/SkandarEverest/refresh-golang/db/sqlc"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       *db.SQLStore
	App      *echo.Echo
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log)

	// setup controller
	UserHandler := handler.NewUserHandler(userUseCase, config.Log, config.Validate)
	user := config.App.Group("/users")
	user.POST("/create", UserHandler.CreateUser)

}
