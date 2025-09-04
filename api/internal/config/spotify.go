package config

import (
	"os"
	"strconv"
)

type SpotifyConfig struct {
	ClientID      string
	ClientSecret  string
	Market        string  // "JP"
	PopBias       float64 // 0.7 等
	MinPopularity int     // 60 等
}

func LoadSpotifyConfig() SpotifyConfig {
	return SpotifyConfig{
		ClientID:      os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret:  os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Market:        firstNonEmpty(os.Getenv("SPOTIFY_MARKET"), "JP"),
		PopBias:       parseFloatEnv("RECO_POPULARITY_BIAS", 0.7),
		MinPopularity: parseIntEnv("RECO_MIN_POP", 60),
	}
}

// --- 以下、共通ユーティリティ（既存に類似あれば流用可）
func firstNonEmpty(s, def string) string {
	if s != "" {
		return s
	}
	return def
}
func parseFloatEnv(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
func parseIntEnv(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
