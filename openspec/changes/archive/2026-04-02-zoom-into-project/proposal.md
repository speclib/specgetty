## Why

When browsing a project's specs, changes, or config, the detail panel shares screen real estate with the project list. For projects with long specs or detailed config, users want to see more content without the project list taking up space. A "zoom" mode gives the detail panel the full terminal width.

Relates to: GitHub issue #3 — "zoom into openspec project"

## What Changes

- Press `enter` on a selected project to enter zoom mode: the project list hides and the detail panel takes the full terminal width
- Press `escape` or `enter` again to exit zoom mode and return to the split view
- All tab navigation (left/right, 1-5) and content navigation (j/k) work in zoom mode
- The nav bar updates to show the zoom/unzoom keybinding
- New `--zoom` / `-z` CLI flag: start the app in zoom mode, auto-selecting the project whose `openspec/` directory is in the current working directory (or the directory passed with `--path`)
- New `--path` / `-p` CLI flag: specify an openspec project path to zoom into directly

## Capabilities

### New Capabilities
- `zoom-mode`: Full-screen detail view that hides the project list for maximum content space

### Modified Capabilities

## Impact

- **ui/ui.go**: Add `zoomed bool` to model, update `View()` to render full-width detail panel when zoomed, update keybindings for enter/escape, update nav bar, accept initial zoom state and target path
- **src/main.go**: Add `--zoom` and `--path` CLI flags, pass to `ui.Run()`
