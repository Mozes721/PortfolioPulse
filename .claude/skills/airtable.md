# Airtable Skill — PortfolioPulse

Guidelines for working with the Airtable integration in this project.

## Project-specific schema

Table field names use **snake_case**. Canonical fields (defined as constants in `domain/fields.go`):

| Field | Airtable type | Go constant |
|-------|--------------|-------------|
| `ticker` | Single line text | `FieldTicker` |
| `date` | Date | `FieldDate` |
| `session` | Single select (`open`, `close`) | `FieldSession` |
| `quantity` | Number | `FieldQuantity` |
| `avg_price` | Currency | `FieldAvgPrice` |
| `current_price` | Currency | `FieldCurrentPrice` |
| `pnl` | Currency | `FieldPnl` |
| `pnl_pct` | Percent | `FieldPnlPct` |
| `period` | Single select (`1d`,`1w`,`1m`,`3m`,`1y`,`max`) | `FieldPeriod` |

## When writing Airtable infrastructure code

1. **Never use raw string literals** for field names — always reference the constant from `domain/fields.go`.
2. The Airtable client lives in `infrastructure/` and implements an interface defined in `domain/`.
3. Configuration (`AIRTABLE_API_KEY`, `AIRTABLE_BASE_ID`, `AIRTABLE_TABLE_NAME`) comes from env vars loaded by `config/` — never hardcode.
4. Wrap all Airtable errors with context: `fmt.Errorf("airtable.CreateRecord: %w", err)`.
5. Retry once on HTTP 5xx before returning the error.

## When using the Airtable MCP tools in this session

Use the MCP Airtable tools in this order:
1. `search_bases` — find the PortfolioPulse base by name.
2. `list_tables_for_base(baseId)` — confirm table and field IDs.
3. `get_table_schema` — required before filtering on single/multi-select fields to get choice IDs.
4. `list_records_for_table` or `create_records_for_table` — read or write records.

Always use Airtable internal IDs (base ID, table ID, field ID) — never substitute human-readable names for IDs in API calls.

## Adding a new field

1. Add the constant to `domain/fields.go`.
2. Update the Airtable infrastructure client to map the new field.
3. Update this skill and `CLAUDE.md` schema table to keep them in sync.