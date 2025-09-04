package repository

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type AudioFeatures struct {
	Energy       float64 `json:"energy"`       // 0..1
	Tempo        float64 `json:"tempo"`        // BPM
	Danceability float64 `json:"danceability"` // 0..1
	Valence      float64 `json:"valence"`      // 0..1
	Popularity   int     `json:"popularity"`   // 0..100
}

type Track struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Artists     []string      `json:"artists"`
	PreviewURL  string        `json:"preview_url"`
	ExternalURL string        `json:"external_url"`
	ImageURL    string        `json:"image_url"`
	Features    AudioFeatures `json:"features"`
}

type SpotifyClient interface {
	RecommendPopular(ctx context.Context, target AudioFeatures, market string, minPopularity, limit int) ([]Track, error)
}

type MockSpotify struct {
	pool []TrackMeta
}

type TrackMeta struct {
	ID     string
	Name   string
	Artist string
}

func NewMockSpotify() SpotifyClient {
	pool := []TrackMeta{
		{ID: "6H1RjVyNruCmrBEWRbJmbA", Name: "Lemon", Artist: "米津玄師"},
		{ID: "5YqltLsPpxzOkG8u1S0h6Z", Name: "Pretender", Artist: "Official髭男dism"},
		{ID: "6n7nd5iceYpXVwcx8VPpxF", Name: "アイドル", Artist: "YOASOBI"},
		{ID: "2Z2pdWl1lK5Xg8Xdeu1TEn", Name: "花束", Artist: "back number"},
		{ID: "6RxyyxQCkZtYVyKqZdhhmF", Name: "踊", Artist: "Aimer"},
		{ID: "1X9YHQ1CEV4VpWwfvX2gYd", Name: "夜に駆ける", Artist: "YOASOBI"},
		{ID: "0yK7Jg6R4O6nH4YwZlEGlB", Name: "紅蓮華", Artist: "LiSA"},
		{ID: "3Pyox4Om5v5xux8OQkGZQm", Name: "シンデレラボーイ", Artist: "Saucy Dog"},
		{ID: "0wGXIJtJmwUEDX0z3WsmSf", Name: "残響散歌", Artist: "Aimer"},
		{ID: "2tGvwE8GcFKwNdaxlK1j9u", Name: "白日", Artist: "King Gnu"},
	}
	return &MockSpotify{pool: pool}
}

func (m *MockSpotify) RecommendPopular(ctx context.Context, target AudioFeatures, market string, minPopularity, limit int) ([]Track, error) {
	if limit <= 0 {
		limit = 20
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	basePops := []int{92, 88, 85, 83, 80, 78, 75, 73, 70, 68}

	jitter := func(x, amplitude float64) float64 {
		return clamp01(x + (r.Float64()*2-1)*amplitude)
	}
	jitterTempo := func(bpm, amp float64) float64 {
		return math.Max(60, math.Min(200, bpm+(r.Float64()*2-1)*amp))
	}

	out := make([]Track, 0, limit)
	for i := 0; i < limit; i++ {
		idx := i % len(m.pool)
		meta := m.pool[idx]
		pop := basePops[idx%len(basePops)]
		if pop < minPopularity {
			pop = minPopularity + (idx % 5)
		}
		feat := AudioFeatures{
			Energy:       jitter(target.Energy, 0.08),
			Tempo:        jitterTempo(target.Tempo, 8),
			Danceability: jitter(target.Danceability, 0.08),
			Valence:      jitter(target.Valence, 0.1),
			Popularity:   pop,
		}
		out = append(out, Track{
			ID:          meta.ID,
			Name:        meta.Name,
			Artists:     []string{meta.Artist},
			PreviewURL:  "https://p.scdn.co/mp3-preview/mock.mp3",
			ExternalURL: "https://open.spotify.com/track/" + meta.ID,
			ImageURL:    "https://i.scdn.co/image/ab67616d00001e02mock",
			Features:    feat,
		})
	}
	return out[:limit], nil
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
