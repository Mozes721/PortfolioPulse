package domain

import "errors"

var (
	ErrInvalidPeriod    = errors.New("invalid period token")
	ErrPositionNotFound = errors.New("position not found")
	ErrSnapshotNotFound = errors.New("snapshot not found")
	ErrWeekend          = errors.New("market is closed on weekends")
)
