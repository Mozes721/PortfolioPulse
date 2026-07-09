#!/usr/bin/env bash
# PreToolUse — Bash(git commit *)
# Validates inline -m "..." commit messages against Conventional Commits.
# HEREDOC / multi-line commit bodies are silently skipped (can't be parsed from the command string).

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

# Extract the first -m "..." value (skip HEREDOC commits — not parseable here).
MSG=$(echo "$CMD" | python3 -c "
import sys, re
cmd = sys.stdin.read()
m = re.search(r'-m\s+[\"'"'"']([^\"'"'"']+)[\"'"'"']', cmd)
if m:
    print(m.group(1).split('\n')[0].strip())
" 2>/dev/null || true)

if [[ -z "$MSG" ]]; then
  # HEREDOC or unparseable — allow through.
  exit 0
fi

# Conventional Commits: type(scope): subject  OR  type: subject
PATTERN='^(feat|fix|docs|style|refactor|perf|test|chore|ci|build|revert)(\([a-z0-9/_-]+\))?: .{1,72}$'
if ! echo "$MSG" | grep -Eq "$PATTERN"; then
  echo "Commit message does not follow Conventional Commits format."
  echo "  Got:      $MSG"
  echo "  Expected: type(scope): subject"
  echo "  Types:    feat fix docs style refactor perf test chore ci build revert"
  exit 2
fi

exit 0
