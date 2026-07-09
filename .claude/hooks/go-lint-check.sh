#!/usr/bin/env bash
# PreToolUse (Bash, git commit *): block the commit if golangci-lint finds issues.
set -uo pipefail

cd "${CLAUDE_PROJECT_DIR:-.}"

output=$(golangci-lint run ./... 2>&1)
status=$?

if [[ $status -ne 0 ]]; then
  reason="golangci-lint failed:
$(printf '%s' "$output" | tail -c 3000)

Fix the issues above (or run \`golangci-lint run --fix ./...\`) before committing."
  jq -n --arg reason "$reason" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: "deny", permissionDecisionReason: $reason}}'
fi

exit 0
