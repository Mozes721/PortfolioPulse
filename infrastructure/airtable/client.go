package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mozes721/portfolio-pulse/config"
	"github.com/Mozes721/portfolio-pulse/domain"
)

const (
	baseURL        = "https://api.airtable.com/v0"
	defaultTimeout = 10 * time.Second
)

// Client writes snapshots to Airtable.
type Client struct {
	apiKey    string
	baseID    string
	tableName string
	http      *http.Client
}

// New returns a Client configured from cfg.
func New(cfg config.Config) *Client {
	return &Client{
		apiKey:    cfg.AirtableAPIKey,
		baseID:    cfg.AirtableBaseID,
		tableName: cfg.AirtableTableName,
		http:      &http.Client{Timeout: defaultTimeout},
	}
}

// Append implements domain.HistoryRepository.
// Each Position in the snapshot becomes one Airtable record.
func (c *Client) Append(ctx context.Context, s domain.Snapshot) error {
	records := make([]record, 0, len(s.Positions))
	for _, p := range s.Positions {
		records = append(records, record{
			Fields: map[string]any{
				domain.FieldTicker:       p.Ticker,
				domain.FieldDate:         s.Date.UTC().Format("2006-01-02"),
				domain.FieldSession:      string(s.Session),
				domain.FieldPeriod:       string(s.Period),
				domain.FieldQuantity:     p.Quantity,
				domain.FieldAvgPrice:     p.AvgPrice,
				domain.FieldCurrentPrice: p.CurrentPrice,
				domain.FieldPnL:          p.PnL,
				domain.FieldPnLPct:       p.PnLPct,
			},
		})
	}

	// Airtable accepts up to 10 records per request.
	for i := 0; i < len(records); i += 10 {
		end := i + 10
		if end > len(records) {
			end = len(records)
		}
		if err := c.createRecords(ctx, records[i:end]); err != nil {
			return fmt.Errorf("airtable.Append batch %d: %w", i/10, err)
		}
	}

	return nil
}

type record struct {
	Fields map[string]any `json:"fields"`
}

type createRequest struct {
	Records []record `json:"records"`
}

func (c *Client) createRecords(ctx context.Context, records []record) error {
	payload, err := json.Marshal(createRequest{Records: records})
	if err != nil {
		return fmt.Errorf("airtable.createRecords marshal: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", baseURL, c.baseID, c.tableName)

	var lastErr error
	for attempt := range 2 {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return fmt.Errorf("airtable.createRecords build request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("airtable.createRecords: %w", err)
			if attempt == 0 {
				continue
			}
			return lastErr
		}
		defer resp.Body.Close() //nolint:gocritic // defer in loop is intentional: each iteration has its own resp

		if resp.StatusCode >= 500 && attempt == 0 {
			// Retry once on transient server error.
			lastErr = fmt.Errorf("airtable.createRecords: HTTP %d (retrying)", resp.StatusCode)
			continue
		}

		if resp.StatusCode >= 400 {
			var apiErr struct {
				Error struct {
					Type    string `json:"type"`
					Message string `json:"message"`
				} `json:"error"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&apiErr) //nolint:errcheck // best-effort: enrich error message only
			return fmt.Errorf("airtable.createRecords: HTTP %d — %s: %s",
				resp.StatusCode, apiErr.Error.Type, apiErr.Error.Message)
		}

		return nil
	}

	return lastErr
}
