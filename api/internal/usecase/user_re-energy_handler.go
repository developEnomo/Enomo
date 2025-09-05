package usecase

import (
	"context"
	"errors"

	"enomo/api/internal/repository"
)

type UpdateEnergyValueRequest struct {
	UserID      string
	EnergyValue int16
}

type UpdateEnergyValueResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	EnergyValue int16  `json:"energy_value"`
}

type UserUpdateEnergyValueUsecase struct{ repo *repository.UserRepository }

func NewUserUpdateEnergyValueUsecase(r *repository.UserRepository) *UserUpdateEnergyValueUsecase {
	return &UserUpdateEnergyValueUsecase{repo: r}
}

func (u *UserUpdateEnergyValueUsecase) Update(ctx context.Context, in UpdateEnergyValueRequest) (UpdateEnergyValueResponse, error) {
	if in.UserID == "" {
		return UpdateEnergyValueResponse{}, errors.New("user_id required")
	}
	if in.EnergyValue < 1 || in.EnergyValue > 4 {
		return UpdateEnergyValueResponse{}, errors.New("energy_value must be 1 ~ 4")
	}
	user, err := u.repo.UpdateEnergyValue(ctx, in.UserID, in.EnergyValue)
	if err != nil {
		return UpdateEnergyValueResponse{}, err
	}
	return UpdateEnergyValueResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		EnergyValue: user.EnergyValue,
	}, nil
}
