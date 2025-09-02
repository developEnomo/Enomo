package domain

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"enomo/server/internal/usecase"
)

type RecommendationsHandler struct {
	UC usecase.RecommendationsUsecase
}

func NewRecommendationsHandler(uc usecase.RecommendationsUsecase) *RecommendationsHandler {
	return &RecommendationsHandler{UC: uc}
}

func (h *RecommendationsHandler) Get(c echo.Context) error {
	groupID := c.Param("groupId")

	valence := c.QueryParam("valence")
	if valence == "" {
		valence = "mid"
	}
	popBias, _ := strconv.ParseFloat(c.QueryParam("popularity_bias"), 64)
	if popBias == 0 {
		popBias = 0.7
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}

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
