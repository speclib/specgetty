## Why

The archive tab currently shows "Not yet implemented". Archived changes have the same directory structure as active changes (`.md` files, `specs/` subdirectory). Users want to browse archived changes the same way they browse active ones — list on the left with artifact sub-tabs on the right.

## What Changes

- Reuse the `ChangeInfo` struct and parsing logic for archived changes
- Replace `ArchivedChanges []ArchivedChange` with `ArchivedChanges []ChangeInfo` in `ProjectInfo`
- Add `archiveCursor` and `archiveArtifactTab` to the model
- Render the archive tab using the same split layout as the changes tab, showing archive date in the list
- Refactor the changes tab rendering into a shared function that both tabs call

## Capabilities

### New Capabilities
- `archive-tab`: Browse archived changes with the same layout as the changes tab

### Modified Capabilities

## Impact

- **scanner/scan.go**: Remove `ArchivedChange` struct, parse `openspec/archive/` into `[]ChangeInfo`, add archive date to `ChangeInfo`
- **ui/ui.go**: Add archive cursor/tab fields, refactor `renderChangesTab` into shared `renderChangeList`, wire archive tab
