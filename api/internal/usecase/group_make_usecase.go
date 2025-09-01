package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"enomo/server/internal/repository"
)

type GroupMakeRequest struct {
	Name    string
	OwnerID string
	TrackID *string
}

type GroupMakeResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	TrackID   *string   `json:"track_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMakeUsecase struct {
	groups  *repository.GroupRepository
	members *repository.GroupMemberRepository
}

func NewGroupMakeUsecase(g *repository.GroupRepository, m *repository.GroupMemberRepository) *GroupMakeUsecase {
	return &GroupMakeUsecase{groups: g, members: m}
}

func (u *GroupMakeUsecase) Make(ctx context.Context, in GroupMakeRequest) (GroupMakeResponse, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || len(name) > 50 {
		return GroupMakeResponse{}, errors.New("invalid name")
	}
	if strings.TrimSpace(in.OwnerID) == "" {
		return GroupMakeResponse{}, errors.New("owner_id required")
	}

	g := &repository.Group{
		Name:    name,
		OwnerID: in.OwnerID,
		TrackID: in.TrackID,
	}
	if err := u.groups.Create(ctx, g); err != nil {
		return GroupMakeResponse{}, err
	}

	if err := u.members.Add(ctx, g.ID, in.OwnerID, true); err != nil {
		return GroupMakeResponse{}, err
	}

	return GroupMakeResponse{
		ID:        g.ID,
		Name:      g.Name,
		OwnerID:   g.OwnerID,
		TrackID:   g.TrackID,
		CreatedAt: g.CreatedAt,
	}, nil
}
