<!-- Clauses written bare, with no list marker. Taken from
     /home/pim/mipnix/openspec/specs/airplane-mode/spec.md, the file bean
     specgetty-vfnn reported as showing empty scenarios. 157 lines in the
     local corpus are written this way. The delta header of the original is
     removed here so this fixture isolates the clause shape. -->
# Airplane Mode

## Purpose
Controls the radios of the machine from one toggle in the panel.

## Requirements

### Requirement: airplane-mode-toggle
The panel SHALL display a toggle button for airplane mode.

#### Scenario: enable airplane mode

GIVEN airplane mode is off
WHEN the user clicks the airplane mode toggle
THEN all radios SHALL be blocked via `rfkill block all`

#### Scenario: toggle reflects state

WHEN the panel opens
THEN the airplane mode toggle SHALL reflect whether any radios are blocked
