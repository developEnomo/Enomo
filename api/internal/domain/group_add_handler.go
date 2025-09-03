package domain

import (
	"net/http"

	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupAddHandler struct {
	uc *usecase.GroupAddUsecase
}

func NewGroupAddHandler(uc *usecase.GroupAddUsecase) *GroupAddHandler {
	return &GroupAddHandler{uc: uc}
}

type AddMember struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	IsAdmin *bool  `json:"is_admin,omitempty"`
}

func (h *GroupAddHandler) Add(c echo.Context) error {
	var in AddMember
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	isAdmin := false
	if in.IsAdmin != nil {
		isAdmin = *in.IsAdmin
	}
	out, err := h.uc.AddMember(c.Request().Context(), usecase.AddMemberRequest{
		GroupID: in.GroupID,
		UserID:  in.UserID,
		IsAdmin: isAdmin,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, out)
}
