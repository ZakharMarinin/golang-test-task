package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"testovoe/internal/application"
	"testovoe/internal/config"
	"testovoe/internal/http/handlers"
	"testovoe/internal/http/router"
	"testovoe/internal/storage"
	"testovoe/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := setupLogger(cfg.Env)

	db, err := storage.New(ctx, cfg.Postgres.Addr)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	httpRouter := chi.NewRouter()

	useCase := usecase.NewUseCase(log, db)

	httpHandlers := handlers.NewHTTPHandler(useCase)

	httpRouter.Use(middleware.RequestID)
	httpRouter.Use(middleware.Recoverer)

	router.Router(ctx, httpRouter, httpHandlers)

	app := application.NewApplication(ctx, cfg, log, httpRouter)

	app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	app.Shutdown()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
