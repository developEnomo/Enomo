package usecase

import (
	"context"
	"database/sql"
	"enomo/server/internal/repository"
	"errors"
	"time"
)

var ErrUnauthorized = errors.New("unauthorized")

func IsUnauthorized(err error) bool { return errors.Is(err, ErrUnauthorized) }

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	EnergyValue int16     `json:"energy_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserLoginUsecase struct {
	repo *repository.UserRepository
}

func NewUserLoginUsecase(r *repository.UserRepository) *UserLoginUsecase {
	return &UserLoginUsecase{repo: r}
}

func (u *UserLoginUsecase) Login(ctx context.Context, in LoginInput) (LoginOutput, error) {
	if in.Email == "" || in.Password == "" {
		return LoginOutput{}, ErrUnauthorized
	}

	user, err := u.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginOutput{}, ErrUnauthorized
		}
		return LoginOutput{}, err
	}
	if user.Password != in.Password {
		return LoginOutput{}, ErrUnauthorized
	}

	return LoginOutput{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		EnergyValue: user.EnergyValue,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
