package router

import (
	"enomo/server/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(reg *domain.UserRegisterHandler, login *domain.UserLoginHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	api := e.Group("/api/v1")
	api.POST("/users/register", reg.Register)
	api.POST("/auth/login", login.Login)

	return e
}
