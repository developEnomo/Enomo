package config

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"enomo/api/internal/repository"
)

var (
	spOnce   sync.Once
	spClient repository.SpotifyClient
	spErr    error
)

// truthy: "true", "1", "yes", "on" を true とみなす
func truthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// NewSpotifyClient は .env/環境変数に応じて Mock or Prod を返す。
// 優先度:
// 1) SPOTIFY_USE_MOCK=true/1/yes/on → 強制モック
// 2) CLIENT_ID/SECRET が空（.env未読込を含む） → 自動モック
// 3) それ以外 → Prod(Spotify API)
func NewSpotifyClient(ctx context.Context) (repository.SpotifyClient, error) {
	spOnce.Do(func() {
		useMock := truthy(os.Getenv("SPOTIFY_USE_MOCK"))
		id := strings.TrimSpace(os.Getenv("SPOTIFY_CLIENT_ID"))
		secret := strings.TrimSpace(os.Getenv("SPOTIFY_CLIENT_SECRET"))

		switch {
		case useMock:
			log.Println("[spotify] using MOCK (SPOTIFY_USE_MOCK is truthy)")
			spClient = repository.NewMockSpotify()
		case id == "" || secret == "":
			log.Println("[spotify] using MOCK (CLIENT_ID/SECRET not set)")
			spClient = repository.NewMockSpotify()
		default:
			var prod repository.SpotifyClient
			prod, spErr = repository.NewProdSpotify(ctx)
			if spErr != nil {
				log.Printf("[spotify] prod init failed: %v; falling back to MOCK\n", spErr)
				spClient = repository.NewMockSpotify()
				spErr = nil // フォールバック成功として扱う
				return
			}
			log.Println("[spotify] using PROD (client-credentials)")
			spClient = prod
		}
	})
	return spClient, spErr
}

// デフォルトの market（"JP" など）を環境変数から取得。
// 2文字なら大文字にして返す。それ以外は空文字（未指定）。
func DefaultSpotifyMarket() string {
	m := strings.TrimSpace(os.Getenv("SPOTIFY_MARKET"))
	if len(m) == 2 {
		return strings.ToUpper(m)
	}
	return ""
}
