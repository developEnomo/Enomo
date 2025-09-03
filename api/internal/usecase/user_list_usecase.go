package usecase

import (
	"context"
	"errors"

	"enomo/server/internal/repository"
)

type UserGroupsUsecase struct {
	members *repository.GroupMemberRepository
}

func NewUserGroupsUsecase(m *repository.GroupMemberRepository) *UserGroupsUsecase {
	return &UserGroupsUsecase{members: m}
}

type ListUserGroupsRequest struct {
	UserID string
	Limit  int
	Offset int
}

type GroupWithCountDTO struct {
	GroupID     string `json:"group_id"`
	Name        string `json:"name"`
	MemberCount int    `json:"member_count"`
}

type ListUserGroupsResponse struct {
	UserID string              `json:"user_id"`
	Groups []GroupWithCountDTO `json:"groups"`
}

func (u *UserGroupsUsecase) List(ctx context.Context, in ListUserGroupsRequest) (ListUserGroupsResponse, error) {
	if in.UserID == "" {
		return ListUserGroupsResponse{}, errors.New("user_id required")
	}
	gcs, err := u.members.GroupList(ctx, in.UserID, in.Limit, in.Offset)
	if err != nil {
		return ListUserGroupsResponse{}, err
	}
	out := make([]GroupWithCountDTO, 0, len(gcs))
	for _, g := range gcs {
		out = append(out, GroupWithCountDTO{
			GroupID:     g.GroupID,
			Name:        g.GroupName,
			MemberCount: g.MemberCount,
		})
	}
	return ListUserGroupsResponse{UserID: in.UserID, Groups: out}, nil
}
