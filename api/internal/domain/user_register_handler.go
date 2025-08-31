package domain

import (
	"net/http"

	"enomo/server/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserRegisterHandler struct {
	uc *usecase.UserRegisterUsecase
}

func NewUserRegisterHandler(uc *usecase.UserRegisterUsecase) *UserRegisterHandler {
	return &UserRegisterHandler{uc: uc}
}

type Register struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	EnergyValue int16  `json:"energy_value"`
}

func (h *UserRegisterHandler) Register(c echo.Context) error {
	var req Register
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	out, err := h.uc.Register(c.Request().Context(), usecase.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		EnergyValue: req.EnergyValue,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to register"})
	}
	return c.JSON(http.StatusCreated, out)
}
