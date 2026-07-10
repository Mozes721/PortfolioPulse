# PR Summary — main

_Note: `main` has no open feature branch (already up to date with `origin/main`), so this summarizes the most recent commits (`fb5892f..HEAD`: `87102d6`, `fef06fa`, `be4a7ac`) rather than a `main...HEAD` diff._

## What changed
- Added Claude Code automation: five hooks (`commit-msg-check`, `go-build-check`, `go-lint-check`, `go-prepush-check`, `go-fmt-autofix`) wired up in `.claude/settings.json` to enforce Conventional Commits, run `go build`/`golangci-lint` pre-commit, `go vet`/`go test` pre-push, and auto-`gofmt` on save.
- Added a GitHub Actions CI workflow (`.github/workflows/ci.yml`) with lint (gofmt + golangci-lint) and build/test jobs, later bumped to `golangci-lint-action@v7` for golangci-lint v2 compatibility.
- Migrated `.golangci.yml` to the v2 config schema (versioned `linters`/`formatters` sections, new exclusion presets) to match the CI action bump.
- Fixed Airtable field-name constants (`FieldPnl`→`FieldPnL`, `FieldPnlPct`→`FieldPnLPct`) to follow Go acronym casing, and added a `Period` field end-to-end (`domain/position.go` → `application/snapshot_service.go` → `infrastructure/airtable/client.go`).
- Hardened `infrastructure/airtable/client.go`'s `createRecords` with a single retry on transient errors (network failure or HTTP 5xx), per the project's error-handling standard.
- Corrected `go.mod` from an invalid `go 1.26.2` directive to `go 1.22`.
- Documented the new `/ddd-go` and `/airtable` skills as pointers inside `.claude/rules/architecture.md`.

## Why
The hooks and CI pipeline enforce the project's Go style, lint, and Conventional Commits rules automatically rather than relying on manual review. The Airtable/domain fixes correct a naming-convention bug and complete the `Period` field that Airtable's schema already required, while the retry logic and `go.mod` fix bring the code in line with the error-handling and toolchain standards defined in `CLAUDE.md`.

## Files modified
| File | Change |
|------|--------|
| `.claude/hooks/*.sh` (5 files) | New pre-commit/pre-push/post-edit hooks for build, lint, test, commit-msg, and gofmt |
| `.claude/settings.json` | New — wires the hooks above into Claude Code's PreToolUse/PostToolUse events |
| `.github/workflows/ci.yml` | New — lint + build/test CI jobs, later bumped to `golangci-lint-action@v7` |
| `.golangci.yml` | Migrated to golangci-lint v2 config schema |
| `.claude/rules/architecture.md` | Added references to `/ddd-go` and `/airtable` skills |
| `.claude/skills/airtable.md` | Fixed `FieldPnl`/`FieldPnlPct` constant names to `FieldPnL`/`FieldPnLPct` |
| `domain/position.go` | Added `Period` field to `Snapshot` |
| `application/snapshot_service.go` | Set `Period: domain.Period1D` when building a snapshot |
| `infrastructure/airtable/client.go` | Added `Period` to the Airtable record payload; added single-retry logic for transient errors |
| `go.mod` | Fixed invalid Go version `1.26.2` → `1.22` |

## Testing notes
- No new or updated test files are included in this range — the retry logic in `client.go` and the `Period` plumbing have no unit test coverage.
- Manual verification: hooks were presumably exercised locally via real commits/pushes; CI workflow itself has not yet been confirmed green after the `golangci-lint-action@v7` bump introduced in `fef06fa`.