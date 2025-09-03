package usecase

import (
	"context"
	"errors"
	"unicode/utf8"

	"enomo/server/internal/repository"
)

type UpdateDisplayNameRequest struct {
	UserID      string
	DisplayName string
}

type UpdateDisplayNameResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	EnergyValue int16  `json:"energy_value"`
}

type UserUpdateDisplayNameUsecase struct{ repo *repository.UserRepository }

func NewUserUpdateDisplayNameUsecase(r *repository.UserRepository) *UserUpdateDisplayNameUsecase {
	return &UserUpdateDisplayNameUsecase{repo: r}
}

func (u *UserUpdateDisplayNameUsecase) Update(ctx context.Context, in UpdateDisplayNameRequest) (UpdateDisplayNameResponse, error) {
	if in.UserID == "" || in.DisplayName == "" {
		return UpdateDisplayNameResponse{}, errors.New("user_id and display_name required")
	}
	if utf8.RuneCountInString(in.DisplayName) > 50 {
		return UpdateDisplayNameResponse{}, errors.New("display_name too long")
	}

	user, err := u.repo.UpdateDisplayName(ctx, in.UserID, in.DisplayName)
	if err != nil {
		return UpdateDisplayNameResponse{}, err
	}

	return UpdateDisplayNameResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		EnergyValue: user.EnergyValue,
	}, nil
}
