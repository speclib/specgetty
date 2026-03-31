## 1. Scanner: project info parsing

- [x] 1.1 Add `ProjectInfo` struct to scanner with fields: `SpecCount int`, `ActiveChanges []string`, `ArchivedChanges []ArchivedChange` (name + date)
- [x] 1.2 Add `ParseProjectInfo(dir string) ProjectInfo` that reads openspec/specs/, openspec/changes/, openspec/archive/ to populate stats
- [x] 1.3 Add `ProjectInfo` field to `ProjectStatus` and call `ParseProjectInfo` during `Scan()`
- [x] 1.4 Add tests for `ParseProjectInfo` with various directory structures

## 2. Tab system

- [x] 2.1 Add `detailTab` constants (tabOverview, tabSpecs, tabChanges, tabConfig, tabSearch) and field to model
- [x] 2.2 Add `renderTabHeader()` method that renders the tab bar with active tab highlighted
- [x] 2.3 Add tab switching keybindings: h/l for prev/next, 1-5 for direct selection (only when detail panel focused)
- [x] 2.4 Update nav bar to show tab switching keys

## 3. Overview tab content

- [x] 3.1 Create `renderOverview()` that displays: project path, stats line, active changes list, recently archived list
- [x] 3.2 Create placeholder `renderNotImplemented()` for other tabs
- [x] 3.3 Wire tab content into detail panel: switch on `detailTab` to call the right renderer

## 4. Tests and verification

- [x] 4.1 Add test for `renderTabHeader` active tab highlighting
- [x] 4.2 Verify build and all tests pass
