package router

import (
	"enomo/api/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(
	reg *domain.UserRegisterHandler,
	login *domain.UserLoginHandler,
	logout *domain.UserLogoutHandler,
	list *domain.UserGroupsHandler,
	rename *domain.UserRenameHandler,
	energy *domain.UserEnergyHandler,
	delete *domain.UserDeleteHandler,
	groupMake *domain.GroupMakeHandler,
	groupAdd *domain.GroupAddHandler,
	groupList *domain.GroupListHandler,
	groupDelete *domain.GroupDeleteHandler,
	groupLeave *domain.GroupLeaveHandler,
	reco *domain.RecommendationsHandler,
	groupSettings *domain.GroupSettingsHandler,
	
) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	api := e.Group("/api/v1")
	api.POST("/users/register", reg.Register) //
	api.POST("/users/login", login.Login)     //
	api.POST("/users/logout", logout.Logout)  //
	api.GET("/users/list", list.List)         //
	api.PATCH("/users/rename", rename.Update) //
	api.PATCH("/users/energy", energy.Update) //
	api.POST("/users/delete", delete.Delete)  //

	api.POST("/groups/make", groupMake.Make) //
	api.POST("/groups/add", groupAdd.Add)    //
	api.GET("/groups/list", groupList.List)
	api.POST("/groups/delete", groupDelete.Delete) //
	api.GET("/groups/settings", groupSettings.Get) //
	api.POST("/groups/settings/update", groupSettings.Update) //
	api.POST("/groups/leave", groupLeave.Leave)    //

	api.GET("/groups/:groupId/recommendations", reco.Get)

	return e
}
