package repository

import (
	"context"
	"database/sql"
	"time"

	pq "github.com/lib/pq"
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

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO users (email, password, display_name, energy_value)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, u.Email, u.Password, u.DisplayName, u.EnergyValue).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, email, password, display_name, energy_value, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.Password, &u.DisplayName, &u.EnergyValue, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByIDs(ctx context.Context, ids []string) ([]*User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, email, display_name, energy_value, created_at, updated_at
		FROM users
		WHERE id = ANY($1)
	`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.DisplayName, &u.EnergyValue, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &u)
	}
	return out, rows.Err()
}

func (r *UserRepository) UpdateDisplayName(ctx context.Context, userID, displayName string) (*User, error) {
	var u User
	err := r.DB.QueryRowContext(ctx, `
		UPDATE users
		   SET display_name = $2,
		       updated_at   = now()
		WHERE id = $1
		RETURNING id, email, password, display_name, energy_value, created_at, updated_at
	`, userID, displayName).
		Scan(&u.ID, &u.Email, &u.Password, &u.DisplayName, &u.EnergyValue, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
