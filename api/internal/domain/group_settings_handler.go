package domain

import (
	"net/http"
	"strings"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupSettingsHandler struct {
	UC    usecase.GroupSettingsUsecase
	store TokenStore
}

type TokenStore interface {
	Get(token string) (userID string, ok bool, err error)
}

func NewGroupSettingsHandler(uc usecase.GroupSettingsUsecase, store repository.TokenStore) *GroupSettingsHandler {
	return &GroupSettingsHandler{UC: uc, store: store}
}

// GET /api/v1/groups/settings?group_id=...
func (h *GroupSettingsHandler) Get(c echo.Context) error {
	groupID := c.QueryParam("group_id")
	if strings.TrimSpace(groupID) == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "group_id is required"})
	}
	out, err := h.UC.Get(c.Request().Context(), groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "failed_to_get_settings"})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"group_id":               out.GroupID,
		"group_name":             out.GroupName,
		"refresh_interval_hours": out.RefreshIntervalHours,
	})
}

type updateReq struct {
	GroupID              string  `json:"group_id"`
	GroupName            *string `json:"group_name,omitempty"`
	RefreshIntervalHours *int    `json:"refresh_interval_hours,omitempty"`
}

// POST /api/v1/groups/settings/update
func (h *GroupSettingsHandler) Update(c echo.Context) error {
	var req updateReq
	if err := c.Bind(&req); err != nil || strings.TrimSpace(req.GroupID) == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid_request"})
	}
	// どちらも指定なしは NG
	if req.GroupName == nil && req.RefreshIntervalHours == nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "no_fields_to_update"})
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

	out, err := h.UC.Update(c.Request().Context(), usecase.UpdateGroupSettingsInput{
		GroupID:              req.GroupID,
		GroupName:            req.GroupName,
		RefreshIntervalHours: req.RefreshIntervalHours,
	}, actor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"group_id":               out.GroupID,
		"group_name":             out.GroupName,
		"refresh_interval_hours": out.RefreshIntervalHours,
	})
}
