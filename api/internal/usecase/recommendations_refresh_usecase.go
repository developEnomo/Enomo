package usecase

import (
	"context"
	"errors"
	"sort"

	"enomo/api/internal/repository"
)

type RecommendationsRefreshUsecase interface {
	// 一番スコアが高い1曲を選定して groups.track_id に保存して返す
	Refresh(ctx context.Context, in RecommendationsInput) (repository.Track, error)
}

type recommendationsRefreshUsecase struct {
	energyRepo repository.EnergyReader
	spotify    repository.SpotifyClient
	groups     *repository.GroupRepository
}

func NewRecommendationsRefreshUsecase(
	er repository.EnergyReader,
	sp repository.SpotifyClient,
	gr *repository.GroupRepository,
) RecommendationsRefreshUsecase {
	return &recommendationsRefreshUsecase{
		energyRepo: er,
		spotify:    sp,
		groups:     gr,
	}
}

func (u *recommendationsRefreshUsecase) Refresh(ctx context.Context, in RecommendationsInput) (repository.Track, error) {
	// 1) グループ energy → 0..1 に正規化
	energies, err := u.energyRepo.GroupEnergies(ctx, in.GroupID)
	if err != nil {
		return repository.Track{}, err
	}
	E := robustEnergy(energies)

	// 2) ターゲット特徴量
	target := buildTarget(E, defaultValence(in.ValencePreset))

	// 3) パラメータ既定値
	market := in.Market
	if market == "" {
		market = "JP"
	}
	minPop := in.MinPopularity
	if minPop == 0 {
		minPop = 60
	}
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}

	// 4) 候補（多めに拾う）
	cands, err := u.spotify.RecommendPopular(ctx, repository.AudioFeatures{
		Energy:       target.Energy,
		Tempo:        target.Tempo,
		Danceability: target.Danceability,
		Valence:      target.Valence,
	}, market, minPop, limit*2)
	if err != nil {
		return repository.Track{}, err
	}
	if len(cands) == 0 {
		return repository.Track{}, errors.New("no candidates")
	}

	// 5) プレビュー有りを優先、無ければプレビュー無しでフォールバック（★ここだけ変更）
	withPrev := make([]repository.Track, 0, len(cands))
	withoutPrev := make([]repository.Track, 0, len(cands))
	for _, t := range cands {
		if t.PreviewURL != "" {
			withPrev = append(withPrev, t)
		} else {
			withoutPrev = append(withoutPrev, t)
		}
	}
	pool := withPrev
	if len(pool) == 0 {
		pool = withoutPrev
	}
	if len(pool) == 0 {
		return repository.Track{}, errors.New("no candidates")
	}

	// 6) スコアリング（Execute と揃える）
	popBias := clamp(in.PopularityBias, 0.0, 1.0)
	type scored struct {
		repository.Track
		score float64
	}
	ss := make([]scored, 0, len(pool))
	for _, t := range pool {
		p := float64(t.Features.Popularity) / 100.0
		d := featureDistance(t.Features, target) // 0..1 小さいほど良い（欠損成分は距離計算から除外される実装ならそのまま）
		score := popBias*p + (1.0-popBias)*(1.0-d)
		ss = append(ss, scored{Track: t, score: score})
	}
	sort.Slice(ss, func(i, j int) bool { return ss[i].score > ss[j].score })

	best := ss[0].Track

	// 7) DBに反映（groups.track_id を更新）
	if err := u.groups.UpdateTrackID(ctx, in.GroupID, &best.ID); err != nil {
		return repository.Track{}, err
	}

	// 8) 更新時刻を刻む
	if err := u.groups.TouchRefreshedAt(ctx, in.GroupID); err != nil {
		// ここで失敗しても曲は保存済みなので、ポリシー次第ではログだけにして成功返しでもOK
		return repository.Track{}, err
	}

	return best, nil
}
