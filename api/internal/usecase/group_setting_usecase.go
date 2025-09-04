package usecase

import (
	"context"
	"errors"

	"enomo/server/internal/repository"
)

type GroupSettingsUsecase interface {
	GetHours(ctx context.Context, groupID string) (int, error)
	UpdateHours(ctx context.Context, groupID string, hours int, actorUserID string) (int, error)
}

type groupSettingsUsecase struct {
	groups *repository.GroupRepository
	// members *repository.GroupMemberRepository // 必要なら管理者チェックに
}

func NewGroupSettingsUsecase(g *repository.GroupRepository) GroupSettingsUsecase {
	return &groupSettingsUsecase{groups: g}
}

func (u *groupSettingsUsecase) GetHours(ctx context.Context, groupID string) (int, error) {
	return u.groups.GetRefreshIntervalHours(ctx, groupID)
}

func (u *groupSettingsUsecase) UpdateHours(ctx context.Context, groupID string, hours int, actorUserID string) (int, error) {
	if !isAllowed(hours) {
		return 0, errors.New("refresh_interval_hours must be one of [12,24,48,72,96,120,144,168]")
	}
	// TODO: actorUserID が owner/admin か確認したい場合はここでチェック
	if err := u.groups.UpdateRefreshIntervalHours(ctx, groupID, hours); err != nil {
		return 0, err
	}
	return hours, nil
}

func isAllowed(h int) bool {
	switch h {
	case 12, 24, 48, 72, 96, 120, 144, 168:
		return true
	default:
		return false
	}
}
