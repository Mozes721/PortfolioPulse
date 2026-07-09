#!/usr/bin/env bash
# PreToolUse — Bash(git push *)
# Blocks the push if go vet or go test fails.

set -euo pipefail

INPUT=$(cat)

CMD=$(echo "$INPUT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
print(d.get('tool_input', {}).get('command', ''))
" 2>/dev/null || true)

# Only run on git push commands.
if [[ "$CMD" != *"git push"* ]]; then
  exit 0
fi

cd "$CLAUDE_PROJECT_DIR"

VET_OUTPUT=$(go vet ./... 2>&1) || {
  echo "go vet failed — fix before pushing:"
  echo "$VET_OUTPUT"
  exit 2
}

TEST_OUTPUT=$(go test ./... 2>&1) || {
  echo "go test failed — fix before pushing:"
  echo "$TEST_OUTPUT"
  exit 2
}

exit 0
