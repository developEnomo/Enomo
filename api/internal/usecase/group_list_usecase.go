package usecase

import (
	"context"
	"errors"
	"time"

	"enomo/server/internal/repository"
)

type GroupListUsecase struct {
	members *repository.GroupMemberRepository
	users   *repository.UserRepository
}

func NewGroupListUsecase(members *repository.GroupMemberRepository, users *repository.UserRepository) *GroupListUsecase {
	return &GroupListUsecase{members: members, users: users}
}

type ListMembersRequest struct {
	GroupID string
	Limit   int
	Offset  int
}

type GroupMemberDTO struct {
	GroupID     string    `json:"group_id"`
	UserID      string    `json:"user_id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	EnergyValue int16     `json:"energy_value"`
	IsAdmin     bool      `json:"is_admin"`
	JoinedAt    time.Time `json:"joined_at"`
}

type ListMembersResponse struct {
	GroupID string           `json:"group_id"`
	Members []GroupMemberDTO `json:"members"`
}

func (u *GroupListUsecase) ListMembers(ctx context.Context, in ListMembersRequest) (ListMembersResponse, error) {
	if in.GroupID == "" {
		return ListMembersResponse{}, errors.New("group_id required")
	}
	if in.Limit <= 0 || in.Limit > 50 {
		in.Limit = 20
	}
	if in.Offset < 0 {
		in.Offset = 0
	}

	ms, err := u.members.UserList(ctx, in.GroupID, in.Limit, in.Offset)
	if err != nil {
		return ListMembersResponse{}, err
	}

	ids := make([]string, 0, len(ms))
	for _, m := range ms {
		ids = append(ids, m.UserID)
	}
	us, err := u.users.FindByIDs(ctx, ids)
	if err != nil {
		return ListMembersResponse{}, err
	}

	byID := make(map[string]*repository.User, len(us))
	for _, x := range us {
		byID[x.ID] = x
	}

	out := make([]GroupMemberDTO, 0, len(ms))
	for _, m := range ms {
		if x, ok := byID[m.UserID]; ok {
			out = append(out, GroupMemberDTO{
				GroupID:     m.GroupID,
				UserID:      x.ID,
				DisplayName: x.DisplayName,
				Email:       x.Email,
				EnergyValue: x.EnergyValue,
				IsAdmin:     m.IsAdmin,
				JoinedAt:    m.JoinedAt,
			})
		}
	}
	return ListMembersResponse{GroupID: in.GroupID, Members: out}, nil
}
