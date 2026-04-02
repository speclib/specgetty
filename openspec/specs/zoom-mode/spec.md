# zoom-mode Specification

## Purpose
TBD - created by archiving change zoom-into-project. Update Purpose after archive.
## Requirements
### Requirement: Enter zoom mode
The user SHALL be able to enter zoom mode to view the detail panel at full terminal width.

#### Scenario: Press enter on project list
- **WHEN** the project list is focused and the user presses `enter`
- **THEN** the project list SHALL be hidden, the detail panel SHALL expand to full width, and focus SHALL move to the detail panel

### Requirement: Exit zoom mode
The user SHALL be able to exit zoom mode to return to the split view.

#### Scenario: Press escape while zoomed
- **WHEN** zoom mode is active and the user presses `escape`
- **THEN** the project list SHALL reappear, the layout SHALL return to the split view, and focus SHALL return to the project list

#### Scenario: Press enter while zoomed
- **WHEN** zoom mode is active and the user presses `enter`
- **THEN** zoom mode SHALL be exited (same as pressing escape)

### Requirement: Full functionality in zoom mode
All detail panel features SHALL work in zoom mode.

#### Scenario: Tab navigation while zoomed
- **WHEN** zoom mode is active
- **THEN** left/right arrows and 1-5 number keys SHALL switch tabs, and j/k SHALL navigate within tabs

#### Scenario: Tab key disabled while zoomed
- **WHEN** zoom mode is active and the user presses `tab`
- **THEN** nothing SHALL happen (no panel to switch to)

### Requirement: Zoom mode shows project context
The detail panel SHALL indicate which project is being viewed when zoomed.

#### Scenario: Panel title in zoom mode
- **WHEN** zoom mode is active
- **THEN** the detail panel title SHALL include the project display name (e.g. "Detail (specgetty)")

### Requirement: Scanning scoped to zoomed project
When in zoom mode, rescanning SHALL only scan the zoomed project, not all configured directories.

#### Scenario: Rescan while zoomed
- **WHEN** zoom mode is active and the user presses `s`
- **THEN** only the zoomed project SHALL be rescanned (not the full scan directory list)

#### Scenario: Rescan after unzoom
- **WHEN** the user exits zoom mode and presses `s`
- **THEN** a full scan of all configured directories SHALL be performed

### Requirement: Start zoomed via CLI flag
The application SHALL support starting in zoom mode via command-line flags.

#### Scenario: --zoom flag from project directory
- **WHEN** the user runs `spg --zoom` from a directory containing an `openspec/` subdirectory (or a parent does)
- **THEN** the app SHALL scan only that single project (not all configured scan directories) and start in zoom mode with it selected

#### Scenario: --zoom with --path flag
- **WHEN** the user runs `spg --zoom --path /some/project`
- **THEN** the app SHALL start in zoom mode with the specified project selected

#### Scenario: --zoom in directory without openspec
- **WHEN** the user runs `spg --zoom` from a directory with no openspec project in its path
- **THEN** the app SHALL start in normal mode (no zoom)

