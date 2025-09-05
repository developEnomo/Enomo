package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserLogoutHandler struct {
	uc    *usecase.UserLogoutUsecase
	store repository.TokenStore
}

func NewUserLogoutHandler(uc *usecase.UserLogoutUsecase, store repository.TokenStore) *UserLogoutHandler {
	return &UserLogoutHandler{uc: uc, store: store}
}

type Logout struct {
	UserID string `json:"user_id,omitempty"`
}

func (h *UserLogoutHandler) Logout(c echo.Context) error {
	var in Logout
	_ = c.Bind(&in)

	actorID := ""
	if ck, err := c.Cookie("session_token"); err == nil && ck.Value != "" {
		if uid, ok, err := h.store.Get(ck.Value); err == nil && ok {
			actorID = uid
		}
		c.SetCookie(&http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	}
	targetID := in.UserID
	if targetID == "" {
		targetID = actorID
	}
	if targetID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if actorID != "" && targetID != actorID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
	}

	out, err := h.uc.LogoutByUser(c.Request().Context(), usecase.LogoutByUserRequest{UserID: targetID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "logout failed"})
	}
	return c.JSON(http.StatusOK, out)
}
