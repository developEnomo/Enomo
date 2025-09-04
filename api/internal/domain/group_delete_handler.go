package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupDeleteHandler struct {
	uc    *usecase.GroupDeleteUsecase
	store repository.TokenStore
}

func NewGroupDeleteHandler(uc *usecase.GroupDeleteUsecase, store repository.TokenStore) *GroupDeleteHandler {
	return &GroupDeleteHandler{uc: uc, store: store}
}

type DeleteGroup struct {
	GroupID string `json:"group_id"`
}

func (h *GroupDeleteHandler) Delete(c echo.Context) error {
	var in DeleteGroup
	if err := c.Bind(&in); err != nil || in.GroupID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	ck, err := c.Cookie("session_token")
	if err != nil || ck.Value == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	actorID, ok, err := h.store.Get(ck.Value)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
	}

	out, err := h.uc.Delete(c.Request().Context(), usecase.DeleteGroupRequest{GroupID: in.GroupID, ActorID: actorID})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
