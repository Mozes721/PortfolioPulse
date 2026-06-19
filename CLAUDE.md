# PortfolioPulse — Claude Context

## Project Overview

**portfolio-pulse** is a Go service that fetches Trading212 portfolio positions at market open and close, stores rolling price snapshots in Upstash Redis, and persists historical records to Airtable. It is deployed on Fly.io and scheduled via cron.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| Hot data / cache | Upstash Redis (REST API) — local Redis for dev |
| Historical storage | Airtable |
| Portfolio source | Trading212 API |
| Deployment | Fly.io |
| Linting | golangci-lint |
| Scheduling | Fly.io cron machines |

## Architecture — DDD Layers

```
domain/          — pure Go types, interfaces, value objects; zero external deps
infrastructure/  — Redis client, Airtable client, Trading212 HTTP client
application/     — use-cases / service orchestration; wires domain + infra
cmd/             — entry points (main.go, cron jobs)
config/          — environment variable loading only
```

**Dependency direction:** `cmd → application → domain ← infrastructure`  
Infrastructure implements domain interfaces; application orchestrates them.  
Domain has no imports outside the Go standard library.

## Domain Concepts

- **Position** — a single holding (ticker, quantity, average price, current price, P&L)
- **Snapshot** — a point-in-time capture of all positions (open or close)
- **Period** — a reporting window; valid values: `1d`, `1w`, `1m`, `3m`, `1y`, `max`
- **MarketSession** — `open` or `close`; drives which cron fires

## Market Hours

Polling only occurs Monday–Friday. No snapshots are taken on weekends or UK/US public holidays. Market open = 08:00 UTC (LSE pre-market) and 14:30 UTC (NYSE open). Market close = 16:30 UTC (LSE) and 21:00 UTC (NYSE). The scheduler must check `time.Weekday()` before executing.

## Redis Key Schema

```
snapshot:{session}:{YYYY-MM-DD}   → JSON Snapshot  (TTL: 90 days)
positions:latest                   → JSON []Position (TTL: 24h)
period:{ticker}:{period}           → JSON []PricePoint (TTL: per period)
```

Valid period tokens: `1d`, `1w`, `1m`, `3m`, `1y`, `max` — no others.

## Airtable Schema

Table field names use **snake_case**. Canonical fields:

| Field | Type |
|---|---|
| `ticker` | Single line text |
| `date` | Date |
| `session` | Single select (`open`, `close`) |
| `quantity` | Number |
| `avg_price` | Currency |
| `current_price` | Currency |
| `pnl` | Currency |
| `pnl_pct` | Percent |
| `period` | Single select (`1d`,`1w`,`1m`,`3m`,`1y`,`max`) |

## Configuration

All runtime config is loaded from environment variables. No defaults may be hardcoded for credentials or URLs. See `.env.example` for the full list.

Required env vars:
- `UPSTASH_REDIS_REST_URL` / `UPSTASH_REDIS_REST_TOKEN` — production Redis
- `REDIS_URL` — local development Redis
- `TRADING_212_STORAGE_ACCESS_KEY` / `TRADING_212_SECRET_STORAGE_ACCESS_KEY`
- `AIRTABLE_API_KEY` / `AIRTABLE_BASE_ID` / `AIRTABLE_TABLE_NAME`

## Error Handling Standards

- Never discard errors with `_`. Every `err` must be checked.
- Wrap errors with context: `fmt.Errorf("fetchPositions: %w", err)`.
- Domain errors are typed values defined in `domain/errors.go`; infrastructure errors are wrapped and surfaced to the application layer.
- HTTP calls must have a timeout set on the `http.Client` (default 10s).
- Redis and Airtable operations retry once on transient errors (5xx, timeout) before returning an error to the caller.

## Infrastructure Notes

- **Fly.io** — app defined in `fly.toml`; secrets set via `fly secrets set`, never committed.
- **Cron machines** — two machines: one fires at market open, one at market close.
- **Upstash Redis** — accessed via REST (not TCP); use the official `upstash/go-redis` SDK or plain HTTP with bearer token.
- **Local dev** — `REDIS_URL` points to a local Redis instance; the config layer selects Upstash vs local based on presence of `UPSTASH_REDIS_REST_URL`.

## YouTube Tutorial Series

This project is the subject of a step-by-step Claude Code tutorial series. Each video introduces a new Claude Code feature. Keep code pedagogically clear; avoid premature abstraction.
