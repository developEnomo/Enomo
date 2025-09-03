package router

import (
	"enomo/server/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(
	reg *domain.UserRegisterHandler,
	login *domain.UserLoginHandler,
	logout *domain.UserLogoutHandler,
	groupMake *domain.GroupMakeHandler,
	groupAdd *domain.GroupAddHandler,
	groupList *domain.GroupListHandler,
	groupDelete *domain.GroupDeleteHandler,
	settingsH *domain.GroupSettingsHandler,

) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	api := e.Group("/api/v1")
	api.POST("/users/register", reg.Register)
	api.POST("/users/login", login.Login)
	api.POST("/users/logout", logout.Logout)

	api.POST("/groups/make", groupMake.Make)
	api.POST("/groups/add", groupAdd.Add)
	api.GET("/groups/list", groupList.List)
	api.POST("/groups/delete", groupDelete.Delete)
	api.GET("/groups/settings", settingsH.Get)
	api.POST("/groups/settings/update", settingsH.Update)

	api.GET("/groups/settings", settingsH.Get)
	api.POST("/groups/settings/update", settingsH.Update)

	api.GET("/groups/:groupId/recommendations", reco.Get)
	return e
}
