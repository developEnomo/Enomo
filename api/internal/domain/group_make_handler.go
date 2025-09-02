package domain

import (
	"net/http"

	"enomo/server/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupMakeHandler struct {
	uc *usecase.GroupMakeUsecase
}

func NewGroupMakeHandler(uc *usecase.GroupMakeUsecase) *GroupMakeHandler {
	return &GroupMakeHandler{uc: uc}
}

type MakeGroup struct {
	Name    string  `json:"name"`
	OwnerID string  `json:"owner_id"`
	TrackID *string `json:"track_id,omitempty"`
}

func (h *GroupMakeHandler) Make(c echo.Context) error {
	var req MakeGroup
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	out, err := h.uc.Make(c.Request().Context(), usecase.GroupMakeRequest{
		Name:    req.Name,
		OwnerID: req.OwnerID,
		TrackID: req.TrackID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, out)
}
