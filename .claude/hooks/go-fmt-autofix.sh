#!/usr/bin/env bash
# PostToolUse — Write|Edit
# Runs gofmt -s -w on whichever .go file Claude just wrote or edited.
# Never blocks (always exits 0).

set -euo pipefail

INPUT=$(cat)

# Extract the file path from the tool result JSON.
FILE=$(echo "$INPUT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
# PostToolUse delivers tool_input for Write/Edit
ti = d.get('tool_input', {})
path = ti.get('file_path', '')
print(path)
" 2>/dev/null || true)

if [[ "$FILE" == *.go ]] && [[ -f "$FILE" ]]; then
  gofmt -s -w "$FILE"
fi

exit 0
