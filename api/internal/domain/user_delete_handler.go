package domain

import (
	"errors"
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserDeleteHandler struct {
	uc    *usecase.UserDeleteUsecase
	store repository.TokenStore
}

func NewUserDeleteHandler(uc *usecase.UserDeleteUsecase, store repository.TokenStore) *UserDeleteHandler {
	return &UserDeleteHandler{uc: uc, store: store}
}

type DeleteUser struct {
	UserID *string `json:"user_id,omitempty"`
}

func (h *UserDeleteHandler) Delete(c echo.Context) error {
	var in DeleteUser
	_ = c.Bind(&in) // user_id 省略可

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

	out, err := h.uc.Delete(c.Request().Context(), usecase.DeleteUserRequest{UserID: targetID})
	if err != nil {
		if errors.Is(err, usecase.ErrUserOwnsGroups) {
			return c.JSON(http.StatusConflict, echo.Map{"error": "owner must transfer or delete groups before account deletion"})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out) // {"deleted":1}
}
