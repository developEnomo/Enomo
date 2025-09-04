package domain

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"enomo/api/internal/usecase"
)

type GroupSettingsHandler struct{ UC usecase.GroupSettingsUsecase }

func NewGroupSettingsHandler(uc usecase.GroupSettingsUsecase) *GroupSettingsHandler {
	return &GroupSettingsHandler{UC: uc}
}

// GET /api/v1/groups/settings?group_id=...
func (h *GroupSettingsHandler) Get(c echo.Context) error {
	groupID := c.QueryParam("group_id")
	if groupID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "group_id is required"})
	}
	hours, err := h.UC.GetHours(c.Request().Context(), groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "failed_to_get_settings"})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"group_id":               groupID,
		"refresh_interval_hours": hours,
	})
}

type updateReq struct {
	GroupID string `json:"group_id"`
	Hours   int    `json:"refresh_interval_hours"`
}

// POST /api/v1/groups/settings/update
func (h *GroupSettingsHandler) Update(c echo.Context) error {
	var req updateReq
	if err := c.Bind(&req); err != nil || req.GroupID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid_request"})
	}
	actor := "" // セッションがあれば userID を入れる
	hours, err := h.UC.UpdateHours(c.Request().Context(), req.GroupID, req.Hours, actor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"group_id":               req.GroupID,
		"refresh_interval_hours": hours,
		"updated":                true,
	})
}
