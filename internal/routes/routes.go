package routes

import (
	"github.com/SkandarEverest/refresh-golang/internal/handler"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App            *echo.Echo
	UserController *handler.UserHandler
	AuthMiddleware echo.MiddlewareFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	auth := c.App.Group("/auth")
	auth.POST("/", c.UserController.CreateUser)
}

func (c *RouteConfig) SetupAuthRoute() {
	user := c.App.Group("/user", c.AuthMiddleware)
	user.GET("/", c.UserController.GetUser)
	user.PATCH("/", c.UserController.UpdateUser)
}
