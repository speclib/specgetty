## 1. Scanner: spec data

- [x] 1.1 Add `SpecNames []string` and `SpecContents map[string]string` to `ProjectInfo`
- [x] 1.2 In `ParseProjectInfo`, read `openspec/specs/` subdirectories, sort names, read each `spec.md` into `SpecContents`
- [x] 1.3 Add test for `ParseProjectInfo` spec reading (with specs, without specs, missing spec.md)

## 2. UI: spec cursor and navigation

- [x] 2.1 Add `specCursor int` field to model
- [x] 2.2 When `detailTab == tabSpecs` and `activeView == viewDetail`, j/k/up/down moves `specCursor` instead of `fileCursor`
- [x] 2.3 Reset `specCursor` to 0 when switching projects or switching to specs tab

## 3. UI: specs tab rendering

- [x] 3.1 Add `renderSpecsTab(width, height int) string` with horizontal split: ~30% spec list, ~70% content
- [x] 3.2 Render spec list with selected item highlighted, scrollable if list exceeds height
- [x] 3.3 Render selected spec content using `renderMarkdown`, truncated to panel height
- [x] 3.4 Show "No specs found" when project has no specs, "No spec.md found" when spec dir has no spec.md
- [x] 3.5 Wire `renderSpecsTab` into `renderDetailPanel` for `tabSpecs` case

## 4. Tests and verification

- [x] 4.1 Verify build and all tests pass
