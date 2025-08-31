package domain

import (
	"net/http"

	"enomo/server/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserLoginHandler struct {
	uc *usecase.UserLoginUsecase
}

func NewUserLoginHandler(uc *usecase.UserLoginUsecase) *UserLoginHandler {
	return &UserLoginHandler{uc: uc}
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserLoginHandler) Login(c echo.Context) error {
	var req Login
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	out, err := h.uc.Login(c.Request().Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if usecase.IsUnauthorized(err) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "login failed"})
	}
	return c.JSON(http.StatusOK, out)
}
