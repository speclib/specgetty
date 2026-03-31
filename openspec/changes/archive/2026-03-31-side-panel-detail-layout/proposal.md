## Why

The current layout stacks panels vertically: project list on top, file contents below. This wastes horizontal space and makes the project list feel like the primary view when it should be a navigator. A side-by-side layout with a narrow project panel on the left and a wide detail panel on the right matches the mental model better: pick a project, see its details.

## What Changes

- Restructure the TUI from vertical stacking to a horizontal split: narrow left panel (projects), wide right panel (detail)
- Left panel shows project names (directory basename) instead of full paths
- Right panel shows the project's full path and the existing openspec/ file listing as placeholder content
- Log panel remains at the bottom spanning full width when visible
- Tab switches focus between left and right panels as before

## Capabilities

### New Capabilities
- `side-panel-layout`: Horizontal split layout with project navigator on the left and detail panel on the right

### Modified Capabilities

## Impact

- **ui/ui.go**: Major restructure of `View()`, `recalcLayout()`, all height/width calculation functions, `renderRepoList()`, `renderStatusContent()`
- **No scanner changes** — data model stays the same
- **Tests**: UI layout tests need updating for new panel dimensions
