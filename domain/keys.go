package domain

import (
	"fmt"
	"time"
)

// SnapshotKey returns the Redis key for a session snapshot on a given date.
func SnapshotKey(session MarketSession, date time.Time) string {
	return fmt.Sprintf("snapshot:%s:%s", session, date.UTC().Format("2006-01-02"))
}

// PositionsLatestKey returns the Redis key for the latest positions cache.
func PositionsLatestKey() string {
	return "positions:latest"
}

// PeriodKey returns the Redis key for a ticker's price history over a period.
func PeriodKey(ticker string, period Period) string {
	return fmt.Sprintf("period:%s:%s", ticker, period)
}
