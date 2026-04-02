## 1. Model and state

- [x] 1.1 Add `zoomed bool` field to model
- [x] 1.2 Add `enter` keybinding: when project list focused, set `zoomed = true` and switch focus to detail panel
- [x] 1.3 Add `escape` keybinding: when zoomed, set `zoomed = false` and switch focus to project list
- [x] 1.4 Add `enter` keybinding when zoomed: exit zoom (same as escape)
- [x] 1.5 Disable `tab` key when zoomed

## 2. Scoped scanning

- [x] 2.1 Add `doScanSingle(projectPath string) tea.Cmd` that rescans only one project (ListOpenSpecContents + ParseProjectInfo)
- [x] 2.2 When `zoomed` and `s` is pressed, call `doScanSingle` with the zoomed project path instead of `doScan`
- [x] 2.3 When not zoomed and `s` is pressed, call `doScan` as before (full scan)

## 3. Layout rendering

- [x] 3.1 Update `View()`: when `zoomed`, render only the detail panel at full width (skip left panel and horizontal join)
- [x] 3.2 Update detail panel width calculation: when zoomed, use `m.width - 2` instead of `rightPanelWidth()`
- [x] 3.3 Update `renderPanel` title for detail: when zoomed, show "Detail (<project name>)" using display name

## 4. Nav bar and polish

- [x] 4.1 Update nav bar: show `enter zoom` when not zoomed (project list focused), show `esc back` when zoomed

## 5. CLI flags

- [x] 5.1 Add `--zoom` / `-z` and `--path` / `-p` CLI flags to main.go
- [x] 5.2 Add `findOpenSpecProject(startDir string) string` that walks up from startDir looking for openspec/ directory
- [x] 5.3 Update `ui.Run()` signature to accept `initialZoomPath string`
- [x] 5.4 When `initialZoomPath` is set, use `doScanSingle` as the initial scan instead of full `doScan`; set cursor to 0 and enter zoom mode
- [x] 5.5 If `--zoom` without `--path`, resolve from cwd; if `--path` given, use that directly
- [x] 5.6 When exiting zoom mode after a `--zoom` start, first full scan has not run yet — trigger full `doScan` on unzoom so the project list populates

## 6. Tests and verification

- [x] 6.1 Add test for `findOpenSpecProject` path resolution
- [x] 6.2 Verify build and all tests pass
