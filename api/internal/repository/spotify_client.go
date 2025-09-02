package repository

import "context"

// 楽曲特徴量（必要最小限）
type AudioFeatures struct {
	Energy       float64 // 0..1
	Tempo        float64 // BPM
	Danceability float64 // 0..1
	Valence      float64 // 0..1
	Popularity   int     // 0..100
}

type Track struct {
	ID          string
	Name        string
	Artists     []string
	PreviewURL  string
	ExternalURL string
	ImageURL    string
	Features    AudioFeatures
}

type SpotifyClient interface {
	RecommendPopular(ctx context.Context, target AudioFeatures, market string, minPopularity, limit int) ([]Track, error)
}

// ---- Mock

type MockSpotify struct{}

func NewMockSpotify() SpotifyClient { return &MockSpotify{} }

func (m *MockSpotify) RecommendPopular(ctx context.Context, target AudioFeatures, market string, minPopularity, limit int) ([]Track, error) {
	if limit <= 0 {
		limit = 20
	}
	out := make([]Track, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, Track{
			ID:          "3n3Ppam7vgaVa1iaRUc9Lp",
			Name:        "Popular Mock Song",
			Artists:     []string{"Mock Artist"},
			PreviewURL:  "https://p.scdn.co/mp3-preview/mock.mp3",
			ExternalURL: "https://open.spotify.com/track/3n3Ppam7vgaVa1iaRUc9Lp",
			ImageURL:    "https://i.scdn.co/image/ab67616d00001e02mock",
			Features: AudioFeatures{
				Energy:       target.Energy,
				Tempo:        target.Tempo,
				Danceability: 0.65,
				Valence:      0.5,
				Popularity:   88,
			},
		})
	}
	return out, nil
}
