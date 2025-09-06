package domain

import (
	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GroupNowHandler struct {
	uc    *usecase.GroupNowUsecase
	store repository.TokenStore
}

func NewGroupNowHandler(uc *usecase.GroupNowUsecase, store repository.TokenStore) *GroupNowHandler {
	return &GroupNowHandler{uc: uc, store: store}
}

type setBody struct {
	GroupID string `json:"group_id"`
}

func (h *GroupNowHandler) Set(c echo.Context) error {
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

	var in setBody
	if err := c.Bind(&in); err != nil || in.GroupID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "group_id required"})
	}
	if err := h.uc.Set(c.Request().Context(), usecase.GroupNowSetRequest{UserID: userID, GroupID: in.GroupID}); err != nil {
		if err == usecase.ErrNotGroupMember {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "not a group member"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to set"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *GroupNowHandler) Get(c echo.Context) error {
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

	out, err := h.uc.Get(c.Request().Context(), userID)
	if err != nil {
		if err == usecase.ErrActiveGroupNotSet {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "active group not set"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get"})
	}
	return c.JSON(http.StatusOK, out)
}

func (h *GroupNowHandler) Clear(c echo.Context) error {
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

	if err := h.uc.Clear(c.Request().Context(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to clear"})
	}
	return c.NoContent(http.StatusNoContent)
}
