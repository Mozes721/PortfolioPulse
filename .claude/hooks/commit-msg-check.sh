#!/usr/bin/env bash
# PreToolUse — Bash(git commit *)
# Validates commit messages against Conventional Commits.
# Handles both inline -m "..." and HEREDOC (-m "$(cat <<'EOF' ... EOF)") forms.

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

# Extract the subject line from either commit form:
#   inline:  git commit -m "type: subject"
#   heredoc: git commit -m "$(cat <<'EOF'\ntype: subject\n...\nEOF\n)"
MSG=$(echo "$CMD" | python3 -c "
import sys, re

cmd = sys.stdin.read()

# 1. Try to extract heredoc body between EOF markers.
heredoc = re.search(r\"<<'?EOF'?\s*\n(.*?)\nEOF\", cmd, re.DOTALL)
if heredoc:
    # First non-empty line is the subject.
    for line in heredoc.group(1).splitlines():
        line = line.strip()
        if line:
            print(line)
            break
    sys.exit(0)

# 2. Fall back to inline -m '...' or -m \"...\"
inline = re.search(r'-m\s+[\"\']([^\"\']+)[\"\']', cmd)
if inline:
    print(inline.group(1).split('\n')[0].strip())
" 2>/dev/null || true)

if [[ -z "$MSG" ]]; then
  # Unparseable form — allow through rather than false-block.
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
