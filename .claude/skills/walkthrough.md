# Code Walkthrough Skill

Produce a plain-English, step-by-step walkthrough of a Go file, package, or the whole codebase — written for someone who reads Go fine but wants the bigger picture (what the code does and why) rather than a syntax explanation. Useful for narrating a tutorial video or onboarding into a part of the codebase you haven't touched yet.

## Input

- If the user names a specific file or package, walk through that.
- If no target is given, walk through the whole repository starting from `cmd/`, following the call chain down through `application/` → `domain/` → `infrastructure/` (per the DDD layering in CLAUDE.md).

## Steps

1. Read the target file(s) in full — do not summarize from memory or partial reads.
2. If the target calls into other packages, follow those calls one level deep (read the called function's signature and body) so the summary is accurate, but do not recursively expand every transitive dependency.
3. Identify the entry point (e.g. `main()`, an HTTP handler, a cron job) and trace execution in the order it actually runs, not file/declaration order.
4. For each step in the trace, note:
   - What runs
   - Why it runs (the business reason, not "because Go executes top to bottom")
   - Any non-obvious behavior (retries, TTLs, weekday guards, error wrapping) and the constraint driving it
5. Write the summary using the Output Format below.

## Output Format

```markdown
# Walkthrough — <target>

## What this does
<1-2 sentence plain-English purpose>

## Step by step
1. **<short step name>** (`file.go:line`) — <what happens and why>
2. ...

## Things to know
- <gotchas, invariants, or config dependencies a viewer/reader would otherwise miss>
```

## Rules
- Assume Go fluency: don't explain syntax, receivers, or `error`-handling idioms unless they encode a project-specific rule (e.g. the retry-once policy).
- Do explain domain vocabulary the first time it appears (`Snapshot`, `MarketSession`, `Period` tokens) — one clause is enough.
- Cite `file.go:line` for every step so it doubles as a script cue for screen recording.
- Keep each step to 1-3 sentences. If a step needs more, it's really two steps.
- Do not editorialize about code quality unless asked — this is a summary skill, not a review skill.
