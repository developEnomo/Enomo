package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserEnergyHandler struct {
	uc    *usecase.UserUpdateEnergyValueUsecase
	store repository.TokenStore
}

func NewUserEnergyHandler(uc *usecase.UserUpdateEnergyValueUsecase, store repository.TokenStore) *UserEnergyHandler {
	return &UserEnergyHandler{uc: uc, store: store}
}

type UpdateEnergy struct {
	UserID      *string `json:"user_id,omitempty"`
	EnergyValue *int    `json:"energy_value,omitempty"`
}

func (h *UserEnergyHandler) Update(c echo.Context) error {
	var in UpdateEnergy
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	if in.EnergyValue == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "energy_value required"})
	}
	if *in.EnergyValue <= 0 || *in.EnergyValue > 5 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "energy_value must be between 1 and 32767"})
	}

	var actorID string
	if ck, err := c.Cookie("session_token"); err == nil && ck.Value != "" {
		if uid, ok, err := h.store.Get(ck.Value); err == nil && ok {
			actorID = uid
		}
	}
	if actorID == "" && (in.UserID == nil || *in.UserID == "") {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "no session and no user_id"})
	}

	targetID := actorID
	if in.UserID != nil && *in.UserID != "" {
		targetID = *in.UserID
	}
	if actorID != "" && targetID != actorID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
	}

	ev := int16(*in.EnergyValue)
	out, err := h.uc.Update(c.Request().Context(), usecase.UpdateEnergyValueRequest{
		UserID:      targetID,
		EnergyValue: ev,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
