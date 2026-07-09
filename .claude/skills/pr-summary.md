# PR Summary Skill

Generate a structured summary of all changes on the current branch relative to main, suitable for a pull request description or a `summary.md` file.

## Steps

1. Run `git log main..HEAD --oneline` to list commits on this branch.
2. Run `git diff main...HEAD --stat` to see which files changed and how many lines.
3. Run `git diff main...HEAD` to read the full diff.
4. Analyse the diff and produce a summary using the structure below.

## Output Format

Write the summary to `summary.md` in the repo root (create or overwrite), then print it to the user.

```markdown
# PR Summary — <branch-name>

## What changed
- <bullet per logical change group, one sentence each>

## Why
<one short paragraph explaining the motivation — infer from commit messages and code context>

## Files modified
| File | Change |
|------|--------|
| path/to/file.go | <brief description> |

## Testing notes
- <any test files added or changed>
- <manual verification steps if no tests cover the change>
```

## Rules
- Group related file changes into a single bullet; do not list every file individually in "What changed".
- Keep the "Why" section to 2–3 sentences max.
- If no tests were added for new logic, flag it explicitly under "Testing notes".
- Do not include secrets, tokens, or real credential values in the summary.