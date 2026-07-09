# DDD Go Skill — PortfolioPulse

Reference guide for applying Domain-Driven Design in this Go codebase. Use this skill whenever adding a new feature, domain concept, or infrastructure integration.

## Layer map

```
domain/          — pure Go: types, interfaces, value objects, domain errors
                   ZERO external imports (stdlib only)

infrastructure/  — all external I/O: Redis, Airtable, Trading212 HTTP
                   implements interfaces defined in domain/

application/     — use-cases / orchestration
                   wires domain interfaces to infrastructure; no raw HTTP or DB calls

cmd/             — entry points (main.go, cron jobs)
                   thin: load config → wire deps → call application layer → exit

config/          — reads env vars, returns typed struct; no business logic
```

**Dependency direction:** `cmd → application → domain ← infrastructure`

Circular imports are forbidden. The linter enforces this.

## Checklist for a new domain concept

- [ ] Define the type in `domain/` (e.g. `domain/position.go`)
- [ ] Define its interface in `domain/` if infrastructure must implement it (e.g. `domain/repository.go`)
- [ ] Add any new domain errors to `domain/errors.go` as typed sentinel values
- [ ] Add Redis key helpers to `domain/keys.go` — no ad-hoc key strings elsewhere
- [ ] Add Airtable field constants to `domain/fields.go` — no raw string literals in infra
- [ ] Implement the interface in `infrastructure/`
- [ ] Wire the implementation in the relevant `application/` use-case
- [ ] Keep `cmd/` thin — only config loading, wiring, and a single application call

## Error handling rules

```go
// GOOD
result, err := repo.Save(ctx, snapshot)
if err != nil {
    return fmt.Errorf("snapshotService.Save: %w", err)
}

// BAD — never discard errors
result, _ := repo.Save(ctx, snapshot)
```

- Domain errors: typed sentinels in `domain/errors.go`
- Infrastructure errors: wrap with `fmt.Errorf("client.Operation: %w", err)`
- Retry transient HTTP errors (5xx, timeout) exactly once before surfacing

## Context and I/O

- `context.Context` MUST be the first parameter of every function that does I/O
- HTTP clients MUST have an explicit `Timeout` (default 10 s)
- Log with `log/slog` key-value pairs — no `fmt.Println` in production paths
- No `panic` in library code; no `init()` functions

## Market session guard

The weekday check lives in the **application layer**, not in `cmd/`:

```go
func (s *SnapshotService) TakeSnapshot(ctx context.Context, session domain.MarketSession) error {
    if wd := time.Now().UTC().Weekday(); wd == time.Saturday || wd == time.Sunday {
        slog.Info("skipping snapshot: weekend", "weekday", wd)
        return nil
    }
    // proceed...
}
```

## Naming conventions

| Kind | Style | Example |
|------|-------|---------|
| Exported type/func | CamelCase | `MarketSession`, `FetchPositions` |
| Unexported | camelCase | `parseResponse` |
| Acronyms | uppercase | `HTTPClient`, `APIKey`, `RedisURL` |
| Package | lowercase single word | `domain`, `infrastructure` |
| Airtable fields | snake_case constant | `FieldAvgPrice = "avg_price"` |
| Redis keys | from `domain/keys.go` helpers | `SnapshotKey(session, date)` |

## Quick reference: valid period tokens

`1d` `1w` `1m` `3m` `1y` `max` — any other value is a `domain.ErrInvalidPeriod`.