package domain

import (
	"net/http"

	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserLogoutHandler struct {
	uc *usecase.UserLogoutUsecase
}

func NewUserLogoutHandler(uc *usecase.UserLogoutUsecase) *UserLogoutHandler {
	return &UserLogoutHandler{uc: uc}
}

type Logout struct {
	UserID string `json:"user_id"`
}

func (h *UserLogoutHandler) Logout(c echo.Context) error {
	var req Logout
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id required"})
	}
	out, err := h.uc.LogoutByUser(c.Request().Context(), usecase.LogoutByUserRequest{UserID: req.UserID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "logout failed"})
	}

	// クライアントに session_token が残っていれば失効させる（任意・常にクリア）
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, out)
}
