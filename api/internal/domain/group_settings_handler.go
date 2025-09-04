package domain

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"enomo/server/internal/usecase"
)

type GroupSettingsHandler struct{ UC usecase.GroupSettingsUsecase }

func NewGroupSettingsHandler(uc usecase.GroupSettingsUsecase) *GroupSettingsHandler {
	return &GroupSettingsHandler{UC: uc}
}

// 取得: GET /api/v1/groups/settings?group_id=...
func (h *GroupSettingsHandler) Get(c echo.Context) error {
	gid := c.QueryParam("group_id")
	if gid == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "group_id is required"})
	}
	hours, err := h.UC.GetHours(c.Request().Context(), gid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "failed_to_get_settings"})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"group_id":               gid,
		"refresh_interval_hours": hours,
		"preset":                 hoursToPreset(hours), // "12h" | "1d" | "2d" ... "7d"
	})
}

type updateReq struct {
	GroupID string `json:"group_id"`
	// どちらか一方を受け付ける: hours または preset
	Hours  *int   `json:"refresh_interval_hours"`
	Preset string `json:"preset"` // "12h","1d".."7d"
}

// 更新: POST /api/v1/groups/settings/update
func (h *GroupSettingsHandler) Update(c echo.Context) error {
	var req updateReq
	if err := c.Bind(&req); err != nil || req.GroupID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid_request"})
	}
	var hours int
	if req.Hours != nil {
		hours = *req.Hours
	} else if req.Preset != "" {
		p, err := presetToHours(req.Preset)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid_preset"})
		}
		hours = p
	} else {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "either refresh_interval_hours or preset is required"})
	}
	actor := "" // セッションがあれば取得
	hours, err := h.UC.UpdateHours(c.Request().Context(), req.GroupID, hours, actor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"group_id":               req.GroupID,
		"refresh_interval_hours": hours,
		"preset":                 hoursToPreset(hours),
		"updated":                true,
	})
}

func presetToHours(p string) (int, error) {
	p = strings.ToLower(strings.TrimSpace(p))
	switch p {
	case "12h":
		return 12, nil
	case "1d":
		return 24, nil
	case "2d":
		return 48, nil
	case "3d":
		return 72, nil
	case "4d":
		return 96, nil
	case "5d":
		return 120, nil
	case "6d":
		return 144, nil
	case "7d":
		return 168, nil
	default:
		// 数値文字列も許容（"24" など）
		if n, err := strconv.Atoi(p); err == nil {
			return n, nil
		}
		return 0, echo.NewHTTPError(http.StatusBadRequest, "unsupported preset")
	}
}

func hoursToPreset(h int) string {
	switch h {
	case 12:
		return "12h"
	case 24:
		return "1d"
	case 48:
		return "2d"
	case 72:
		return "3d"
	case 96:
		return "4d"
	case 120:
		return "5d"
	case 144:
		return "6d"
	case 168:
		return "7d"
	default:
		return strconv.Itoa(h) + "h"
	}
}
