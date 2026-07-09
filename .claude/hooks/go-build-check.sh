#!/usr/bin/env bash
# PreToolUse (Bash, git commit *): block the commit if `go build ./...` fails.
set -uo pipefail

cd "${CLAUDE_PROJECT_DIR:-.}"

output=$(go build ./... 2>&1)
status=$?

if [[ $status -ne 0 ]]; then
  reason="go build failed:
$(printf '%s' "$output" | tail -c 3000)

Fix the build before committing."
  jq -n --arg reason "$reason" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: "deny", permissionDecisionReason: $reason}}'
fi

exit 0
