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

	FindActiveByUser(ctx context.Context, userID string) (token string, ok bool, err error)
	Renew(token string, ttl time.Duration) error
	DeleteByUser(ctx context.Context, userID string) (int64, error)
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

func (s *PostgresTokenStore) FindActiveByUser(ctx context.Context, userID string) (string, bool, error) {
	var token string
	var exp time.Time
	err := s.DB.QueryRowContext(ctx, `
		SELECT token, expires_at
		FROM sessions
		WHERE user_id = $1 AND expires_at > now()
		ORDER BY expires_at DESC
		LIMIT 1
	`, userID).Scan(&token, &exp)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return token, true, nil
}

func (s *PostgresTokenStore) Renew(token string, ttl time.Duration) error {
	_, err := s.DB.Exec(`
		UPDATE sessions
		SET expires_at = $2
		WHERE token = $1
	`, token, time.Now().Add(ttl))
	return err
}

func (s *PostgresTokenStore) DeleteByUser(ctx context.Context, userID string) (int64, error) {
	res, err := s.DB.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
