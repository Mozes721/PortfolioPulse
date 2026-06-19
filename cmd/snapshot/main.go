package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/Mozes721/portfolio-pulse/application"
	"github.com/Mozes721/portfolio-pulse/config"
	"github.com/Mozes721/portfolio-pulse/domain"
	"github.com/Mozes721/portfolio-pulse/infrastructure/airtable"
	"github.com/Mozes721/portfolio-pulse/infrastructure/redis"
	"github.com/Mozes721/portfolio-pulse/infrastructure/trading212"
)

func main() {
	// godotenv is a no-op in production where Fly.io injects secrets directly.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	// MARKET_SESSION is set by the Fly.io cron machine: "open" or "close".
	session := domain.MarketSession(os.Getenv("MARKET_SESSION"))
	if session == "" {
		session = domain.SessionOpen
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	redisClient := redis.New(cfg)
	t212Client := trading212.New(cfg)

	var history domain.HistoryRepository
	if cfg.AirtableEnabled() {
		history = airtable.New(cfg)
	} else {
		slog.Info("Airtable not configured — history persistence disabled")
	}

	svc := application.NewSnapshotService(t212Client, redisClient, redisClient, history)

	if err := svc.Run(ctx, session); err != nil {
		if errors.Is(err, domain.ErrWeekend) {
			slog.Info("skipping snapshot — market closed on weekend")
			return
		}
		slog.ErrorContext(ctx, "snapshot failed", "err", err)
		os.Exit(1)
	}
}
