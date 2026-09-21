<!-- A main spec carrying a delta header. 102 of the 544 live specs on this
     machine do. openspec calls it an error: it "truncates the parsed
     ## Requirements section", so validate, list and archive see nothing here.
     Taken from /home/pim/mipnix/openspec/specs/airplane-mode/spec.md. -->
# Airplane Mode

## ADDED Requirements

### Requirement: airplane-mode-toggle
The panel SHALL display a toggle button for airplane mode.

#### Scenario: enable airplane mode
- **WHEN** the user clicks the toggle
- **THEN** all radios SHALL be blocked
