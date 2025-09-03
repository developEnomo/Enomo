package domain

import (
	"net/http"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserRenameHandler struct {
	uc    *usecase.UserUpdateDisplayNameUsecase
	store repository.TokenStore
}

func NewUserRenameHandler(uc *usecase.UserUpdateDisplayNameUsecase, store repository.TokenStore) *UserRenameHandler {
	return &UserRenameHandler{uc: uc, store: store}
}

type RenameUser struct {
	UserID      *string `json:"user_id,omitempty"`
	DisplayName string  `json:"display_name"`
}

func (h *UserRenameHandler) Update(c echo.Context) error {
	var in RenameUser
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
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

	out, err := h.uc.Update(c.Request().Context(), usecase.UpdateDisplayNameRequest{
		UserID:      targetID,
		DisplayName: in.DisplayName,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
