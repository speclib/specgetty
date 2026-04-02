## Why

In zoom mode, users focus on a single project and manually press "s" to rescan when files change. This is tedious when actively editing openspec artifacts in an editor alongside the TUI. The TUI should automatically detect changes to the openspec directory and rescan, keeping the view up-to-date without manual intervention.

## What Changes

- Add filesystem watching (fsnotify — uses inotify on Linux, kqueue on macOS) on the zoomed project's `openspec/` directory
- Automatically trigger a single-project rescan when files within `openspec/` are created, modified, or deleted
- Start watching when entering zoom mode, stop when exiting zoom mode
- Debounce rapid file changes to avoid excessive rescans

## Capabilities

### New Capabilities
- `fs-watch`: Filesystem watching integration that monitors an openspec directory tree for changes and emits rescan messages to the bubbletea update loop

### Modified Capabilities

## Impact

- New dependency: `github.com/fsnotify/fsnotify` (cross-platform file system notifications — inotify on Linux, kqueue on macOS)
- Modified files: `src/ui/ui.go` (watch lifecycle, message handling), potentially a new `src/watcher/` package
- Scanner package unchanged — existing `doScanSingle()` is reused for rescans
