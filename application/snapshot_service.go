package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Mozes721/portfolio-pulse/domain"
)

// SnapshotService orchestrates fetching positions and persisting a snapshot.
type SnapshotService struct {
	fetcher   domain.PortfolioFetcher
	cache     domain.PositionCache
	snapshots domain.SnapshotStore
	history   domain.HistoryRepository
}

// NewSnapshotService wires the four domain dependencies together.
func NewSnapshotService(
	fetcher domain.PortfolioFetcher,
	cache domain.PositionCache,
	snapshots domain.SnapshotStore,
	history domain.HistoryRepository,
) *SnapshotService {
	return &SnapshotService{
		fetcher:   fetcher,
		cache:     cache,
		snapshots: snapshots,
		history:   history,
	}
}

// Run fetches the current portfolio and stores a snapshot for the given session.
// Returns domain.ErrWeekend when called on Saturday or Sunday.
func (s *SnapshotService) Run(ctx context.Context, session domain.MarketSession) error {
	if err := guardWeekday(time.Now().UTC()); err != nil {
		return err
	}

	positions, err := s.fetcher.FetchPositions(ctx)
	if err != nil {
		return fmt.Errorf("snapshotService.Run fetch: %w", err)
	}

	snap := domain.Snapshot{
		Date:      time.Now().UTC(),
		Session:   session,
		Period:    domain.Period1D,
		Positions: positions,
	}

	if err := s.cache.SetLatest(ctx, positions); err != nil {
		slog.ErrorContext(ctx, "position cache update failed", "err", err)
	}

	if err := s.snapshots.Save(ctx, snap); err != nil {
		return fmt.Errorf("snapshotService.Run save snapshot: %w", err)
	}

	if s.history != nil {
		if err := s.history.Append(ctx, snap); err != nil {
			return fmt.Errorf("snapshotService.Run append history: %w", err)
		}
	} else {
		slog.InfoContext(ctx, "Airtable not configured — skipping history append")
	}

	slog.InfoContext(ctx, "snapshot stored",
		"session", session,
		"positions", len(positions),
		"date", snap.Date.Format("2006-01-02"),
	)

	return nil
}

// guardWeekday returns domain.ErrWeekend for Saturday and Sunday.
func guardWeekday(t time.Time) error {
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return domain.ErrWeekend
	default:
		return nil
	}
}
