# detail-tabs Specification

## Purpose
TBD - created by archiving change overview-tab-with-stats. Update Purpose after archive.
## Requirements
### Requirement: Tab header in detail panel
The detail panel SHALL display a tab header row showing available views.

#### Scenario: Tab header rendering
- **WHEN** the detail panel is displayed
- **THEN** a header row SHALL show tabs: overview, specs, changes, config, search — with the active tab visually highlighted

### Requirement: Tab switching
The user SHALL be able to switch between tabs when the detail panel is focused.

#### Scenario: Switch tab with number keys
- **WHEN** the detail panel is focused and user presses 1-5
- **THEN** the corresponding tab SHALL become active (1=overview, 2=specs, 3=changes, 4=config, 5=search)

#### Scenario: Switch tab with h/l keys
- **WHEN** the detail panel is focused and user presses h or l
- **THEN** the active tab SHALL move left or right respectively

### Requirement: Unimplemented tabs show placeholder
Tabs other than overview SHALL display a "Not yet implemented" message.

#### Scenario: Selecting an unimplemented tab
- **WHEN** the user switches to the specs, changes, config, or search tab
- **THEN** the detail panel content SHALL display "Not yet implemented" centered in the panel

