#!/usr/bin/env bash
# PreToolUse (Bash, git push *): block the push if `go vet` or `go test` fail.
# `go test ./...` covers every package including cmd/snapshot and cmd/server,
# so Upstash/Trading212/Airtable client tests run here too once they exist.
set -uo pipefail

cd "${CLAUDE_PROJECT_DIR:-.}"

vet_output=$(go vet ./... 2>&1)
if [[ $? -ne 0 ]]; then
  reason="go vet failed:
$(printf '%s' "$vet_output" | tail -c 3000)

Fix before pushing."
  jq -n --arg reason "$reason" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: "deny", permissionDecisionReason: $reason}}'
  exit 0
fi

test_output=$(go test ./... 2>&1)
if [[ $? -ne 0 ]]; then
  reason="go test failed:
$(printf '%s' "$test_output" | tail -c 3000)

Fix before pushing."
  jq -n --arg reason "$reason" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: "deny", permissionDecisionReason: $reason}}'
  exit 0
fi

exit 0
