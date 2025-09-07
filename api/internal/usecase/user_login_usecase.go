package usecase

import (
	"context"
	"database/sql"
	"enomo/api/internal/repository"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("unauthorized")

func IsUnauthorized(err error) bool { return errors.Is(err, ErrUnauthorized) }

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
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

func (u *UserLoginUsecase) Login(ctx context.Context, in LoginRequest) (LoginResponse, error) {
	if in.Email == "" || in.Password == "" {
		return LoginResponse{}, ErrUnauthorized
	}

	user, err := u.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResponse{}, ErrUnauthorized
		}
		return LoginResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return LoginResponse{}, ErrUnauthorized
	}

	return LoginResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		EnergyValue: user.EnergyValue,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
