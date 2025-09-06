package repository

import (
	"context"
	"database/sql"
)

type ActiveGroupRepository struct{ DB *sql.DB }

func NewActiveGroupRepository(db *sql.DB) *ActiveGroupRepository {
	return &ActiveGroupRepository{DB: db}
}

func (r *ActiveGroupRepository) Set(ctx context.Context, userID, groupID string) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO user_active_groups (user_id, group_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET group_id = EXCLUDED.group_id, updated_at = now()
	`, userID, groupID)
	return err
}
func (r *ActiveGroupRepository) Get(ctx context.Context, userID string) (string, error) {
	var gid string
	err := r.DB.QueryRowContext(ctx, `SELECT group_id FROM user_active_groups WHERE user_id = $1`, userID).Scan(&gid)
	return gid, err
}
func (r *ActiveGroupRepository) Clear(ctx context.Context, userID string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM user_active_groups WHERE user_id = $1`, userID)
	return err
}
