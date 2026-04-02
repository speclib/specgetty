## 1. Dependencies and Package Setup

- [x] 1.1 Add `github.com/fsnotify/fsnotify` dependency via `go get`
- [x] 1.2 Create `src/watcher/watcher.go` with `Watcher` struct, `New()`, `Close()`, and `Events()` channel

## 2. Core Watcher Implementation

- [x] 2.1 Implement recursive directory walking to add all subdirs under a given `openspec/` path to fsnotify
- [x] 2.2 Implement debounce logic (200ms) — collect fsnotify events and emit a single notification on the output channel after the debounce window
- [x] 2.3 Handle `Create` events for new subdirectories — dynamically add them to the watch set
- [x] 2.4 Implement `Close()` to stop the watcher goroutine and close the events channel cleanly

## 3. UI Integration

- [x] 3.1 Add `fsChangeMsg` type and a `tea.Cmd` that listens on the watcher's events channel
- [x] 3.2 Handle `fsChangeMsg` in the Update loop — trigger `doScanSingle()` on the zoomed project
- [x] 3.3 Start watcher on zoom enter (both keyboard Enter and `--zoom` CLI flag) with the project's openspec path
- [x] 3.4 Stop watcher on zoom exit (both Escape and Enter exit paths)
- [x] 3.5 Store `*watcher.Watcher` on the model and nil-check before close

## 4. Testing

- [x] 4.1 Write unit tests for the watcher package — verify events are emitted on file changes and debounce collapses rapid changes
- [x] 4.2 Verify watcher starts/stops correctly by entering and exiting zoom mode manually
- [x] 4.3 Test new subdirectory detection — create a directory while watching and verify it gets watched
