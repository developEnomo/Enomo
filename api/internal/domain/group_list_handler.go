package domain

import (
	"net/http"
	"strconv"

	"enomo/api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type GroupListHandler struct{ uc *usecase.GroupListUsecase }

func NewGroupListHandler(uc *usecase.GroupListUsecase) *GroupListHandler {
	return &GroupListHandler{uc: uc}
}

func (h *GroupListHandler) List(c echo.Context) error {
	groupID := c.QueryParam("group_id")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	out, err := h.uc.ListMembers(c.Request().Context(), usecase.ListMembersRequest{
		GroupID: groupID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
