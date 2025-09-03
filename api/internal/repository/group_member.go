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

type GroupWithMeta struct {
	GroupID     string
	GroupName   string
	MemberCount int
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

func (r *GroupMemberRepository) UserList(ctx context.Context, groupID string, limit, offset int) ([]GroupMember, error) {
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

func (r *GroupMemberRepository) GroupList(ctx context.Context, userID string, limit, offset int) ([]GroupWithMeta, error) {
	rows, err := r.DB.QueryContext(ctx, `
	WITH user_groups AS (
	SELECT DISTINCT group_id
		FROM group_members
	WHERE user_id = $1
	ORDER BY group_id
	LIMIT $2 OFFSET $3
	)
	SELECT
	ug.group_id,
	g.name,
	COUNT(m.user_id) AS member_count
	FROM user_groups ug
	JOIN groups g             ON g.id = ug.group_id
	LEFT JOIN group_members m ON m.group_id = ug.group_id
	GROUP BY ug.group_id, g.name
	ORDER BY MAX(g.created_at) DESC
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []GroupWithMeta
	for rows.Next() {
		var g GroupWithMeta
		if err := rows.Scan(&g.GroupID, &g.GroupName, &g.MemberCount); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *GroupMemberRepository) Remove(ctx context.Context, groupID, userID string) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `
		DELETE FROM group_members
		WHERE group_id = $1 AND user_id = $2
	`, groupID, userID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (r *GroupMemberRepository) IsMember(ctx context.Context, groupID, userID string) (bool, bool, error) {
	var isAdmin bool
	err := r.DB.QueryRowContext(ctx, `
		SELECT role
		FROM group_members
		WHERE group_id = $1 AND user_id = $2
	`, groupID, userID).Scan(&isAdmin)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, isAdmin, nil
}
