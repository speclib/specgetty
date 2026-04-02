## Context

Archived changes live in `openspec/archive/` and have the same internal structure as active changes: `.md` files, `specs/` subdirectory. The current `ArchivedChange` struct only stores name and date. We need full `ChangeInfo` data.

## Goals / Non-Goals

**Goals:**
- Parse archived changes with the same `ChangeInfo` struct used for active changes
- Add an `ArchiveDate` field to `ChangeInfo` (populated from dir mtime for archive entries, zero for active changes)
- Refactor the rendering into a shared function so both tabs use identical layout
- Archive list shows the date alongside the name

**Non-Goals:**
- Different visual treatment for archived vs active changes
- Sorting archived changes differently (keep alphabetical, matching active changes)

## Decisions

### 1. Add ArchiveDate to ChangeInfo, remove ArchivedChange struct
`ChangeInfo` gains `ArchiveDate time.Time`. For active changes it's zero. For archived changes it's the dir mtime. This eliminates the separate `ArchivedChange` struct.

### 2. Shared renderChangeList function
Extract the core rendering from `renderChangesTab` into `renderChangeList(changes []ChangeInfo, cursor int, artifactTab int, width int, height int, showDate bool) string`. Both `renderChangesTab` and `renderArchiveTab` call this with their respective data and cursor.

### 3. Archive list label includes date
Active changes: `my-change (3/5)`. Archived changes: `my-change (2026-01-15)`.

## Risks / Trade-offs

- **Trade-off**: Parsing all archived change contents at scan time could be slow for projects with many archived changes (e.g. OpenSpec has 77). Acceptable for now — can add lazy loading later if needed.
