package redis

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

const defaultTimeout = 10 * time.Second

// Client sends commands to Upstash Redis over REST or a local Redis-compatible endpoint.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// New returns a Client configured from cfg.
func New(cfg config.Config) *Client {
	c := &Client{
		http: &http.Client{Timeout: defaultTimeout},
	}
	if cfg.IsUpstash() {
		c.baseURL = cfg.UpstashRedisURL
		c.token = cfg.UpstashRedisToken
	} else {
		c.baseURL = cfg.RedisURL
	}
	return c
}

// SetLatest implements domain.PositionCache.
func (c *Client) SetLatest(ctx context.Context, positions []domain.Position) error {
	data, err := json.Marshal(positions)
	if err != nil {
		return fmt.Errorf("redis.SetLatest marshal: %w", err)
	}
	return c.set(ctx, domain.PositionsLatestKey(), string(data), 24*time.Hour)
}

// GetLatest implements domain.PositionCache.
func (c *Client) GetLatest(ctx context.Context) ([]domain.Position, error) {
	raw, err := c.get(ctx, domain.PositionsLatestKey())
	if err != nil {
		return nil, fmt.Errorf("redis.GetLatest: %w", err)
	}
	var positions []domain.Position
	if err := json.Unmarshal([]byte(raw), &positions); err != nil {
		return nil, fmt.Errorf("redis.GetLatest unmarshal: %w", err)
	}
	return positions, nil
}

// Save implements domain.SnapshotStore.
func (c *Client) Save(ctx context.Context, s domain.Snapshot) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("redis.Save marshal: %w", err)
	}
	return c.set(ctx, domain.SnapshotKey(s.Session, s.Date), string(data), 90*24*time.Hour)
}

// Get implements domain.SnapshotStore.
func (c *Client) Get(ctx context.Context, session domain.MarketSession, date time.Time) (domain.Snapshot, error) {
	raw, err := c.get(ctx, domain.SnapshotKey(session, date))
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("redis.Get: %w", err)
	}
	var s domain.Snapshot
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return domain.Snapshot{}, fmt.Errorf("redis.Get unmarshal: %w", err)
	}
	return s, nil
}

// set issues a SET key value EX seconds command via the Upstash REST API.
func (c *Client) set(ctx context.Context, key, value string, ttl time.Duration) error {
	_, err := c.do(ctx, "SET", key, value, "EX", int(ttl.Seconds()))
	return err
}

// get issues a GET key command via the Upstash REST API and returns the string result.
func (c *Client) get(ctx context.Context, key string) (string, error) {
	raw, err := c.do(ctx, "GET", key)
	if err != nil {
		return "", err
	}
	var result string
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("redis.get decode result: %w", err)
	}
	return result, nil
}

// do sends a single Redis command as a JSON array to the Upstash REST endpoint.
func (c *Client) do(ctx context.Context, args ...any) (json.RawMessage, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("redis.do marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("redis.do build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("redis.do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("redis.do: server error %d", resp.StatusCode)
	}

	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  string          `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("redis.do decode envelope: %w", err)
	}
	if envelope.Error != "" {
		return nil, fmt.Errorf("redis: %s", envelope.Error)
	}

	return envelope.Result, nil
}
