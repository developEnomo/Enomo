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
	out, err := h.uc.Make(c.Request().Context(), usecase.GroupMakeRequest{
		Name:    req.Name,
		TrackID: req.TrackID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, out)
}
