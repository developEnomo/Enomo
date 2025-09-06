package repository

import "context"

type AudioFeatures struct {
	Energy       float64 `json:"energy"`
	Tempo        float64 `json:"tempo"`
	Danceability float64 `json:"danceability"`
	Valence      float64 `json:"valence"`
	Popularity   int     `json:"popularity"`
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
