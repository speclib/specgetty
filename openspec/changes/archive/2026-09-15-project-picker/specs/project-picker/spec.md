## Purpose
An overlay listing every discovered OpenSpec project with its statistics, so the
user can find one and switch to it from anywhere in the application.

## ADDED Requirements

### Requirement: Opening and closing the picker
The picker SHALL be opened with a key from anywhere in the application and
dismissed without changing the current project.

#### Scenario: Open from the project view
- **WHEN** the user presses the picker key in the project view
- **THEN** the picker SHALL open as an overlay drawn over the current view

#### Scenario: Open while a change is open
- **WHEN** the user presses the picker key while a change is open
- **THEN** the picker SHALL open as an overlay

#### Scenario: Dismiss
- **WHEN** the picker is open and the user presses `esc`
- **THEN** the picker SHALL close, the current project SHALL be unchanged, and
  the view underneath SHALL be as it was

### Requirement: Selecting a project
Selecting a project SHALL switch to it and land on the change list.

#### Scenario: Select from the project view
- **WHEN** a project is under the cursor and the user presses `enter`
- **THEN** the picker SHALL close and that project SHALL be open at the project
  view with the changes tab active

#### Scenario: Select while a change was open
- **WHEN** the picker was opened while a change was open, and the user selects a
  different project
- **THEN** the application SHALL land on the project view of the selected
  project rather than attempting to reopen a change

#### Scenario: Filesystem watching follows the selection
- **WHEN** the user switches to a different project
- **THEN** watching SHALL stop for the previous project and start for the newly
  selected one

### Requirement: Picker lists projects with their statistics
The picker SHALL show each project as a table row with its name and statistics.

#### Scenario: Row content
- **WHEN** the picker is open
- **THEN** each row SHALL show the project name, its spec count, its active
  change count, its archived change count and its aggregate task progress

#### Scenario: Statistics are current
- **WHEN** the picker opens
- **THEN** the statistics shown SHALL be read from disk at that moment rather
  than from the cache

#### Scenario: Project names as basenames
- **WHEN** all discovered projects have unique directory basenames
- **THEN** each row SHALL show only the basename

#### Scenario: Duplicate basenames
- **WHEN** two or more projects share a directory basename
- **THEN** the parent directory name SHALL be appended to disambiguate, for
  example `specgetty (cVibeCoding)`

#### Scenario: No projects
- **WHEN** no projects have been discovered
- **THEN** the picker SHALL say so and name the key that refreshes it

### Requirement: Picker filter
The picker SHALL support filtering with the same query grammar as the change
list.

#### Scenario: Open the filter
- **WHEN** the picker is open and the user presses `/`
- **THEN** a prompt SHALL appear and subsequent letter keys SHALL be entered as
  query text

#### Scenario: Fuzzy match on project names
- **WHEN** the query has no sigil
- **THEN** it SHALL be matched as a fuzzy subsequence against project names,
  ordered with the strongest match first

#### Scenario: Literal match on project names
- **WHEN** the query begins with `'`
- **THEN** the rest SHALL be matched as a literal substring of the project name

#### Scenario: Match inside projects
- **WHEN** the query begins with `:`
- **THEN** the rest SHALL be matched as a literal substring against both the file
  paths and the file contents within each project

#### Scenario: Smart case
- **WHEN** the query contains no uppercase character
- **THEN** matching SHALL be case-insensitive, and case-sensitive otherwise

#### Scenario: Nothing matches
- **WHEN** a non-empty query matches no project
- **THEN** the picker SHALL say so and echo the query

### Requirement: Refreshing the picker
The picker SHALL offer an action that rediscovers projects from disk.

#### Scenario: Refresh
- **WHEN** the user presses `r` in the picker
- **THEN** the configured scan directories SHALL be walked again, the cache
  SHALL be rewritten, and the list SHALL show what was found

#### Scenario: Progress during a refresh
- **WHEN** a refresh is running
- **THEN** the picker SHALL indicate that it is working

### Requirement: Picker keys outrank the view beneath it
While the picker is open its keys SHALL take precedence over the view it covers,
and be outranked by a confirmation modal.

#### Scenario: Keys belong to the picker
- **WHEN** the picker is open and the user presses a key bound in the change
  list, such as `a` or `f`
- **THEN** the change list action SHALL NOT run

#### Scenario: A confirmation modal outranks the picker
- **WHEN** a confirmation modal is awaiting an answer
- **THEN** the picker key SHALL NOT open the picker
