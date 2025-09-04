package domain

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"enomo/api/internal/usecase"
)

type RecommendationsHandler struct {
	UC       usecase.RecommendationsUsecase        // 一覧取得
	RefreshU usecase.RecommendationsRefreshUsecase // 更新保存
}

func NewRecommendationsHandler(
	uc usecase.RecommendationsUsecase,
	ruc usecase.RecommendationsRefreshUsecase,
) *RecommendationsHandler {
	return &RecommendationsHandler{UC: uc, RefreshU: ruc}
}

// GET /api/v1/groups/:groupId/recommendations
func (h *RecommendationsHandler) Get(c echo.Context) error {
	groupID := c.Param("groupId")

	// --- 入力パース ---
	valence := c.QueryParam("valence")
	switch valence {
	case "low", "mid", "high":
		// ok
	default:
		valence = "mid"
	}

	popBias := 0.7
	if s := c.QueryParam("popularity_bias"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			if v < 0 {
				v = 0
			}
			if v > 1 {
				v = 1
			}
			popBias = v
		}
	}

	limit := 20
	if s := c.QueryParam("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			if v < 1 {
				v = 1
			}
			if v > 50 {
				v = 50
			}
			limit = v
		}
	}

	// --- 実行 ---
	out, err := h.UC.Execute(c.Request().Context(), usecase.RecommendationsInput{
		GroupID:        groupID,
		ValencePreset:  valence,
		PopularityBias: popBias,
		Market:         "JP",
		MinPopularity:  60,
		Limit:          limit,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"code":    "INTERNAL_ERROR",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, out)
}

// POST /api/v1/groups/:groupId/recommendations/refresh
func (h *RecommendationsHandler) Refresh(c echo.Context) error {
	groupID := c.Param("groupId")

	// --- 入力パース ---
	valence := c.QueryParam("valence")
	switch valence {
	case "low", "mid", "high":
		// ok
	default:
		valence = "mid"
	}

	popBias := 0.7
	if s := c.QueryParam("popularity_bias"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			if v < 0 {
				v = 0
			}
			if v > 1 {
				v = 1
			}
			popBias = v
		}
	}

	limit := 20
	if s := c.QueryParam("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			if v < 1 {
				v = 1
			}
			if v > 50 {
				v = 50
			}
			limit = v
		}
	}

	// --- 実行 ---
	track, err := h.RefreshU.Refresh(c.Request().Context(), usecase.RecommendationsInput{
		GroupID:        groupID,
		ValencePreset:  valence,
		PopularityBias: popBias,
		Market:         "JP",
		MinPopularity:  60,
		Limit:          limit,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"code":    "INTERNAL_ERROR",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"group_id": groupID,
		"track":    track,
		"message":  "groups.track_id updated",
	})
}
