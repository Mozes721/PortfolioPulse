package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Mozes721/portfolio-pulse/application"
	"github.com/Mozes721/portfolio-pulse/config"
	"github.com/Mozes721/portfolio-pulse/domain"
	"github.com/Mozes721/portfolio-pulse/infrastructure/airtable"
	"github.com/Mozes721/portfolio-pulse/infrastructure/redis"
	"github.com/Mozes721/portfolio-pulse/infrastructure/trading212"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	redisClient := redis.New(cfg)
	t212Client := trading212.New(cfg)

	var history domain.HistoryRepository
	if cfg.AirtableEnabled() {
		history = airtable.New(cfg)
	}

	svc := application.NewSnapshotService(t212Client, redisClient, redisClient, history)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /snapshot/{session}", handleSnapshot(svc))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("server started", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func handleSnapshot(svc *application.SnapshotService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := domain.MarketSession(r.PathValue("session"))
		if session != domain.SessionOpen && session != domain.SessionClose {
			http.Error(w, "session must be 'open' or 'close'", http.StatusBadRequest)
			return
		}

		if err := svc.Run(r.Context(), session); err != nil {
			if errors.Is(err, domain.ErrWeekend) {
				http.Error(w, "market closed on weekends", http.StatusConflict)
				return
			}
			slog.ErrorContext(r.Context(), "snapshot failed", "err", err)
			http.Error(w, "snapshot failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","session":"%s"}`, session)
	}
}
