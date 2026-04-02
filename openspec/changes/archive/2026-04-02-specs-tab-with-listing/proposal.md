## Why

The specs tab currently shows "Not yet implemented". Specs are the core artifact of OpenSpec projects — users need to browse them quickly. A split view with spec names on the left and rendered markdown on the right matches the existing side-panel pattern and gives immediate access to spec content.

## What Changes

- Implement the `specs` tab with a horizontal split inside the detail panel
- Left side: vertical list of spec names (directory names from `openspec/specs/`)
- Right side: rendered markdown of the selected spec's `spec.md` file
- Spec list is navigable with j/k when the specs tab is active
- Markdown rendering reuses the existing `renderMarkdown` function from the config tab

## Capabilities

### New Capabilities
- `specs-tab`: Browse and read specs in a split list+content view within the detail panel

### Modified Capabilities

## Impact

- **scanner/scan.go**: Add `SpecNames` (list of spec directory names) and `SpecContents` (map of name→content) to `ProjectInfo`
- **ui/ui.go**: New `renderSpecsTab()` with split layout, spec cursor, j/k navigation scoped to specs tab
- **No new dependencies** — reuses existing markdown renderer
