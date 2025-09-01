package repository

import (
	"context"
	"database/sql"
	"time"
)

type TokenStore interface {
	Set(token, userID string, ttl time.Duration) error
	Get(token string) (userID string, ok bool, err error)
	Delete(token string) error
	CleanupExpired(ctx context.Context) error
}

type PostgresTokenStore struct{ DB *sql.DB }

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore { return &PostgresTokenStore{DB: db} }

func (s *PostgresTokenStore) Set(token, userID string, ttl time.Duration) error {
	exp := time.Now().Add(ttl)
	_, err := s.DB.Exec(`
		INSERT INTO sessions (token, user_id, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, expires_at = EXCLUDED.expires_at
	`, token, userID, exp)
	return err
}

func (s *PostgresTokenStore) Get(token string) (string, bool, error) {
	var userID string
	var exp time.Time
	err := s.DB.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = $1`, token).Scan(&userID, &exp)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if time.Now().After(exp) {
		_ = s.Delete(token)
		return "", false, nil
	}
	return userID, true, nil
}

func (s *PostgresTokenStore) Delete(token string) error {
	_, err := s.DB.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func (s *PostgresTokenStore) CleanupExpired(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < now()`)
	return err
}
