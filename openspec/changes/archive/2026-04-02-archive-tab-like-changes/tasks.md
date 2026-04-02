## 1. Scanner: unify change parsing

- [x] 1.1 Add `ArchiveDate time.Time` field to `ChangeInfo`
- [x] 1.2 Remove `ArchivedChange` struct
- [x] 1.3 Replace `ArchivedChanges []ArchivedChange` with `ArchivedChanges []ChangeInfo` in `ProjectInfo`
- [x] 1.4 Extract a `parseChangeDir(dir string, name string) ChangeInfo` helper from the active changes parsing loop
- [x] 1.5 Use `parseChangeDir` for both active changes and archived changes; set `ArchiveDate` from dir mtime for archived entries
- [x] 1.6 Update tests for the new `ArchivedChanges` type

## 2. UI: archive cursor and navigation

- [x] 2.1 Add `archiveCursor int` and `archiveArtifactTab int` to model
- [x] 2.2 When `detailTab == tabArchive`, j/k moves `archiveCursor`, left/right moves `archiveArtifactTab`
- [x] 2.3 Reset `archiveCursor` and `archiveArtifactTab` on project switch

## 3. UI: shared rendering

- [x] 3.1 Extract `renderChangeList(changes []ChangeInfo, cursor int, artifactTab int, width int, height int, showDate bool) string` from `renderChangesTab`
- [x] 3.2 Update `renderChangesTab` to call `renderChangeList` with active changes data and `showDate=false`
- [x] 3.3 Add `renderArchiveTab` that calls `renderChangeList` with archived changes data and `showDate=true`
- [x] 3.4 Wire `renderArchiveTab` into `renderDetailPanel` for `tabArchive` case

## 4. Tests and verification

- [x] 4.1 Verify build and all tests pass
