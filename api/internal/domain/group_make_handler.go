package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupMakeHandler struct {
	uc    *usecase.GroupMakeUsecase
	store repository.TokenStore
}

func NewGroupMakeHandler(uc *usecase.GroupMakeUsecase, store repository.TokenStore) *GroupMakeHandler {
	return &GroupMakeHandler{uc: uc, store: store}
}

type MakeGroup struct {
	Name    string  `json:"name"`
	TrackID *string `json:"track_id,omitempty"`
}

func (h *GroupMakeHandler) Make(c echo.Context) error {
	var req MakeGroup
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}

	ck, err := c.Cookie("session_token")
	if err != nil || ck.Value == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	actorID, ok, err := h.store.Get(ck.Value)
	if err != nil || !ok || actorID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
	}

	out, err := h.uc.Make(c.Request().Context(), usecase.GroupMakeRequest{
		Name:    req.Name,
		OwnerID: actorID, // ← ここで Cookie 由来のユーザーIDをセット
		TrackID: req.TrackID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, out)
}
