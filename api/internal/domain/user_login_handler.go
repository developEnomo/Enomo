package domain

import (
	"net/http"
	"os"
	"time"

	"enomo/server/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
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

	token, err := newJWT(out.ID, 15*time.Minute)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "token issue failed"})
	}
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int((15 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   isHTTPS(c),
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, out)
}

var jwtKey = []byte(getenv("JWT_SECRET", "dev-secret"))

func newJWT(userID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtKey)
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

func getenv(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fb
}
