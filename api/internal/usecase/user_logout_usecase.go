package usecase

import (
	"context"

	"enomo/server/internal/repository"
)

type UserLogoutUsecase struct {
	store repository.TokenStore
}

func NewUserLogoutUsecase(store repository.TokenStore) *UserLogoutUsecase {
	return &UserLogoutUsecase{store: store}
}

type LogoutByUserRequest struct {
	UserID string
}

type LogoutByUserResponse struct {
	Deleted int64 `json:"deleted"`
}

func (u *UserLogoutUsecase) LogoutByUser(ctx context.Context, in LogoutByUserRequest) (LogoutByUserResponse, error) {
	if in.UserID == "" {
		return LogoutByUserResponse{Deleted: 0}, nil
	}
	n, err := u.store.DeleteByUser(ctx, in.UserID)
	if err != nil {
		return LogoutByUserResponse{}, err
	}
	return LogoutByUserResponse{Deleted: n}, nil
}
