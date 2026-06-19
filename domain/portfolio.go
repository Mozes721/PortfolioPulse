package domain

import (
	"context"
	"time"
)

// PortfolioFetcher fetches live positions from the broker.
type PortfolioFetcher interface {
	FetchPositions(ctx context.Context) ([]Position, error)
}

// PositionCache caches the latest positions in hot storage.
type PositionCache interface {
	SetLatest(ctx context.Context, positions []Position) error
	GetLatest(ctx context.Context) ([]Position, error)
}

// SnapshotStore persists and retrieves snapshots from hot storage.
type SnapshotStore interface {
	Save(ctx context.Context, s Snapshot) error
	Get(ctx context.Context, session MarketSession, date time.Time) (Snapshot, error)
}

// HistoryRepository persists snapshots to long-term storage.
type HistoryRepository interface {
	Append(ctx context.Context, s Snapshot) error
}
