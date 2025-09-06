package usecase

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"enomo/api/internal/repository"
)

type GroupSettingsUsecase interface {
	// 追加：両方まとめて取得
	Get(ctx context.Context, groupID string) (GroupSettings, error)
	// 追加：片方または両方を更新
	Update(ctx context.Context, in UpdateGroupSettingsInput, actorUserID string) (GroupSettings, error)

	// 互換：従来の呼び出しがあれば生かしておく
	GetHours(ctx context.Context, groupID string) (int, error)
	UpdateHours(ctx context.Context, groupID string, hours int, actorUserID string) (int, error)
}

type groupSettingsUsecase struct {
	groups *repository.GroupRepository
	// members *repository.GroupMemberRepository // 権限チェックを広げるならここに足す
}

func NewGroupSettingsUsecase(g *repository.GroupRepository) GroupSettingsUsecase {
	return &groupSettingsUsecase{groups: g}
}

type GroupSettings struct {
	GroupID              string
	GroupName            string
	RefreshIntervalHours int
}

type UpdateGroupSettingsInput struct {
	GroupID              string
	GroupName            *string
	RefreshIntervalHours *int
}

// --- read ---

func (u *groupSettingsUsecase) Get(ctx context.Context, groupID string) (GroupSettings, error) {
	if strings.TrimSpace(groupID) == "" {
		return GroupSettings{}, errors.New("group_id required")
	}
	name, err := u.groups.GetName(ctx, groupID)
	if err != nil {
		return GroupSettings{}, err
	}
	h, err := u.groups.GetRefreshIntervalHours(ctx, groupID)
	if err != nil {
		return GroupSettings{}, err
	}
	return GroupSettings{GroupID: groupID, GroupName: name, RefreshIntervalHours: h}, nil
}

// --- update (new) ---

func (u *groupSettingsUsecase) Update(ctx context.Context, in UpdateGroupSettingsInput, actorUserID string) (GroupSettings, error) {
	if strings.TrimSpace(in.GroupID) == "" {
		return GroupSettings{}, errors.New("group_id required")
	}
	// TODO: actorUserID の権限チェックを加えるならここで

	if in.GroupName == nil && in.RefreshIntervalHours == nil {
		return GroupSettings{}, errors.New("no_fields_to_update")
	}

	// name
	if in.GroupName != nil {
		n := strings.TrimSpace(*in.GroupName)
		if n == "" {
			return GroupSettings{}, errors.New("group_name required")
		}
		if utf8.RuneCountInString(n) > 50 {
			return GroupSettings{}, errors.New("group_name too long")
		}
		if err := u.groups.UpdateName(ctx, in.GroupID, n); err != nil {
			return GroupSettings{}, err
		}
	}

	// hours
	if in.RefreshIntervalHours != nil {
		if !isAllowed(*in.RefreshIntervalHours) {
			return GroupSettings{}, errors.New("refresh_interval_hours must be one of [12,24,48,72,96,120,144,168]")
		}
		if err := u.groups.UpdateRefreshIntervalHours(ctx, in.GroupID, *in.RefreshIntervalHours); err != nil {
			return GroupSettings{}, err
		}
	}

	// 最新値を返す
	name, err := u.groups.GetName(ctx, in.GroupID)
	if err != nil {
		return GroupSettings{}, err
	}
	h, err := u.groups.GetRefreshIntervalHours(ctx, in.GroupID)
	if err != nil {
		return GroupSettings{}, err
	}
	return GroupSettings{GroupID: in.GroupID, GroupName: name, RefreshIntervalHours: h}, nil
}

// --- 互換API ---

func (u *groupSettingsUsecase) GetHours(ctx context.Context, groupID string) (int, error) {
	return u.groups.GetRefreshIntervalHours(ctx, groupID)
}

func (u *groupSettingsUsecase) UpdateHours(ctx context.Context, groupID string, hours int, actorUserID string) (int, error) {
	if !isAllowed(hours) {
		return 0, errors.New("refresh_interval_hours must be one of [12,24,48,72,96,120,144,168]")
	}
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
