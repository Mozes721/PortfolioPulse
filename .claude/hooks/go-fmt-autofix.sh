#!/usr/bin/env bash
# PostToolUse (Write|Edit): gofmt -s -w whatever Go file Claude just touched.
# Best-effort — never blocks, never fails the turn.
set -uo pipefail

input=$(cat)
file=$(jq -r '.tool_response.filePath // .tool_input.file_path // empty' <<<"$input")

[[ "$file" == *.go ]] || exit 0
[[ -f "$file" ]] || exit 0

gofmt -s -w "$file" 2>/dev/null || true
exit 0
