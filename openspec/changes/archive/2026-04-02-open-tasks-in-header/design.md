## Context

The persistent header in `renderDetailPanel` currently shows a stats line:
```
Specs: 3  Changes: 2 active  Archived: 5
```

Per-change task stats (`TasksTotal`, `TasksDone`) are already parsed in `ChangeInfo` but only displayed in the changes tab list. There is no project-level aggregation.

## Goals / Non-Goals

**Goals:**
- Show aggregate open task count in the persistent header stats line
- Keep the header compact (single stats line)

**Non-Goals:**
- Individual task listing in the header
- Per-change breakdown in the header
- New "tasks" tab implementation

## Decisions

**Aggregate at scan time, not render time**
Add `TasksTotal` and `TasksDone` fields to `ProjectInfo`. Compute them in `ParseProjectInfo` by summing across all active `ChangeInfo` entries. This avoids re-computing on every render frame.

Alternative: compute in the UI render loop. Rejected because it mixes data aggregation with rendering and would repeat the sum on every frame.

**Conditionally display**
Only append the tasks stat when `TasksTotal > 0`. Projects with no tasks.md files keep the current clean stats line.

**Format: `Tasks: done/total`**
Append `Tasks: 4/12` to the existing stats line, matching the `Key: value` pattern already used.

## Risks / Trade-offs

- [Minimal risk] Stats line could get long on narrow terminals → the line already wraps gracefully via lipgloss padding, and one more stat is marginal.
