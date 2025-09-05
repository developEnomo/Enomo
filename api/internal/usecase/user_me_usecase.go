package usecase

import (
	"context"
	"database/sql"
	"enomo/api/internal/repository"
	"errors"
)

var ErrMeNotFound = errors.New("user not found")

type MeRequest struct {
	UserID string
}

type MeResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	EnergyValue int16  `json:"energy_value"`
}

type UserMeUsecase struct {
	repo *repository.UserRepository
}

func NewUserMeUsecase(r *repository.UserRepository) *UserMeUsecase {
	return &UserMeUsecase{repo: r}
}

func (u *UserMeUsecase) Get(ctx context.Context, in MeRequest) (MeResponse, error) {
	if in.UserID == "" {
		return MeResponse{}, errors.New("user_id required")
	}
	user, err := u.repo.FindByID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MeResponse{}, ErrMeNotFound
		}
		return MeResponse{}, err
	}
	return MeResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		EnergyValue: user.EnergyValue,
	}, nil
}
