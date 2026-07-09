#!/usr/bin/env bash
# PreToolUse (Bash, git commit *): enforce Conventional Commits on the subject line.
# Only inspects an inline `-m "..."` argument — multi-line HEREDOC commit bodies
# (the norm for Claude-authored commits) aren't parseable here and are skipped.
set -uo pipefail

input=$(cat)
command=$(jq -r '.tool_input.command // empty' <<<"$input")

msg=$(grep -oP '(?<=-m[ =])("[^"]*"|'"'"'[^'"'"']*'"'"')' <<<"$command" | head -1)
msg="${msg%\"}"; msg="${msg#\"}"; msg="${msg%\'}"; msg="${msg#\'}"

[[ -n "$msg" ]] || exit 0

pattern='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-z0-9_-]+\))?: .+'
if ! [[ "$msg" =~ $pattern ]]; then
  reason="Commit subject \"$msg\" doesn't follow Conventional Commits (type(scope): subject). Allowed types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert."
  jq -n --arg reason "$reason" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: "deny", permissionDecisionReason: $reason}}'
fi

exit 0
