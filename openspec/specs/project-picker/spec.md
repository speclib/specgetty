# project-picker Specification

## Purpose
An overlay listing every discovered OpenSpec project with its statistics, so the
user can find one and switch to it from anywhere in the application.

## Requirements

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

### Requirement: Picker columns are separated by two blank columns
The project picker draws its rows with the same table renderer as the change
list, and its columns SHALL be separated by the same two blank columns.

#### Scenario: A project name that fills its column
- **WHEN** a project name is as wide as its column allows
- **THEN** two blank columns SHALL separate it from the spec count beside it

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

### Requirement: The picker can be paged
The project picker SHALL move its cursor by a page and by half a page, beside
the jumps to either end it already has.

#### Scenario: Paging the picker
- **WHEN** the picker holds the keyboard and the user presses `pgdown`, `pgup`,
  `ctrl+f`, `ctrl+b`, `ctrl+d` or `ctrl+u`
- **THEN** the selected project SHALL move by a page or half a page

#### Scenario: The jumps are unchanged
- **WHEN** the user presses `g` or `G`
- **THEN** the selection SHALL move to the first or last project, as it already
  does

#### Scenario: Within the filter
- **WHEN** a filter is narrowing the picker
- **THEN** these keys SHALL move within what the filter left

### Requirement: The picker names its matchers and suggests the contents search
The picker's prompt SHALL name the three matchers while it is focused and empty,
and a name search that matched no project SHALL suggest the same term as a
contents search.

The picker's contents search reaches the file paths and the file contents of
every project found on the machine, which is the least guessable capability in
the application.

#### Scenario: The picker's prompt is opened
- **WHEN** the picker's prompt is focused and its query is empty
- **THEN** it SHALL name the default matcher, the literal-name sigil and the
  contents sigil

#### Scenario: A failed name search in the picker
- **WHEN** a query with no sigil, or one beginning with `'`, matches no project
- **THEN** the message SHALL say so and SHALL suggest the same term prefixed
  with the contents sigil

#### Scenario: A failed contents search in the picker
- **WHEN** a query beginning with `:` matches no project
- **THEN** the message SHALL say so and SHALL suggest nothing

#### Scenario: One wording for two meanings
- **WHEN** the naming is shown in either the picker or the change list
- **THEN** it SHALL use the same words, even though the contents sigil reaches
  artifact and spec text in a change and file paths and contents in a project
