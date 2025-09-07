package usecase

import (
	"context"
	"enomo/api/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email       string
	Password    string
	DisplayName string
	EnergyValue int16
}

type RegisterResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	EnergyValue int16     `json:"energy_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRegisterUsecase struct {
	repo *repository.UserRepository
}

func NewUserRegisterUsecase(r *repository.UserRepository) *UserRegisterUsecase {
	return &UserRegisterUsecase{repo: r}
}

func (u *UserRegisterUsecase) Register(ctx context.Context, in RegisterRequest) (RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{}, err
	}
	user := &repository.User{
		Email:       in.Email,
		Password:    string(hashedPassword),
		DisplayName: in.DisplayName,
		EnergyValue: in.EnergyValue,
	}

	if user.EnergyValue == 0 {
		user.EnergyValue = 3
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		EnergyValue: user.EnergyValue,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
