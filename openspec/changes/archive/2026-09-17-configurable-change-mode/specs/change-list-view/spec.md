## MODIFIED Requirements

### Requirement: Active and archived changes share one list
Active changes and archived changes SHALL be shown in a single list with a
filter selecting which are included. The three states SHALL be named `active`,
`archived` and `active+archived` wherever the user meets them, whether on screen
or in configuration.

#### Scenario: Filter modes
- **WHEN** the change list is displayed
- **THEN** the filter SHALL be in one of three modes: active only, archived
  only, or both, with the current mode visible

#### Scenario: Default mode
- **WHEN** the change list is first displayed for a project
- **THEN** the filter SHALL be in the configured default mode, which is active
  only when nothing is configured

#### Scenario: Default mode from configuration
- **WHEN** a default mode is set in the configuration file
- **THEN** the change list SHALL open in that mode

#### Scenario: Default mode from the command line
- **WHEN** a default mode is given on the command line
- **AND** one is also set in the configuration file
- **THEN** the command line value SHALL be used

#### Scenario: An unrecognised mode
- **WHEN** a configured or supplied mode is not one of the three names
- **THEN** the application SHALL report the unknown value together with the
  valid ones, rather than falling back silently

#### Scenario: Returning to the default
- **WHEN** the user has changed the mode and then selects a different project
- **THEN** the filter SHALL return to the configured default

#### Scenario: Both mode
- **WHEN** the filter is in both mode
- **THEN** active and archived changes SHALL appear in the same list, and each
  row SHALL indicate whether it is active or archived

#### Scenario: No changes in the selected mode
- **WHEN** the filter mode selects no changes
- **THEN** the list SHALL state which mode is active and that it contains no
  changes
