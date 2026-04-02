## Context

The TUI currently requires manual "s" key press to rescan a project's openspec directory. In zoom mode, users typically edit openspec artifacts in a separate editor and switch back to the TUI to check the result — requiring a manual rescan each time. The bubbletea architecture uses a message-based update loop, so filesystem events need to be translated into tea.Msg values.

## Goals / Non-Goals

**Goals:**
- Automatically rescan the zoomed project when its `openspec/` directory tree changes
- Debounce rapid changes (e.g., editor writing multiple files) into a single rescan
- Start/stop watching cleanly with zoom enter/exit lifecycle
- Cross-platform support (Linux via inotify, macOS via kqueue) using fsnotify

**Non-Goals:**
- Watching all projects simultaneously in non-zoom mode (too many watchers, low value)
- Watching files outside the `openspec/` directory
- Incremental/partial rescans — full `doScanSingle()` is fast enough
- Hot-reloading the UI without a full rescan cycle

## Decisions

### 1. Use `fsnotify/fsnotify` for filesystem events

**Choice**: fsnotify v1.x (stable, widely used)
**Rationale**: Standard Go library for cross-platform fs notifications. Uses inotify on Linux, kqueue on macOS. No CGO required. Alternative `rjeczalik/notify` offers recursive watching but is less maintained and heavier.

### 2. New `watcher` package

**Choice**: Create `src/watcher/watcher.go` with a `Watcher` struct
**Rationale**: Keeps filesystem concerns separate from UI logic. The watcher owns the fsnotify instance, directory walking for recursive watch setup, and debounce logic. Exposes a simple channel-based API that the UI consumes.

### 3. Channel-to-message bridge via tea.Cmd

**Choice**: Use a `tea.Cmd` that blocks on the watcher's event channel to produce a `fsChangeMsg` in the bubbletea update loop. Re-issue the cmd after each event to keep listening.
**Rationale**: This is the standard bubbletea pattern for external event sources (same as spinner ticks). No goroutine leaks — the cmd naturally stops when the channel closes.

### 4. Debounce at 200ms in the watcher

**Choice**: After receiving an fsnotify event, wait 200ms for more events before emitting a single notification on the output channel.
**Rationale**: Editors often write files in multiple steps (write temp, rename, chmod). 200ms collapses these into one rescan while still feeling instant to users. The debounce runs inside the watcher goroutine, keeping the UI simple.

### 5. Recursive watch via directory walking

**Choice**: On watcher start, walk the `openspec/` tree and add each subdirectory to fsnotify. Watch for `Create` events to add newly created subdirectories.
**Rationale**: fsnotify doesn't support recursive watching natively. The openspec tree is shallow (typically < 20 directories), so the overhead is negligible. New subdirectories (e.g., new change or spec) are picked up dynamically.

### 6. Lifecycle tied to zoom enter/exit

**Choice**: Start watcher on zoom enter, close on zoom exit. Store `*watcher.Watcher` on the model.
**Rationale**: Simple lifecycle. No watching overhead when browsing multiple projects. On zoom exit the watcher closes its channel, which naturally terminates the listening cmd.

## Risks / Trade-offs

- **[inotify limit]** → Linux has a per-user inotify watch limit (default 8192). The openspec tree is tiny (~20 dirs), so this is not a practical concern. No mitigation needed.
- **[Missed events during debounce]** → If a directory is created and immediately populated, we might process the rescan before all files are written. → The 200ms debounce window handles most cases; worst case the next save triggers another rescan.
- **[Watcher goroutine leak]** → If zoom exit forgets to close the watcher. → The `Close()` method is called in both zoom-exit paths (Escape and Enter). Defensive: watcher checks for closed channel.
