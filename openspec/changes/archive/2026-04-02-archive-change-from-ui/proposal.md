## Why

Users currently need to leave specgetty and run `openspec archive` from the command line to archive a completed change. Archiving directly from the TUI while browsing changes removes that context switch. This requires the `openspec` CLI to be installed; clear feedback must be given on success, failure, or when the CLI is missing. Users should be warned when archiving a change that still has open tasks.

## What Changes

- Add `a` keybinding in the changes tab to archive the currently selected change
- Detect whether the `openspec` CLI is available on PATH
- Show a confirmation modal before archiving, with a warning when tasks are incomplete
- Execute `openspec archive <name> -y` in the background and display the result in a feedback modal
- Rescan the project after a successful archive so the changes and archive tabs update
- Add `a archive` hint to the nav bar when the changes tab is active

## Capabilities

### New Capabilities

- `archive-from-ui`: Archiving an OpenSpec change from within the TUI, including confirmation flow, task warnings, command execution, and result feedback

### Modified Capabilities

## Impact

- `src/ui/ui.go`: New keybinding, confirmation modal state, exec integration, feedback modal, nav bar update
- New dependency: `os/exec` (first usage in the codebase)
