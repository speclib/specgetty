## Context

Specgetty is a Bubbletea TUI that browses OpenSpec projects. It currently has no external command execution — all data is read-only via filesystem scanning. This change introduces the first write action: archiving a change by shelling out to `openspec archive`.

The existing modal pattern (scanning spinner, error overlay) provides a foundation for confirmation and feedback UI. The existing `doScanSingle` command provides the rescan-after-mutation pattern.

## Goals / Non-Goals

**Goals:**
- Archive a selected change from the changes tab with a single keybinding
- Warn users when incomplete tasks exist before archiving
- Provide clear feedback on success, failure, or missing CLI
- Rescan the project after successful archive so the UI updates

**Non-Goals:**
- Un-archiving changes from the archive tab
- Configuring `openspec archive` flags (e.g. `--skip-specs`) from the UI
- Streaming live output during archive execution

## Decisions

**Background exec with modal feedback (not `tea.Exec`)**
Use `exec.Command` in a `tea.Cmd` goroutine and capture combined output. Display the result in a modal overlay. This keeps the TUI visible throughout and avoids the terminal handoff complexity of `tea.Exec`. The `openspec archive -y` command is fast (sub-second) and non-interactive, so live output streaming is unnecessary.

Alternative: `tea.Exec` temporarily hands the terminal to the child process. Rejected because it causes a visual flash/redraw for a sub-second operation and adds complexity for no user benefit.

**Pass `-y` flag to skip openspec's own confirmation**
The TUI handles its own confirmation modal (with task warnings), so `openspec archive`'s interactive prompt must be bypassed. The `-y` flag does this.

**Confirmation modal with two states**
A single modal type that conditionally shows a task warning. When there are incomplete tasks, the modal lists them and asks "Archive anyway?". When all tasks are done (or there are no tasks), a simple "Archive <name>?" confirmation is shown. This avoids a multi-step wizard for a simple action.

**State machine for the archive flow**
Add an `archiveState` field to the model with states: `idle`, `confirming`, `running`, `result`. This keeps the flow explicit and testable. The result state holds success/failure message and auto-clears on any keypress.

```
idle ──[a key]──▶ confirming ──[y]──▶ running ──[done]──▶ result ──[any key]──▶ idle
                      │                                        │
                      └──[n/esc]──▶ idle                       └──[if success]──▶ rescan
```

**Detect `openspec` via `exec.LookPath` at archive time**
Check on each archive attempt rather than at startup. This handles the case where the user installs openspec while specgetty is running. Show a specific modal message when not found.

**Run `openspec archive` from the project directory**
Set `cmd.Dir` to the project path so openspec resolves the change correctly.

## Risks / Trade-offs

- [Risk] `openspec archive` modifies the filesystem while the TUI is running → Mitigation: rescan immediately after success to sync state; the archive operation is atomic from the user's perspective.
- [Risk] Long-running archive (e.g. large spec merge) could block the UI → Mitigation: the command runs in a `tea.Cmd` goroutine so the UI stays responsive; a "Running..." modal is shown.
- [Trade-off] Using `-y` means we skip openspec's own validation warnings → Acceptable because we show our own task warning, and openspec's validation errors still surface as command failures in the result modal.
