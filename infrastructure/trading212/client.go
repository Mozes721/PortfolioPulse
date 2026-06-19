package trading212

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Mozes721/portfolio-pulse/config"
	"github.com/Mozes721/portfolio-pulse/domain"
)

const (
	apiBase        = "https://live.trading212.com/api/v0"
	defaultTimeout = 10 * time.Second
)

// Client fetches live portfolio data from Trading212.
type Client struct {
	basicAuth string
	http      *http.Client
}

// New returns a Client configured from cfg.
// Trading212 uses HTTP Basic Auth: base64(accessKey:secretKey).
func New(cfg config.Config) *Client {
	creds := base64.StdEncoding.EncodeToString([]byte(cfg.Trading212Key + ":" + cfg.Trading212Secret))
	return &Client{
		basicAuth: "Basic " + creds,
		http:      &http.Client{Timeout: defaultTimeout},
	}
}

// FetchPositions implements domain.PortfolioFetcher.
func (c *Client) FetchPositions(ctx context.Context) ([]domain.Position, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+"/equity/portfolio", nil)
	if err != nil {
		return nil, fmt.Errorf("trading212.FetchPositions build request: %w", err)
	}
	req.Header.Set("Authorization", c.basicAuth)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("trading212.FetchPositions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("trading212.FetchPositions: HTTP %d — %s", resp.StatusCode, body)
	}

	var raw []apiPosition
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("trading212.FetchPositions decode: %w", err)
	}

	positions := make([]domain.Position, 0, len(raw))
	for _, r := range raw {
		pnl := (r.CurrentPrice - r.AveragePrice) * r.Quantity
		pnlPct := 0.0
		if r.AveragePrice > 0 {
			pnlPct = (r.CurrentPrice - r.AveragePrice) / r.AveragePrice * 100
		}
		positions = append(positions, domain.Position{
			Ticker:       r.Ticker,
			Quantity:     r.Quantity,
			AvgPrice:     r.AveragePrice,
			CurrentPrice: r.CurrentPrice,
			PnL:          pnl,
			PnLPct:       pnlPct,
		})
	}

	return positions, nil
}

// apiPosition mirrors the Trading212 equity/portfolio response shape.
type apiPosition struct {
	Ticker       string  `json:"ticker"`
	Quantity     float64 `json:"quantity"`
	AveragePrice float64 `json:"averagePrice"`
	CurrentPrice float64 `json:"currentPrice"`
}
