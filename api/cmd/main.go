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

	"enomo/server/internal/domain"
	"enomo/server/internal/repository"
	"enomo/server/internal/router"
	"enomo/server/internal/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	store := repository.NewPostgresTokenStore(db)

	regUC := usecase.NewUserRegisterUsecase(userRepo)
	loginUC := usecase.NewUserLoginUsecase(userRepo)
	logoutUC := usecase.NewUserLogoutUsecase(store)
	makeUC := usecase.NewGroupMakeUsecase(groupRepo)

	regH := domain.NewUserRegisterHandler(regUC)
	loginH := domain.NewUserLoginHandler(loginUC, store)
	logoutH := domain.NewUserLogoutHandler(logoutUC)
	makeH := domain.NewGroupMakeHandler(makeUC)

	e := router.New(regH, loginH, logoutH, makeH)

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
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
