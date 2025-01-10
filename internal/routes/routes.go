package routes

import (
	"github.com/SkandarEverest/refresh-golang/internal/handler"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App            *echo.Echo
	UserController *handler.UserHandler
	FileController *handler.FileHandler
	AuthMiddleware echo.MiddlewareFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	auth := c.App.Group("/auth")
	auth.POST("/", c.UserController.AuthenticateUser)
}

func (c *RouteConfig) SetupAuthRoute() {
	//user
	user := c.App.Group("/user", c.AuthMiddleware)
	user.GET("/", c.UserController.GetUser)
	user.PATCH("/", c.UserController.UpdateUser)

	//file
	file := c.App.Group("/file", c.AuthMiddleware)
	file.POST("/upload", c.FileController.UploadFile)
}
