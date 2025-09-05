package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupSettingsHandler struct {
	UC    usecase.GroupSettingsUsecase
	store repository.TokenStore
}

func NewGroupSettingsHandler(uc usecase.GroupSettingsUsecase, store repository.TokenStore) *GroupSettingsHandler {
	return &GroupSettingsHandler{UC: uc, store: store}
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
	actor := ""
	if ck, err := c.Cookie("session_token"); err == nil && ck.Value != "" {
		if uid, ok, err := h.store.Get(ck.Value); err == nil && ok {
			actor = uid
		}
	}
	if actor == "" {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}

	hours, err := h.UC.UpdateHours(c.Request().Context(), req.GroupID, req.Hours, actor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"group_id": req.GroupID, "refresh_interval_hours": hours})
}
