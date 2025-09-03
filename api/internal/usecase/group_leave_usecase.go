package usecase

import (
	"context"
	"errors"

	"enomo/server/internal/repository"
)

type GroupLeaveUsecase struct {
	members *repository.GroupMemberRepository
	groups  *repository.GroupRepository
}

func NewGroupLeaveUsecase(m *repository.GroupMemberRepository, g *repository.GroupRepository) *GroupLeaveUsecase {
	return &GroupLeaveUsecase{members: m, groups: g}
}

type LeaveMemberRequest struct {
	GroupID      string
	TargetUserID string
	ActorUserID  string
}

type LeaveMemberResponse struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	Deleted int64  `json:"deleted"`
	Action  string `json:"action"` // "left" or "kicked"
}

func (u *GroupLeaveUsecase) Leave(ctx context.Context, in LeaveMemberRequest) (LeaveMemberResponse, error) {
	if in.GroupID == "" || in.TargetUserID == "" || in.ActorUserID == "" {
		return LeaveMemberResponse{}, errors.New("group_id, target_user_id, actor_user_id required")
	}
	ownerID, err := u.groups.OwnerID(ctx, in.GroupID)
	if err != nil {
		return LeaveMemberResponse{}, err
	}

	action := "left"
	if in.TargetUserID != in.ActorUserID {
		isMember, isAdmin, err := u.members.IsMember(ctx, in.GroupID, in.ActorUserID)
		if err != nil {
			return LeaveMemberResponse{}, err
		}
		if !(in.ActorUserID == ownerID || (isMember && isAdmin)) {
			return LeaveMemberResponse{}, errors.New("forbidden")
		}
		action = "kicked"
	}

	// オーナー本人の退会は制限したい場合はここで弾く
	if in.TargetUserID == ownerID {
		return LeaveMemberResponse{}, errors.New("owner cannot leave; delete or transfer")
	}

	n, err := u.members.Remove(ctx, in.GroupID, in.TargetUserID)
	if err != nil {
		return LeaveMemberResponse{}, err
	}
	return LeaveMemberResponse{
		GroupID: in.GroupID,
		UserID:  in.TargetUserID,
		Deleted: n,
		Action:  action,
	}, nil
}
