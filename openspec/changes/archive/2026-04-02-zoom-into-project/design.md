## Context

The current layout is a horizontal split: narrow project list on the left (~30%), detail panel on the right (~70%). When reading long specs or config files, the detail panel can feel cramped.

```
Normal mode:                          Zoom mode:
┌─ Projects ──┐┌─ Detail ────────┐   ┌─ Detail (specgetty) ──────────────────┐
│ specgetty    ││ [overview] ...  │   │ [overview] [specs] [changes] ...      │
│ other-proj   ││ ...content...   │   │                                       │
│              ││                 │   │ ...full-width content...               │
└──────────────┘└─────────────────┘   │                                       │
                                      └───────────────────────────────────────┘
```

## Goals / Non-Goals

**Goals:**
- Toggle zoom with `enter` (zoom in) and `escape` or `enter` (zoom out)
- Detail panel uses full terminal width minus borders in zoom mode
- All existing detail panel functionality works: tabs, navigation, content rendering
- Panel title shows project name when zoomed for context
- `tab` key is disabled in zoom mode (only detail panel visible, no panel switching)

**Non-Goals:**
- Zooming into individual specs or changes (just the project detail level)
- Remembering zoom state across project switches

## Decisions

### 1. Simple boolean toggle on model
Add `zoomed bool` to the model. `View()` checks this flag to decide whether to render split or full-width.

### 2. Enter to zoom, escape to unzoom
`enter` when project list is focused enters zoom mode and switches focus to detail. `escape` exits zoom mode. `enter` while zoomed also exits.

### 3. Panel title includes project name when zoomed
Since the project list is hidden, the detail panel title changes from "Detail" to the project display name (e.g. "Detail (specgetty)") so users know which project they're viewing.

### 4. Reset zoom on project list focus
If somehow focus returns to the project list, zoom is automatically cancelled.

### 5. CLI flags: --zoom and --path
`--zoom` / `-z` starts the app in zoom mode. Without `--path`, it checks if the current working directory (or a parent) contains an `openspec/` directory and selects that project. With `--path` / `-p`, it zooms into the specified project path directly.

The `ui.Run()` signature gains `initialZoom string` — empty means normal mode, non-empty is the project path to zoom into. `main.go` resolves the path from flags before calling Run.

### 6. Scanning in zoom mode
When zoomed (whether via CLI flag or enter key), pressing `s` to rescan only rescans the single zoomed project, not all configured scan directories. This makes rescan near-instant and avoids unnecessary filesystem traversal.

Implementation: when `zoomed`, `doScan()` skips the full `Walk()` and instead calls `ListOpenSpecContents` + `ParseProjectInfo` directly on the zoomed project path.

### 7. Skip full scan when starting with --zoom
When `--zoom` is used, the app should not do a full directory walk at all. Instead, it directly scans just the single target project (resolved from cwd or `--path`). The full scan directories from config are ignored on startup. If the user later exits zoom mode and presses `s`, the full scan runs then.

### 8. Path resolution for --zoom without --path
Walk up from `cwd` looking for an `openspec/` directory. If found, use that directory's parent as the project path. If the project isn't in the scan results, scan the parent directory directly.

## Risks / Trade-offs

- **Trade-off**: Using `escape` for unzoom means it can't be used for other purposes in the detail panel later. Acceptable since `escape` is a natural "go back" key.
- **Risk**: `--zoom` without `--path` in a directory with no openspec project → Mitigation: Fall back to normal (non-zoomed) mode with a log message.
