# PortfolioPulse — Enforceable Claude Rules

These rules apply to every file generated or modified in this project.
Violating any rule is a hard error — do not proceed without fixing it.

---

## Go Code Style

- All Go code MUST be `gofmt`-formatted. Tabs for indentation, no trailing whitespace.
- Code MUST pass `golangci-lint run` with the project's `.golangci.yml`. Do not suppress linter warnings with `//nolint` unless a comment explains the exact reason.
- Use standard Go naming: `CamelCase` for exported identifiers, `camelCase` for unexported. Acronyms stay uppercase (`HTTP`, `API`, `URL`, `ID`).
- One package per directory. Package names are lowercase, single words, matching the directory name.
- Keep functions short. If a function exceeds ~40 lines, consider splitting it.

## Error Handling

- NEVER ignore errors. `_, err := f()` followed by no `err` check is forbidden.
- ALWAYS wrap errors with context: `fmt.Errorf("operationName: %w", err)`.
- Domain-layer errors MUST be typed sentinel values defined in `domain/errors.go`, not raw `errors.New` strings scattered across packages.
- HTTP clients MUST have an explicit `Timeout` set. Default: 10 seconds.
- Retry transient errors (HTTP 5xx, network timeout) exactly once before surfacing to the caller.

## DDD Layer Structure

- `domain/` — pure Go: types, interfaces, value objects, domain errors. Zero external imports (stdlib only).
- `infrastructure/` — all external I/O: Redis, Airtable, Trading212 HTTP. Implements interfaces defined in `domain/`.
- `application/` — use-cases only. Wires domain interfaces to infrastructure implementations. No raw HTTP or DB calls here.
- `cmd/` — `main.go` and cron entry points. Thin: load config, wire dependencies, call application layer, exit.
- `config/` — reads env vars and returns a typed config struct. No business logic.
- Cross-layer imports: `cmd` may import `application`, `config`. `application` may import `domain`. `infrastructure` may import `domain`. `domain` imports nothing outside stdlib. Circular imports are forbidden.
- For the full checklist when adding a new domain concept, naming conventions, and code patterns: use the `/ddd-go` skill.

## Trading212 API

- Official docs: https://docs.trading212.com/api
- Base URL: `https://live.trading212.com/api/v0`
- Authentication: HTTP Basic Auth — `Authorization: Basic base64(accessKey:secretKey)`.
  - `TRADING_212_STORAGE_ACCESS_KEY` is the access key (username)
  - `TRADING_212_SECRET_STORAGE_ACCESS_KEY` is the secret key (password)
- Do NOT use a plain bearer token or `Authorization: Bearer` — that is the deprecated approach and returns 401.
- Portfolio positions endpoint: `GET /equity/portfolio`
- Always consult https://docs.trading212.com/api before adding new Trading212 endpoints.

## Configuration & Secrets

- ALL configuration (URLs, tokens, keys, table names) MUST come from environment variables.
- NEVER hardcode credentials, URLs, or any environment-specific value in source code.
- Redis connection details are loaded from `UPSTASH_REDIS_REST_URL` + `UPSTASH_REDIS_REST_TOKEN`. This is used for **both production and local dev** — no local Redis instance or Docker is required. `REDIS_URL` exists as a fallback only. The `config/` package selects between them; callers never inspect env vars directly.
- `fly.toml` MUST NOT contain secret values. Secrets are set via `fly secrets set`.
- `.env` and `.env.local` are never committed.

## Redis Key Conventions

- All Redis keys follow the schema defined in CLAUDE.md.
- Period tokens MUST be one of: `1d`, `1w`, `1m`, `3m`, `1y`, `max`. Any other value is a domain error.
- Key construction helpers live in `domain/keys.go`; no ad-hoc key string formatting outside that file.

## Airtable Field Naming

- All Airtable field names use **snake_case** (e.g. `avg_price`, `pnl_pct`, `current_price`).
- Never use camelCase or PascalCase for Airtable field names.
- Field names are defined as typed constants in `domain/fields.go`; raw string literals are not permitted in infrastructure code.
- For the full field schema, MCP tool usage order, and infrastructure code rules: use the `/airtable` skill.

## Market Hours Awareness

- Snapshot polling MUST NOT execute on weekends (`time.Saturday`, `time.Sunday`).
- Before any API call to Trading212, verify `time.Now().UTC().Weekday()` is Mon–Fri.
- The weekday guard lives in the application layer, not in `main.go` or cron entry points.
- Do not hardcode holiday lists; weekend-only guard is the minimum required check.

## Testing

- Unit tests live in `_test.go` files alongside the package under test.
- Test fixtures MUST NOT contain real API keys, real account IDs, or real position data.
- Use table-driven tests for functions with multiple input/output cases.
- Infrastructure clients MUST be hidden behind domain interfaces so they can be faked in tests without hitting real APIs.

## General

- No `panic` in library code. Reserve `panic` for `main()` startup failures where the process cannot safely continue.
- No `init()` functions.
- Context (`context.Context`) MUST be the first parameter of every function that performs I/O.
- Log with structured key-value pairs (`log/slog`); no `fmt.Println` in production paths.
