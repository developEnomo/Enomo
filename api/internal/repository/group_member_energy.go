package repository

import (
	"context"
	"database/sql"
)

type EnergyReader interface {
	GroupEnergies(ctx context.Context, groupID string) ([]int, error)
}

type GroupEnergyRepository struct{ DB *sql.DB }

func NewGroupEnergyRepository(db *sql.DB) *GroupEnergyRepository {
	return &GroupEnergyRepository{DB: db}
}

func (r *GroupEnergyRepository) GroupEnergies(ctx context.Context, groupID string) ([]int, error) {
	const q = `
SELECT u.energy_value
FROM group_members gm
JOIN users u ON u.id = gm.user_id
WHERE gm.group_id = $1
`
	rows, err := r.DB.QueryContext(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
