## Purpose
Describes the key that hands the file a pane is showing to the user's own
editor: where the key applies, how the editor is chosen from the environment,
and what the application does while the editor has the terminal.

## ADDED Requirements

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

### Requirement: The file opened is the one the pane is showing
The path handed to the editor SHALL be the file whose content is on screen,
resolved under the root the project's content was read from.

#### Scenario: A store-backed project
- **GIVEN** a project that declares a store, so the content comes from the store
  and not from the repository the user started in
- **WHEN** a file is opened from any pane
- **THEN** the path SHALL be under the store, and SHALL resolve to a file that
  exists

#### Scenario: An archived change
- **WHEN** an artifact of an archived change is opened
- **THEN** the path SHALL be the directory as it exists on disk, including the
  `YYYY-MM-DD-` prefix that the displayed name does not carry

### Requirement: The editor comes from the environment
The editor SHALL be taken from `$VISUAL` when it is set and not empty, and from
`$EDITOR` otherwise. A value carrying arguments SHALL work, the value being
split on whitespace into a command and its arguments, with the file appended
last.

#### Scenario: Both are set
- **WHEN** both `$VISUAL` and `$EDITOR` are set and a file is opened
- **THEN** `$VISUAL` SHALL be used

#### Scenario: Only EDITOR is set
- **WHEN** `$VISUAL` is unset or empty and `$EDITOR` is set
- **THEN** `$EDITOR` SHALL be used

#### Scenario: A value with arguments
- **WHEN** the chosen variable holds a command followed by arguments, such as an
  editor invoked with a wait flag
- **THEN** the command SHALL be run with those arguments, and the file path
  SHALL be the last argument

#### Scenario: Neither is set
- **WHEN** neither `$VISUAL` nor `$EDITOR` is set and the user presses `E`
- **THEN** no process SHALL be started, and the nav bar SHALL report that no
  editor is configured

#### Scenario: No editor is guessed
- **WHEN** neither variable is set
- **THEN** the application SHALL NOT fall back to any particular editor

### Requirement: The application yields the terminal while the editor runs
The application SHALL stop drawing while the editor has the terminal and resume
when the editor exits, so that the editor is not drawn over.

#### Scenario: While the editor runs
- **WHEN** an editor has been started
- **THEN** the application SHALL not draw, and the editor SHALL have the terminal
  to itself

#### Scenario: After the editor exits
- **WHEN** the editor exits
- **THEN** the application SHALL draw again at the size the terminal now has,
  and the same pane SHALL be on screen

#### Scenario: The edit is picked up
- **WHEN** the editor exits after the file was changed
- **THEN** the open project SHALL be read again, so the change is on screen
  without the user asking for it

#### Scenario: The editor could not be started
- **WHEN** the editor named by the environment cannot be run, for instance
  because no such command exists
- **THEN** the application SHALL resume, and the nav bar SHALL report why

#### Scenario: The editor exits with an error
- **WHEN** the editor runs and exits with a non-zero status
- **THEN** the application SHALL resume and report it, rather than treating the
  edit as having succeeded silently

### Requirement: The key is advertised where it applies
The nav bar SHALL show a hint for `E` exactly when the key would do something.

#### Scenario: A pane with a file
- **WHEN** a pane showing one file holds the keyboard
- **THEN** the nav bar SHALL list `E`

#### Scenario: A pane without one
- **WHEN** a pane that shows no single file is on screen
- **THEN** the nav bar SHALL NOT list `E`
