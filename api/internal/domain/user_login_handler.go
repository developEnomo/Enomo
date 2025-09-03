package domain

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserLoginHandler struct {
	uc    *usecase.UserLoginUsecase
	store repository.TokenStore
}

func NewUserLoginHandler(uc *usecase.UserLoginUsecase, store repository.TokenStore) *UserLoginHandler {
	return &UserLoginHandler{uc: uc, store: store}
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newSessionToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func isHTTPS(c echo.Context) bool {
	if c.Scheme() == "https" {
		return true
	}
	if c.Request().Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}
	return false
}

func (h *UserLoginHandler) Login(c echo.Context) error {
	var req Login
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	out, err := h.uc.Login(c.Request().Context(), usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if usecase.IsUnauthorized(err) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "login failed"})
	}

	ctx := c.Request().Context()
	const ttl = 60 * time.Minute

	if exist, ok, _ := h.store.FindActiveByUser(ctx, out.ID); ok {
		_ = h.store.Renew(exist, ttl)
		c.SetCookie(&http.Cookie{
			Name:     "session_token",
			Value:    exist,
			Path:     "/",
			MaxAge:   int(ttl.Seconds()),
			HttpOnly: true,
			Secure:   isHTTPS(c),
			SameSite: http.SameSiteLaxMode,
		})
		return c.JSON(http.StatusOK, out)
	}

	token := newSessionToken()
	if err := h.store.Set(token, out.ID, ttl); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "session save failed"})
	}
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   isHTTPS(c),
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, out)
}
