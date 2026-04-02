## 1. Model: archive state machine

- [x] 1.1 Add `archiveState int` constants: `archiveIdle`, `archiveConfirming`, `archiveRunning`, `archiveResult`
- [x] 1.2 Add fields to model: `archiveState int`, `archiveChangeName string`, `archiveResultMsg string`, `archiveResultOk bool`

## 2. Archive command execution

- [x] 2.1 Add `doArchiveChange(projectPath string, changeName string) tea.Cmd` that checks `exec.LookPath("openspec")`, runs `openspec archive <name> -y` with `cmd.Dir` set to projectPath, captures combined output, and returns a result message
- [x] 2.2 Add `archiveResultMsg` message type for the tea.Cmd return value (success bool, output string)

## 3. Keybinding and Update logic

- [x] 3.1 Handle `a` key: when `detailTab == tabChanges` and changes exist, set `archiveState = archiveConfirming` with the selected change name
- [x] 3.2 Handle `y` key in `archiveConfirming`: transition to `archiveRunning`, call `doArchiveChange`
- [x] 3.3 Handle `n`/`escape` key in `archiveConfirming`: transition to `archiveIdle`
- [x] 3.4 Handle `archiveResultMsg` in Update: set `archiveState = archiveResult`, store result; if success, trigger `doScanSingle` rescan
- [x] 3.5 Handle any key in `archiveResult`: transition to `archiveIdle`
- [x] 3.6 Suppress normal keybindings while `archiveState != archiveIdle`

## 4. Modal rendering

- [x] 4.1 Add confirmation modal: show change name, task warning (with incomplete count) when `TasksDone < TasksTotal`, and "Archive? (y/n)" prompt
- [x] 4.2 Add "openspec not installed" modal when LookPath fails (returned as archiveResult with ok=false)
- [x] 4.3 Add result modal: show success or failure message with command output, "press any key to dismiss"
- [x] 4.4 Wire modals into `View()` as overlays when `archiveState != archiveIdle`

## 5. Nav bar

- [x] 5.1 Add `a archive` hint to nav bar when changes tab is active and project has active changes

## 6. Tests and verification

- [x] 6.1 Add test: `a` key sets archiveConfirming state when on changes tab with changes
- [x] 6.2 Add test: `y` in confirming transitions to archiveRunning
- [x] 6.3 Add test: `n` in confirming transitions to archiveIdle
- [x] 6.4 Verify build and all tests pass
