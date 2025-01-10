package config

import (
	"github.com/SkandarEverest/refresh-golang/internal/handler"
	"github.com/SkandarEverest/refresh-golang/internal/middleware"
	"github.com/SkandarEverest/refresh-golang/internal/routes"
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

	// user
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Config)
	userHandler := handler.NewUserHandler(userUseCase, config.Log, config.Validate)

	// file
	fileUseCase := usecase.NewFileUseCase(config.Log, config.Config)
	fileHandler := handler.NewFileHandler(fileUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.Auth(config.Config)

	// setup route
	routeConfig := routes.RouteConfig{
		App:            config.App,
		UserController: userHandler,
		FileController: fileHandler,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()

}
