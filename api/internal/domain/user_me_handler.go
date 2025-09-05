package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserMeHandler struct {
	uc    *usecase.UserMeUsecase
	store repository.TokenStore
}

func NewUserMeHandler(uc *usecase.UserMeUsecase, store repository.TokenStore) *UserMeHandler {
	return &UserMeHandler{uc: uc, store: store}
}

func (h *UserMeHandler) Me(c echo.Context) error {
	ck, err := c.Cookie("session_token")
	if err != nil || ck.Value == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, ok, err := h.store.Get(ck.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "session error"})
	}
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
	}

	out, err := h.uc.Get(c.Request().Context(), usecase.MeRequest{UserID: userID})
	if err != nil {
		if err == usecase.ErrMeNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch user"})
	}
	return c.JSON(http.StatusOK, out)
}
