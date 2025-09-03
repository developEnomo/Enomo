package usecase

import (
	"context"
	"errors"

	"enomo/api/internal/repository"
)

type GroupAddUsecase struct {
	members *repository.GroupMemberRepository
}

func NewGroupAddUsecase(members *repository.GroupMemberRepository) *GroupAddUsecase {
	return &GroupAddUsecase{members: members}
}

type AddMemberRequest struct {
	GroupID string
	UserID  string
	IsAdmin bool
}

type AddMemberResponse struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

func (u *GroupAddUsecase) AddMember(ctx context.Context, in AddMemberRequest) (AddMemberResponse, error) {
	if in.GroupID == "" || in.UserID == "" {
		return AddMemberResponse{}, errors.New("group_id and user_id required")
	}
	if err := u.members.Add(ctx, in.GroupID, in.UserID, in.IsAdmin); err != nil {
		return AddMemberResponse{}, err
	}
	return AddMemberResponse{GroupID: in.GroupID, UserID: in.UserID}, nil
}
