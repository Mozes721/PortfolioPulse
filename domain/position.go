package domain

import "time"

// MarketSession identifies whether a snapshot was taken at open or close.
type MarketSession string

const (
	SessionOpen  MarketSession = "open"
	SessionClose MarketSession = "close"
)

// Position is a single portfolio holding at a point in time.
type Position struct {
	Ticker       string  `json:"ticker"`
	Quantity     float64 `json:"quantity"`
	AvgPrice     float64 `json:"avg_price"`
	CurrentPrice float64 `json:"current_price"`
	PnL          float64 `json:"pnl"`
	PnLPct       float64 `json:"pnl_pct"`
}

// Snapshot is an immutable capture of all positions at a market session boundary.
type Snapshot struct {
	Date      time.Time     `json:"date"`
	Session   MarketSession `json:"session"`
	Period    Period        `json:"period"`
	Positions []Position    `json:"positions"`
}
