package usecase

import (
	"context"
	"database/sql"
	"errors"

	"enomo/api/internal/repository"
)

var (
	ErrActiveGroupNotSet = errors.New("active group not set")
	ErrNotGroupMember    = errors.New("user is not a member of the group")
)

type GroupNowSetRequest struct {
	UserID  string
	GroupID string
}
type GroupNowResponse struct {
	GroupID string `json:"group_id"`
}

type GroupNowUsecase struct {
	activeRepo *repository.ActiveGroupRepository
	memberRepo *repository.GroupMemberRepository
}

func NewGroupNowUsecase(active *repository.ActiveGroupRepository, member *repository.GroupMemberRepository) *GroupNowUsecase {
	return &GroupNowUsecase{activeRepo: active, memberRepo: member}
}

func (u *GroupNowUsecase) Set(ctx context.Context, in GroupNowSetRequest) error {
	ok, err := u.memberRepo.ExistsByUserAndGroup(ctx, in.UserID, in.GroupID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotGroupMember
	}
	return u.activeRepo.Set(ctx, in.UserID, in.GroupID)
}

func (u *GroupNowUsecase) Get(ctx context.Context, userID string) (GroupNowResponse, error) {
	gid, err := u.activeRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GroupNowResponse{}, ErrActiveGroupNotSet
		}
		return GroupNowResponse{}, err
	}
	return GroupNowResponse{GroupID: gid}, nil
}

func (u *GroupNowUsecase) Clear(ctx context.Context, userID string) error {
	return u.activeRepo.Clear(ctx, userID)
}
