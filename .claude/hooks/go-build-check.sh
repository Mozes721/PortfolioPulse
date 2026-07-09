#!/usr/bin/env bash
# PreToolUse — Bash(git commit *)
# Blocks the commit if go build fails.

set -euo pipefail

INPUT=$(cat)

CMD=$(echo "$INPUT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
print(d.get('tool_input', {}).get('command', ''))
" 2>/dev/null || true)

# Only run on git commit commands.
if [[ "$CMD" != *"git commit"* ]]; then
  exit 0
fi

cd "$CLAUDE_PROJECT_DIR"

OUTPUT=$(go build ./... 2>&1) || {
  echo "go build failed — fix before committing:"
  echo "$OUTPUT"
  exit 2
}

exit 0
