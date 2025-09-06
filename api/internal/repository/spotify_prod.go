package repository

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type ProdSpotify struct {
	app *spotify.Client // app-only client (no user OAuth)
}

func NewProdSpotify(ctx context.Context) (*ProdSpotify, error) {
	id := os.Getenv("SPOTIFY_CLIENT_ID")
	secret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if id == "" || secret == "" {
		return nil, errors.New("SPOTIFY_CLIENT_ID / SPOTIFY_CLIENT_SECRET が未設定")
	}
	conf := &clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     spotifyauth.TokenURL,
	}
	httpClient := conf.Client(ctx) // 自動トークン取得/更新
	app := spotify.New(httpClient, spotify.WithRetry(true))
	return &ProdSpotify{app: app}, nil
}

// ====== 許可APIのみでの近似レコメンド実装 ======
// ・/v1/search で候補trackを集め、/v1/tracks でPopularity等を取得して整形
// ・Audio Features/Recommendations は使わない（0扱い/未設定扱い）
func (p *ProdSpotify) RecommendPopular(
	ctx context.Context,
	target AudioFeatures,
	market string,
	minPopularity, limit int,
) ([]Track, error) {

	// ---- 正規化 ----
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if minPopularity < 0 {
		minPopularity = 0
	}
	if minPopularity > 100 {
		minPopularity = 100
	}
	mkt := strings.ToUpper(strings.TrimSpace(market))
	if len(mkt) != 2 {
		mkt = ""
	}

	// ---- キーワード（市場ごと）で候補を検索 ----
	queries := seedQueriesForMarket(mkt) // ex: JP -> ["j-pop","anime","japanese pop","j rock"]
	opts := []spotify.RequestOption{spotify.Limit(20)}
	if mkt != "" {
		opts = append(opts, spotify.Country(mkt))
	}

	seen := make(map[spotify.ID]struct{}, limit*5)
	var ids []spotify.ID

	for _, q := range queries {
		// トラック検索（軽量・早い）
		res, err := p.app.Search(ctx, q, spotify.SearchTypeTrack, opts...)
		if err != nil || res == nil || res.Tracks == nil {
			continue
		}
		for _, t := range res.Tracks.Tracks {
			if t.ID == "" {
				continue
			}
			if _, ok := seen[t.ID]; ok {
				continue
			}
			seen[t.ID] = struct{}{}
			ids = append(ids, t.ID)
			if len(ids) >= limit*4 { // 余剰気味に集める
				break
			}
		}
		if len(ids) >= limit*4 {
			break
		}
	}

	// 検索でほぼ拾えなかった場合の保険：汎用キーワード
	if len(ids) == 0 {
		fallback := []string{"pop", "dance", "rock", "j-pop", "anime"}
		for _, q := range fallback {
			res, err := p.app.Search(ctx, q, spotify.SearchTypeTrack, opts...)
			if err != nil || res == nil || res.Tracks == nil {
				continue
			}
			for _, t := range res.Tracks.Tracks {
				if t.ID == "" {
					continue
				}
				if _, ok := seen[t.ID]; ok {
					continue
				}
				seen[t.ID] = struct{}{}
				ids = append(ids, t.ID)
				if len(ids) >= limit*3 {
					break
				}
			}
			if len(ids) >= limit*3 {
				break
			}
		}
	}

	if len(ids) == 0 {
		return []Track{}, nil
	}

	// ---- /v1/tracks でFullTrackを取得（Popularity/画像/外部URLなど） ----
	full := batchGetTracks(ctx, p.app, ids, mkt)

	// ---- 整形＋Popularityフィルタ ----
	out := make([]Track, 0, len(full))
	for _, ft := range full {
		if ft == nil {
			continue
		}
		pop := int(ft.Popularity)
		if pop < minPopularity {
			continue
		}
		artists := make([]string, 0, len(ft.Artists))
		for _, a := range ft.Artists {
			artists = append(artists, a.Name)
		}
        image := ""
        if len(ft.Album.Images) > 0 {
            image = ft.Album.Images[0].URL
        }
		external := ""
		if ft.ExternalURLs != nil {
			if u, ok := ft.ExternalURLs["spotify"]; ok {
				external = u
			}
		}
		// AudioFeaturesは使わない → 0扱い（Popularityのみセット）
		out = append(out, Track{
			ID:          ft.ID.String(),
			Name:        ft.Name,
			Artists:     artists,
			PreviewURL:  ft.PreviewURL,
			ExternalURL: external,
			ImageURL:    image,
			Features: AudioFeatures{
				Popularity: pop,
				// Energy/Danceability/Valence/Tempo は未知（0）
			},
		})
	}

	if len(out) == 0 {
		return []Track{}, nil
	}

	// ---- 人気順ベースで軽くシャッフル→上位を確保（多様性少し持たせる） ----
	sort.SliceStable(out, func(i, j int) bool { return out[i].Features.Popularity > out[j].Features.Popularity })
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 上位20%を軽くシャッフル
	k := len(out) / 5
	if k > 2 {
		perm := r.Perm(k)
		buf := make([]Track, k)
		for i, pi := range perm {
			buf[i] = out[pi]
		}
		copy(out[:k], buf)
	}

	// ---- limitで切る ----
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// ---- helpers ----

func batchGetTracks(ctx context.Context, cli *spotify.Client, ids []spotify.ID, market string) []*spotify.FullTrack {
	chunk := 50
	opts := []spotify.RequestOption{}
	if market != "" {
		opts = append(opts, spotify.Country(market))
	}
	out := make([]*spotify.FullTrack, 0, len(ids))
	for i := 0; i < len(ids); i += chunk {
		j := i + chunk
		if j > len(ids) {
			j = len(ids)
		}
		part, err := cli.GetTracks(ctx, ids[i:j], opts...)
		if err != nil {
			continue
		}
		out = append(out, part...)
	}
	return out
}

func seedQueriesForMarket(market string) []string {
	switch strings.ToUpper(strings.TrimSpace(market)) {
	case "JP":
		return []string{"j-pop", "anime", "japanese pop", "j rock"}
	case "US":
		return []string{"pop", "dance", "hip hop", "rock"}
	case "GB":
		return []string{"uk pop", "dance pop", "indie rock", "grime"}
	default:
		return []string{"pop", "dance", "rock"}
	}
}
