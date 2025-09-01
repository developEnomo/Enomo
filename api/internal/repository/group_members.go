package repository

import (
	"context"
	"database/sql"
	"time"
)

type GroupMember struct {
	GroupID  string
	UserID   string
	IsAdmin  bool
	JoinedAt time.Time
}

type GroupMemberRepository struct{ DB *sql.DB }

func NewGroupMemberRepository(db *sql.DB) *GroupMemberRepository {
	return &GroupMemberRepository{DB: db}
}

func (r *GroupMemberRepository) Add(ctx context.Context, groupID, userID string, isAdmin bool) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (group_id, user_id) DO UPDATE
		SET role = EXCLUDED.role
	`, groupID, userID, isAdmin)
	return err
}
