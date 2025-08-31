package repository

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID          string
	Email       string
	Password    string
	DisplayName string
	EnergyValue int16
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRepository struct{ DB *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO users (email, password, display_name, energy_value)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, u.Email, u.Password, u.DisplayName, u.EnergyValue).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}