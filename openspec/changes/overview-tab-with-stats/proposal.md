## Why

After the side-panel layout is in place, the detail panel just shows a raw file listing. Users need at-a-glance project health: how many specs, active changes, open tasks. The detail panel should present a structured overview with tab navigation for future views.

## What Changes

- Add a tab header row to the detail panel: `[overview] [specs] [changes] [config] [search]`
- Only the `overview` tab is functional; others show "Not yet implemented"
- Overview tab shows: project path, stats (spec count, active changes, open tasks), active changes list, recently archived changes list
- Add OpenSpec content parsing to the scanner: count specs, list changes (active vs archived), count tasks
- Tab navigation with left/right arrow keys or number keys when detail panel is focused

## Capabilities

### New Capabilities
- `overview-tab`: Structured overview of an OpenSpec project with stats, active changes, and recently archived changes
- `detail-tabs`: Tab navigation system in the detail panel for switching between views

### Modified Capabilities

## Impact

- **scanner package**: New `ProjectInfo` struct and parsing logic to extract stats from openspec/ contents
- **ui/ui.go**: Tab header rendering, overview content rendering, tab switching keybindings
- **Dependencies**: May need to read YAML/Markdown files for change status and archive dates
