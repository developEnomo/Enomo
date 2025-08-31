package router

import (
	"enomo/server/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(userHandler *domain.UserHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"status": "ok"})
	})

	api := e.Group("/api/v1")
	api.POST("/users/register", userHandler.Register)

	return e
}