package domain

import (
	"net/http"

	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupDeleteHandler struct {
	uc *usecase.GroupDeleteUsecase
}

func NewGroupDeleteHandler(uc *usecase.GroupDeleteUsecase) *GroupDeleteHandler {
	return &GroupDeleteHandler{uc: uc}
}

type DeleteGroup struct {
	GroupID string `json:"group_id"`
}

func (h *GroupDeleteHandler) Delete(c echo.Context) error {
	var in DeleteGroup
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}
	out, err := h.uc.Delete(c.Request().Context(), usecase.DeleteGroupRequest{GroupID: in.GroupID})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
