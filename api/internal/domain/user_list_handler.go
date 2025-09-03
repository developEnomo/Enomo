package domain

import (
	"net/http"
	"strconv"

	"enomo/server/internal/repository"
	"enomo/server/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserGroupsHandler struct {
	uc    *usecase.UserGroupsUsecase
	store repository.TokenStore
}

func NewUserGroupsHandler(uc *usecase.UserGroupsUsecase, store repository.TokenStore) *UserGroupsHandler {
	return &UserGroupsHandler{uc: uc, store: store}
}

func (h *UserGroupsHandler) List(c echo.Context) error {
	uid := c.QueryParam("user_id")
	if uid == "" {
		ck, err := c.Cookie("session_token")
		if err != nil || ck.Value == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "no session"})
		}
		userID, ok, err := h.store.Get(ck.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "session error"})
		}
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
		}
		uid = userID
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	out, err := h.uc.List(c.Request().Context(), usecase.ListUserGroupsRequest{
		UserID: uid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
