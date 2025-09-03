package domain

import (
	"net/http"

	"enomo/server/internal/repository"
	"enomo/server/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupLeaveHandler struct {
	uc    *usecase.GroupLeaveUsecase
	store repository.TokenStore
}

func NewGroupLeaveHandler(uc *usecase.GroupLeaveUsecase, store repository.TokenStore) *GroupLeaveHandler {
	return &GroupLeaveHandler{uc: uc, store: store}
}

type LeaveGroup struct {
	GroupID string  `json:"group_id"`
	UserID  *string `json:"user_id,omitempty"` // 省略時は自分を対象にする
}

func (h *GroupLeaveHandler) Leave(c echo.Context) error {
	var in LeaveGroup
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}

	cookie, err := c.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "no session"})
	}
	actorID, ok, err := h.store.Get(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "session error"})
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
	}

	targetID := actorID
	if in.UserID != nil && *in.UserID != "" {
		targetID = *in.UserID
	}

	out, err := h.uc.Leave(c.Request().Context(), usecase.LeaveMemberRequest{
		GroupID:      in.GroupID,
		TargetUserID: targetID,
		ActorUserID:  actorID,
	})
	if err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
