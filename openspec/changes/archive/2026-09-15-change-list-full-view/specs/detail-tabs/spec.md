## MODIFIED Requirements

### Requirement: Tab header in detail panel
The detail panel SHALL display a tab header row showing the available views,
which no longer include a separate archive tab.

#### Scenario: Tab header rendering
- **WHEN** the detail panel is displayed
- **THEN** a header row SHALL show tabs: changes, specs, config, with the active
  tab visually highlighted

### Requirement: Tab switching
The user SHALL be able to switch between tabs when the detail panel is focused
and no search prompt is open.

#### Scenario: Switch tab with number keys
- **WHEN** the detail panel is focused and the user presses 1-3
- **THEN** the corresponding tab SHALL become active (1=changes, 2=specs,
  3=config)

#### Scenario: Switch tab with arrow keys
- **WHEN** the detail panel is focused with the change list displayed and the user
  presses left or right
- **THEN** the active tab SHALL move left or right respectively

#### Scenario: Switch tab with h/l keys
- **WHEN** the detail panel is focused and user presses h or l
- **THEN** the active tab SHALL move left or right respectively

#### Scenario: Number keys while a change is open
- **WHEN** a change is open and the user presses a number key
- **THEN** the project tab bar SHALL NOT change
