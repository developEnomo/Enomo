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

func (r *GroupMemberRepository) List(ctx context.Context, groupID string, limit, offset int) ([]GroupMember, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT group_id, user_id, role, joined_at
		FROM group_members
		WHERE group_id = $1
		ORDER BY joined_at DESC
		LIMIT $2 OFFSET $3
	`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []GroupMember
	for rows.Next() {
		var m GroupMember
		if err := rows.Scan(&m.GroupID, &m.UserID, &m.IsAdmin, &m.JoinedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
