## Why

The `config` tab in the detail panel currently shows "Not yet implemented". Users need to quickly see a project's configuration — either `project.md` (rendered as parsed markdown) or `config.yaml` (with YAML syntax highlighting). This gives at-a-glance understanding of a project's setup without opening files externally.

## What Changes

- Implement the `config` tab to read and display the project's OpenSpec configuration file
- If `openspec/project.md` exists, render it as parsed markdown (headers, lists, bold/italic via lipgloss)
- If `openspec/config.yaml` exists (and no project.md), render it with YAML syntax highlighting (keys, values, comments in different colors)
- If neither exists, show "No project configuration found"

## Capabilities

### New Capabilities
- `config-tab-display`: Display project configuration in the config tab with format-appropriate rendering

### Modified Capabilities

## Impact

- **scanner/scan.go**: Add `ConfigFile` field to `ProjectInfo` indicating which config file exists (project.md or config.yaml)
- **ui/ui.go**: Implement `renderConfigTab()` with markdown parsing and YAML highlighting
- **No new dependencies** — use lipgloss styling for both markdown and YAML rendering
