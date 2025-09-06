package usecase

import (
	"context"
	"math"
	"sort"

	"enomo/api/internal/config"
	"enomo/api/internal/repository"
)

type RecommendationsInput struct {
	GroupID        string
	ValencePreset  string  // "low" | "mid" | "high"
	PopularityBias float64 // 既定 0.7
	Market         string  // 既定 config.DefaultSpotifyMarket() or "JP"
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

	// 3) ターゲット特徴量（AudioFeaturesなしでも使える目標）
	target := buildTarget(E, in.ValencePreset)

	// 既定値の適用
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
		if def := config.DefaultSpotifyMarket(); def != "" {
			market = def
		} else {
			market = "JP"
		}
	}

	// 4) 候補取得（SpotifyClientは許可APIのみで実装）
	cands, err := u.spotify.RecommendPopular(ctx, repository.AudioFeatures{
		Energy:       target.Energy,       // 無視されてもOK
		Tempo:        target.Tempo,        // 無視されてもOK
		Danceability: target.Danceability, // 無視されてもOK
		Valence:      target.Valence,      // 無視されてもOK
	}, market, minPop, limit*2)
	if err != nil {
		return RecommendationsOutput{}, err
	}

	// 5) プレビューあり優先
	filtered := make([]repository.Track, 0, len(cands))
	for _, t := range cands {
		if t.PreviewURL != "" {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == 0 {
		filtered = cands
	}

	// 6) スコアリング（AudioFeatures欠損は自動で距離計算から除外）
	popBias := clamp(in.PopularityBias, 0.0, 1.0)
	s := make([]scored, 0, len(filtered))
	for _, t := range filtered {
		p := float64(t.Features.Popularity) / 100.0
		d := featureDistance(t.Features, target) // 0..1 小さいほど良い（欠損は除外）
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
		Energy:       0.25 + 0.6*E, // 0.25..0.85
		Tempo:        80 + 80*E,    // 80..160
		Danceability: 0.4 + 0.4*E,  // 0.4..0.8
		Valence:      v,            // 0.3 / 0.5 / 0.75
	}
}

// AudioFeaturesが欠損（0）なら、その成分は距離計算から除外する。
// すべて欠損の場合は中庸(=0.5)の距離を返す。
func featureDistance(f repository.AudioFeatures, t targetFeat) float64 {
	used := 0
	d2 := 0.0

	// Energy (0..1)
	if f.Energy > 0 {
		d2 += sq(f.Energy - t.Energy)
		used++
	}
	// Tempo (≈80±80 を0..1に正規化)
	if f.Tempo > 0 {
		tempoN := (f.Tempo - 80.0) / 80.0
		tempoT := (t.Tempo - 80.0) / 80.0
		d2 += sq(tempoN - tempoT)
		used++
	}
	// Danceability (0..1)
	if f.Danceability > 0 {
		d2 += sq(f.Danceability - t.Danceability)
		used++
	}
	// Valence (0..1)
	if f.Valence > 0 {
		d2 += sq(f.Valence - t.Valence)
		used++
	}

	if used == 0 {
		return 0.5 // 何も分からない → 中庸
	}
	// 成分数に応じて平均化し、0..1に収める
	return math.Min(1.0, math.Sqrt(d2)/math.Sqrt(float64(used)))
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
