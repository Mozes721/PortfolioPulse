package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	UpstashRedisURL   string
	UpstashRedisToken string
	RedisURL          string

	Trading212Key    string
	Trading212Secret string

	AirtableAPIKey    string
	AirtableBaseID    string
	AirtableTableName string
}

// Load reads all required env vars and returns a validated Config.
func Load() (Config, error) {
	cfg := Config{
		UpstashRedisURL:   os.Getenv("UPSTASH_REDIS_REST_URL"),
		UpstashRedisToken: os.Getenv("UPSTASH_REDIS_REST_TOKEN"),
		RedisURL:          os.Getenv("REDIS_URL"),
		Trading212Key:     os.Getenv("TRADING_212_STORAGE_ACCESS_KEY"),
		Trading212Secret:  os.Getenv("TRADING_212_SECRET_STORAGE_ACCESS_KEY"),
		AirtableAPIKey:    os.Getenv("AIRTABLE_API_KEY"),
		AirtableBaseID:    os.Getenv("AIRTABLE_BASE_ID"),
		AirtableTableName: os.Getenv("AIRTABLE_TABLE_NAME"),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	var missing []string

	require := func(val, name string) {
		if strings.TrimSpace(val) == "" {
			missing = append(missing, name)
		}
	}

	require(c.Trading212Key, "TRADING_212_STORAGE_ACCESS_KEY")
	require(c.Trading212Secret, "TRADING_212_SECRET_STORAGE_ACCESS_KEY")
	// Airtable is optional — snapshots are still stored in Redis when not configured.

	if strings.TrimSpace(c.UpstashRedisURL) == "" && strings.TrimSpace(c.RedisURL) == "" {
		missing = append(missing, "UPSTASH_REDIS_REST_URL or REDIS_URL")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}

// IsUpstash reports whether Upstash REST credentials are present.
func (c Config) IsUpstash() bool {
	return c.UpstashRedisURL != "" && c.UpstashRedisToken != ""
}

// AirtableEnabled reports whether all three Airtable env vars are set.
func (c Config) AirtableEnabled() bool {
	return c.AirtableAPIKey != "" && c.AirtableBaseID != "" && c.AirtableTableName != ""
}
