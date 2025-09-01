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
