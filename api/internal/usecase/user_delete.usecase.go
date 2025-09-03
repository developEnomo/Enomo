package usecase

import (
	"context"
	"errors"

	"enomo/server/internal/repository"
)

var ErrUserOwnsGroups = errors.New("user owns groups")

type DeleteUserRequest struct {
	UserID string
}

type DeleteUserResponse struct {
	Deleted int64 `json:"deleted"`
}

type UserDeleteUsecase struct {
	users  *repository.UserRepository
	groups *repository.GroupRepository
}

func NewUserDeleteUsecase(u *repository.UserRepository, g *repository.GroupRepository) *UserDeleteUsecase {
	return &UserDeleteUsecase{users: u, groups: g}
}

func (u *UserDeleteUsecase) Delete(ctx context.Context, in DeleteUserRequest) (DeleteUserResponse, error) {
	if in.UserID == "" {
		return DeleteUserResponse{}, errors.New("user_id required")
	}
	owns, err := u.groups.OwnsAny(ctx, in.UserID)
	if err != nil {
		return DeleteUserResponse{}, err
	}
	if owns {
		return DeleteUserResponse{}, ErrUserOwnsGroups
	}
	n, err := u.users.Delete(ctx, in.UserID)
	if err != nil {
		return DeleteUserResponse{}, err
	}
	return DeleteUserResponse{Deleted: n}, nil
}
