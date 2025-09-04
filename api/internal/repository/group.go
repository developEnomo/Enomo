package repository

import (
	"context"
	"database/sql"
	"time"
)

type Group struct {
	ID        string
	Name      string
	OwnerID   string
	TrackID   *string
	CreatedAt time.Time
}

type GroupRepository struct{ DB *sql.DB }

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

func (r *GroupRepository) Create(ctx context.Context, g *Group) error {
	var track interface{}
	if g.TrackID == nil || *g.TrackID == "" {
		track = nil
	} else {
		track = *g.TrackID
	}
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO groups (name, owner_id, track_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, g.Name, g.OwnerID, track).Scan(&g.ID, &g.CreatedAt)
}

func (r *GroupRepository) Delete(ctx context.Context, id string) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM groups WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (r *GroupRepository) OwnerID(ctx context.Context, id string) (string, error) {
	var owner string
	err := r.DB.QueryRowContext(ctx, `SELECT owner_id FROM groups WHERE id = $1`, id).Scan(&owner)
	return owner, err
}

func (r *GroupRepository) GetOwnerID(ctx context.Context, id string) (string, error) {
	return r.OwnerID(ctx, id)
}

func (r *GroupRepository) OwnsAny(ctx context.Context, ownerID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM groups WHERE owner_id = $1)
	`, ownerID).Scan(&exists)
	return exists, err
}

// 更新頻度を取得
func (r *GroupRepository) GetRefreshIntervalHours(ctx context.Context, groupID string) (int, error) {
	var h int
	err := r.DB.QueryRowContext(ctx, `
		SELECT refresh_interval_hours
		  FROM groups
		 WHERE id = $1
	`, groupID).Scan(&h)
	return h, err
}

// 更新頻度を変更
func (r *GroupRepository) UpdateRefreshIntervalHours(ctx context.Context, groupID string, hours int) error {
	_, err := r.DB.ExecContext(ctx, `
UPDATE groups
   SET refresh_interval_hours = $1,
       next_refresh_at = refreshed_at + make_interval(hours => $1)
 WHERE id = $2
`, hours, groupID)
	return err
}

// Refresh実行後にrefreshed_atを刻む
func (r *GroupRepository) TouchRefreshedAt(ctx context.Context, groupID string) error {
	_, err := r.DB.ExecContext(ctx, `
UPDATE groups
   SET refreshed_at   = now(),
       next_refresh_at = now() + make_interval(hours => refresh_interval_hours)
 WHERE id = $1
`, groupID)
	return err
}

// groups.track_id を1曲だけ更新（nilでクリアも可）
func (r *GroupRepository) UpdateTrackID(ctx context.Context, groupID string, trackID *string) error {
	const q = `UPDATE groups SET track_id = $2 WHERE id = $1;`
	res, err := r.DB.ExecContext(ctx, q, groupID, trackID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
