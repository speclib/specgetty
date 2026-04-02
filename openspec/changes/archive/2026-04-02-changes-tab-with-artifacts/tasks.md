## 1. Scanner: change data

- [x] 1.1 Add `ChangeInfo` struct with Name, ArtifactFiles, ArtifactContents, TasksTotal, TasksDone, SpecNames, SpecContents
- [x] 1.2 Add `Changes []ChangeInfo` to `ProjectInfo`
- [x] 1.3 In `ParseProjectInfo`, read `openspec/changes/` subdirs: discover .md files, read contents, parse tasks.md checkboxes, read specs/
- [x] 1.4 Add `parseTaskStats(content string) (total, done int)` helper that counts `- [x]` and `- [x]` patterns
- [x] 1.5 Add tests for change parsing (with artifacts, with tasks, without tasks, with specs, empty)

## 2. UI: change cursor and navigation

- [x] 2.1 Add `changeCursor int` and `changeArtifactTab int` fields to model
- [x] 2.2 When `detailTab == tabChanges` and `activeView == viewDetail`, j/k moves `changeCursor`
- [x] 2.3 Left/right arrows move `changeArtifactTab` when inside the changes tab (not the main detail tabs)
- [x] 2.4 Reset `changeCursor` and `changeArtifactTab` to 0 when switching projects or switching to changes tab

## 3. UI: changes tab rendering

- [x] 3.1 Add `renderChangesTab(width, height int) string` with horizontal split: ~30% change list, ~70% artifact content
- [x] 3.2 Render change list with selected item highlighted, showing task progress `(done/total)` when available
- [x] 3.3 Render artifact sub-tab header from discovered .md file stems + "specs" at the end
- [x] 3.4 Render selected artifact as markdown; for tasks, prepend a "Tasks: N/M complete" header
- [x] 3.5 Render specs sub-tab as list of spec names with first spec's content (reuse renderMarkdown)
- [x] 3.6 Show "No active changes" when project has no changes, "No specs in this change" when change has no specs
- [x] 3.7 Wire `renderChangesTab` into `renderDetailPanel` for `tabChanges` case

## 4. Tests and verification

- [x] 4.1 Add test for `parseTaskStats` with various checkbox patterns
- [x] 4.2 Verify build and all tests pass
