package domain

import (
	"net/http"
	"strings"

	"enomo/api/internal/repository"
	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupAddHandler struct {
	uc    *usecase.GroupAddUsecase
	store repository.TokenStore
}

func NewGroupAddHandler(uc *usecase.GroupAddUsecase, store repository.TokenStore) *GroupAddHandler {
	return &GroupAddHandler{uc: uc, store: store}
}

type AddMember struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	IsAdmin *bool  `json:"is_admin,omitempty"`
}

func (h *GroupAddHandler) Add(c echo.Context) error {
	var in AddMember
	if err := c.Bind(&in); err != nil || in.GroupID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}

	ck, err := c.Cookie("session_token")
	if err != nil || ck.Value == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	actorID, ok, err := h.store.Get(ck.Value)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
	}

	targetID := actorID
	if uid := strings.TrimSpace(in.UserID); uid != "" && uid != actorID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "adding others requires admin"})
	}

	isAdmin := false
	if in.IsAdmin != nil {
		isAdmin = *in.IsAdmin
	}

	out, err := h.uc.AddMember(c.Request().Context(), usecase.AddMemberRequest{
		GroupID: in.GroupID,
		UserID:  targetID,
		IsAdmin: isAdmin,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
