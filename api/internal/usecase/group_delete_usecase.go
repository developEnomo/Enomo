package usecase

import (
	"context"
	"errors"

	"enomo/api/internal/repository"
)

type GroupDeleteUsecase struct{ groups *repository.GroupRepository }

func NewGroupDeleteUsecase(g *repository.GroupRepository) *GroupDeleteUsecase {
	return &GroupDeleteUsecase{groups: g}
}

type DeleteGroupRequest struct {
	GroupID string
}

type DeleteGroupResponse struct {
	Deleted int64 `json:"deleted"` // 1なら削除済み、0なら見つからず
}

func (u *GroupDeleteUsecase) Delete(ctx context.Context, in DeleteGroupRequest) (DeleteGroupResponse, error) {
	if in.GroupID == "" {
		return DeleteGroupResponse{}, errors.New("group_id required")
	}
	n, err := u.groups.Delete(ctx, in.GroupID)
	if err != nil {
		return DeleteGroupResponse{}, err
	}
	return DeleteGroupResponse{Deleted: n}, nil
}
