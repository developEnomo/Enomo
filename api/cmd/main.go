package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"enomo/api/internal/config"
	"enomo/api/internal/domain"
	"enomo/api/internal/repository"
	"enomo/api/internal/router"
	"enomo/api/internal/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

func main() {
	port := getenv("PORT", "8080")
	dsn := getenv("DATABASE_URL", "postgres://enomo:enomo@localhost:5432/enomo_dev?sslmode=disable")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// ベースコンテキスト
	ctx := context.Background()

	// --- repositories ---
	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	memberRepo := repository.NewGroupMemberRepository(db)
	store := repository.NewPostgresTokenStore(db)
	activeRepo := repository.NewActiveGroupRepository(db)
	energyRepo := repository.NewGroupEnergyRepository(db)

	// Spotify クライアント（.env未設定→自動でMock／設定あり→Prod）
	spotifyClient, err := config.NewSpotifyClient(ctx)
	if err != nil {
		log.Fatalf("spotify init error: %v", err)
	}

	// --- usecases ---
	meUC := usecase.NewUserMeUsecase(userRepo)
	regUC := usecase.NewUserRegisterUsecase(userRepo)
	loginUC := usecase.NewUserLoginUsecase(userRepo)
	logoutUC := usecase.NewUserLogoutUsecase(store)
	glistUC := usecase.NewUserGroupsUsecase(memberRepo)
	renameUC := usecase.NewUserUpdateDisplayNameUsecase(userRepo)
	energyUC := usecase.NewUserUpdateEnergyValueUsecase(userRepo)
	udelUC := usecase.NewUserDeleteUsecase(userRepo, groupRepo)
	nowUC := usecase.NewGroupNowUsecase(activeRepo, memberRepo)
	makeUC := usecase.NewGroupMakeUsecase(groupRepo, memberRepo)
	addUC := usecase.NewGroupAddUsecase(memberRepo)
	ulistUC := usecase.NewGroupListUsecase(memberRepo, userRepo)
	gdelUC := usecase.NewGroupDeleteUsecase(groupRepo)
	leaveUC := usecase.NewGroupLeaveUsecase(memberRepo, groupRepo)
	recoUC := usecase.NewRecommendationsUsecase(energyRepo, spotifyClient)
	recoRefreshUC := usecase.NewRecommendationsRefreshUsecase(energyRepo, spotifyClient, groupRepo)
	settingsUC := usecase.NewGroupSettingsUsecase(groupRepo)

	// --- handlers ---
	meH := domain.NewUserMeHandler(meUC, store)
	regH := domain.NewUserRegisterHandler(regUC)
	loginH := domain.NewUserLoginHandler(loginUC, store)
	logoutH := domain.NewUserLogoutHandler(logoutUC, store)
	glistH := domain.NewUserGroupsHandler(glistUC, store)
	renameH := domain.NewUserRenameHandler(renameUC, store)
	energyH := domain.NewUserEnergyHandler(energyUC, store)
	udelH := domain.NewUserDeleteHandler(udelUC, store)
	nowH := domain.NewGroupNowHandler(nowUC, store)
	makeH := domain.NewGroupMakeHandler(makeUC, store)
	addH := domain.NewGroupAddHandler(addUC, store)
	ulistH := domain.NewGroupListHandler(ulistUC)
	gdelH := domain.NewGroupDeleteHandler(gdelUC, store)
	leaveH := domain.NewGroupLeaveHandler(leaveUC, store)
	recoH := domain.NewRecommendationsHandler(recoUC, recoRefreshUC)
	settingsH := domain.NewGroupSettingsHandler(settingsUC, store)

	// --- router ---
	e := router.New(
		meH,
		regH,
		loginH,
		logoutH,
		glistH,
		renameH,
		energyH,
		udelH,
		nowH,
		makeH,
		addH,
		ulistH,
		gdelH,
		leaveH,
		recoH,
		settingsH,
	)

	// --- start / graceful shutdown ---
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ★ 自動更新ジョブを起動
	jctx, cancelJobs := context.WithCancel(context.Background())
	defer cancelJobs()
	go startAutoRefresher(jctx, groupRepo, recoRefreshUC)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// 自動更新ループ：毎分チェックし、期限到来のグループを順に更新
func startAutoRefresher(ctx context.Context, groupRepo *repository.GroupRepository, refresher usecase.RecommendationsRefreshUsecase) {
	ticker := time.NewTicker(1 * time.Minute) // チェック間隔。環境変数で可変にしてもOK
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			const batch = 20 // 1回で処理する件数上限（任意）
			ids, err := groupRepo.FindDueGroups(ctx, batch)
			if err != nil {
				log.Printf("[auto-refresh] FindDueGroups error: %v", err)
				continue
			}
			for _, gid := range ids {
				_, err := refresher.Refresh(ctx, usecase.RecommendationsInput{
					GroupID:        gid,
					ValencePreset:  "mid",
					PopularityBias: 0.7,
					Market:         "JP",
					MinPopularity:  60,
					Limit:          20,
				})
				if err != nil {
					log.Printf("[auto-refresh] Refresh group=%s error: %v", gid, err)
					continue
				}
				log.Printf("[auto-refresh] refreshed group=%s", gid)
			}
		}
	}
}
