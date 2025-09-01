package usecase

import (
	"context"
	"enomo/server/internal/repository"
	"errors"
	"strings"
	"time"
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
	repo *repository.GroupRepository
}

func NewGroupMakeUsecase(r *repository.GroupRepository) *GroupMakeUsecase {
	return &GroupMakeUsecase{repo: r}
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
	if err := u.repo.Create(ctx, g); err != nil {
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
