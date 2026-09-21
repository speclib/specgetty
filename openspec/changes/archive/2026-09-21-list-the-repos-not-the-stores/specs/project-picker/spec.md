## ADDED Requirements

### Requirement: Picker lists the directories you work in, with their statistics
The picker SHALL show each project as a table row with its name and statistics.
A row whose content comes from a store SHALL name that store, so that two repos
showing the same specs explain themselves rather than reading as duplicates.

#### Scenario: Row content
- **WHEN** the picker is open
- **THEN** each row SHALL show the project name, its spec count, its active
  change count, its archived change count and its aggregate task progress

#### Scenario: A row that reads from a store
- **WHEN** a row stands for a repo that declares a store
- **THEN** it SHALL name the store its content comes from, and its statistics
  SHALL be that store's

#### Scenario: Two repos sharing one store
- **WHEN** two rows declare the same store
- **THEN** both SHALL name it, and both SHALL show the same statistics

#### Scenario: A row that keeps its own content
- **WHEN** a row stands for a project holding its own specs and changes
- **THEN** it SHALL name no store

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

#### Scenario: A repo is named by its directory, never by its store
- **WHEN** a repo's directory name and the store it declares differ
- **THEN** the row SHALL show the directory name, because that is the directory
  the user works in and the one they will look for

#### Scenario: No projects
- **WHEN** no projects have been discovered
- **THEN** the picker SHALL say so and name the key that refreshes it

## MODIFIED Requirements

### Requirement: The picker opens on the row holding the open project's content
When the picker opens, the cursor SHALL rest on the row for the open project.

#### Scenario: Opening the picker from a store-backed project
- **GIVEN** a project opened from a repo that declares a store
- **WHEN** the picker is opened
- **THEN** the cursor SHALL rest on that repo's row

#### Scenario: Opening the picker from a store opened directly
- **GIVEN** a store opened by path rather than from the picker, which therefore
  has no row of its own
- **WHEN** the picker is opened
- **THEN** the picker SHALL open without failing, and no row SHALL be selected
  on the strength of a path it does not list

#### Scenario: Dismissing without choosing
- **GIVEN** the picker opened from a store-backed project
- **WHEN** the user presses `esc`
- **THEN** the view underneath SHALL be unchanged, origin included

## REMOVED Requirements

### Requirement: Picker lists projects with their statistics
**Reason**: Superseded by "Picker lists the directories you work in, with their
statistics". Three of its scenarios describe how a store row is named, and there
are no store rows any more: a store is where content lives, and the repos that
read from it are what the list is made of. The scenarios that still apply are
carried over unchanged.
