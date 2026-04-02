## Why

The changes tab currently shows "Not yet implemented". Changes are the active work happening in an OpenSpec project — users need to see which changes exist, how much progress has been made on each, and read the individual artifacts (proposal, design, tasks, specs) within a change.

## What Changes

- Implement the `changes` tab with a three-level navigation:
  1. Left side: list of active changes (directories in `openspec/changes/`) with task progress (e.g. `18/20`)
  2. Right side header: sub-tabs for each `.md` file found in the change directory, plus a "specs" sub-tab for the `specs/` subdirectory
  3. Right side body: rendered markdown content of the selected artifact, or spec list+content for the specs sub-tab
- No schema awareness — simply discover all `.md` files in the change directory and show them as sub-tabs
- Parse `tasks.md` for checkbox stats: count `- [x]` and `- [ ]` patterns
- Show task progress both in the change list (compact) and in the change header (detailed)

## Capabilities

### New Capabilities
- `changes-tab`: Browse active changes with artifact sub-navigation and task progress display

### Modified Capabilities

## Impact

- **scanner/scan.go**: Add `ChangeInfo` struct (name, artifact files, task stats, spec names/contents) to `ProjectInfo`; parse change directories at scan time
- **ui/ui.go**: New `renderChangesTab()` with three-level navigation, `changeCursor` and `changeArtifactTab` on model, sub-tab rendering, task stats display
