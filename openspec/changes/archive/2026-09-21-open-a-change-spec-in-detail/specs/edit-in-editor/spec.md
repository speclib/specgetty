## MODIFIED Requirements

### Requirement: E opens the file the pane is showing
A pane that shows exactly one file on disk SHALL offer `E`, which opens that
file in the user's editor. A pane that shows no file, or that shows several
files at once, SHALL NOT act on the key.

#### Scenario: A change's artifact
- **WHEN** a change's proposal, design or tasks artifact is on screen and the
  user presses `E`
- **THEN** that artifact's file SHALL be opened in the editor

#### Scenario: A spec
- **WHEN** a spec is selected on the specs tab and the user presses `E`
- **THEN** that spec's `spec.md` SHALL be opened in the editor

#### Scenario: The project configuration
- **WHEN** the project row of the properties tab is selected and the user
  presses `E`
- **THEN** the configuration file that row reports SHALL be opened in the editor

#### Scenario: A node of a change's delta
- **WHEN** a node is selected in the change spec detail view and the user presses
  `E`
- **THEN** the delta file that node belongs to SHALL be opened in the editor,
  that view having narrowed the several files of the specs sub-tab down to one

#### Scenario: A pane showing several files
- **WHEN** a change's specs sub-tab is on screen, which renders every spec delta
  of the change together, and the user presses `E`
- **THEN** nothing SHALL happen, because there is no single file to open

#### Scenario: A pane showing no file
- **WHEN** a properties row that reports on a schema or on the store is selected,
  those being assembled reports rather than files, and the user presses `E`
- **THEN** nothing SHALL happen

#### Scenario: An overlay holds the keyboard
- **WHEN** the project picker, a confirmation or the search prompt holds the
  keyboard and the user presses `E`
- **THEN** the key SHALL belong to the overlay and no editor SHALL be opened

#### Scenario: The file opens at its top
- **WHEN** a file is opened from any pane
- **THEN** it SHALL be opened at its first line, and no position within the file
  SHALL be passed to the editor
