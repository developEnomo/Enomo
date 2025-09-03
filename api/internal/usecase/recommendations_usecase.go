package usecase

import (
	"context"
	"math"
	"sort"

	"enomo/api/internal/repository"
)

type RecommendationsInput struct {
	GroupID        string
	ValencePreset  string  // "low" | "mid" | "high"
	PopularityBias float64 // 既定 0.7
	Market         string  // 既定 "JP"
	MinPopularity  int     // 既定 60
	Limit          int     // 既定 20
}

type RecommendationsOutput struct {
	GroupID string             `json:"group_id"`
	Params  map[string]any     `json:"params"`
	Tracks  []repository.Track `json:"tracks"`
}

type RecommendationsUsecase interface {
	Execute(ctx context.Context, in RecommendationsInput) (RecommendationsOutput, error)
}

type recommendationsUsecase struct {
	energyRepo repository.EnergyReader
	spotify    repository.SpotifyClient
}

func NewRecommendationsUsecase(er repository.EnergyReader, sp repository.SpotifyClient) RecommendationsUsecase {
	return &recommendationsUsecase{energyRepo: er, spotify: sp}
}

func (u *recommendationsUsecase) Execute(ctx context.Context, in RecommendationsInput) (RecommendationsOutput, error) {
	// 1) グループの energy 値取得
	energies, err := u.energyRepo.GroupEnergies(ctx, in.GroupID)
	if err != nil {
		return RecommendationsOutput{}, err
	}

	// 2) ロバスト平均 → 0..1 正規化
	E := robustEnergy(energies)

	// 3) ターゲット特徴量
	target := buildTarget(E, in.ValencePreset)

	minPop := in.MinPopularity
	if minPop == 0 {
		minPop = 60
	}
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	market := in.Market
	if market == "" {
		market = "JP"
	}

	// 4) 候補取得（モック or 実装差替え）
	cands, err := u.spotify.RecommendPopular(ctx, repository.AudioFeatures{
		Energy:       target.Energy,
		Tempo:        target.Tempo,
		Danceability: target.Danceability,
		Valence:      target.Valence,
	}, market, minPop, limit*2)
	if err != nil {
		return RecommendationsOutput{}, err
	}

	// 5) プレビューありのみ
	filtered := make([]repository.Track, 0, len(cands))
	for _, t := range cands {
		if t.PreviewURL != "" {
			filtered = append(filtered, t)
		}
	}

	// 6) スコアリング
	popBias := clamp(in.PopularityBias, 0.0, 1.0)
	s := make([]scored, 0, len(filtered))
	for _, t := range filtered {
		p := float64(t.Features.Popularity) / 100.0
		d := featureDistance(t.Features, target) // 0..1 小さいほど良い
		score := popBias*p + (1.0-popBias)*(1.0-d)
		s = append(s, scored{Track: t, score: score})
	}
	sort.Slice(s, func(i, j int) bool { return s[i].score > s[j].score })

	// 7) 上位 limit 件
	out := takeTop(s, limit)

	return RecommendationsOutput{
		GroupID: in.GroupID,
		Params: map[string]any{
			"energy":          E,
			"valence_preset":  defaultValence(in.ValencePreset),
			"popularity_bias": popBias,
			"target":          target,
			"market":          market,
		},
		Tracks: out,
	}, nil
}

// ------- ユーティリティ -------

func defaultValence(v string) string {
	if v == "" {
		return "mid"
	}
	return v
}

type targetFeat struct {
	Energy, Tempo, Danceability, Valence float64
}

func buildTarget(E float64, valence string) targetFeat {
	v := 0.5
	switch valence {
	case "low":
		v = 0.3
	case "high":
		v = 0.75
	}
	return targetFeat{
		Energy:       0.25 + 0.6*E,
		Tempo:        80 + 80*E,
		Danceability: 0.4 + 0.4*E,
		Valence:      v,
	}
}

func featureDistance(f repository.AudioFeatures, t targetFeat) float64 {
	tempoN := (f.Tempo - 80.0) / 80.0
	tempoT := (t.Tempo - 80.0) / 80.0
	d2 := sq(f.Energy-t.Energy) +
		sq(tempoN-tempoT) +
		sq(f.Danceability-t.Danceability) +
		sq(f.Valence-t.Valence)
	return math.Min(1.0, math.Sqrt(d2)/2.0)
}

func robustEnergy(vals []int) float64 {
	if len(vals) == 0 {
		return 0.5
	}
	cp := append([]int(nil), vals...)
	sort.Ints(cp)
	med := cp[len(cp)/2]
	sum := 0.0
	for _, v := range cp {
		if v < med-1 {
			v = med - 1
		}
		if v > med+1 {
			v = med + 1
		}
		sum += float64(v)
	}
	avg := sum / float64(len(cp))
	return (avg - 1.0) / 4.0 // 1..5 -> 0..1
}

func sq(x float64) float64 { return x * x }

func clamp(x, a, b float64) float64 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

type scored struct {
	repository.Track
	score float64
}


func takeTop(s []scored, k int) []repository.Track {
	if k > len(s) {
		k = len(s)
	}
	out := make([]repository.Track, 0, k)
	for i := 0; i < k; i++ {
		out = append(out, s[i].Track)
	}
	return out
}
