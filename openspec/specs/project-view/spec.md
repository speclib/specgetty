# project-view Specification

## Purpose
The base view of the application: one OpenSpec project, resolved at startup from
the working directory, filling the terminal.

## Requirements

### Requirement: The project view is the startup view
The application SHALL open directly on a single project, resolved from the
working directory, without walking the configured scan directories. Where the
working directory resolves to a store, the store's content SHALL be what opens.

#### Scenario: Started inside a project
- **WHEN** the user runs `spg` from a directory that contains an `openspec/`
  subdirectory, or whose parent does
- **THEN** that project SHALL be opened immediately, and the configured scan
  directories SHALL NOT be walked

#### Scenario: Started inside a store-backed repo
- **WHEN** the user runs `spg` from a repo whose `openspec/` declares a store and
  holds no `specs/` or `changes/`
- **THEN** the store's specs and changes SHALL be opened immediately, and the
  configured scan directories SHALL NOT be walked

#### Scenario: Started outside any project
- **WHEN** the user runs `spg` from a directory with no OpenSpec project in its
  path
- **THEN** the application SHALL offer to open the project picker

#### Scenario: Declining the offer
- **WHEN** the user declines that offer
- **THEN** the application SHALL show the project view in its empty state rather
  than exiting

### Requirement: Startup view selectable by flag
The application SHALL accept a `--view` flag selecting which view opens first.
Flags that selected a change list filter mode are gone with the modes.

#### Scenario: Default
- **WHEN** no `--view` flag is given
- **THEN** the application SHALL behave as `--view=single`

#### Scenario: Opening at the picker
- **WHEN** the user runs `spg --view=all`
- **THEN** the project picker SHALL be open when the application starts

#### Scenario: Explicit path
- **WHEN** the user runs `spg --path /some/project`
- **THEN** that project SHALL be opened, as though `--view=single` had resolved
  to it

#### Scenario: The zoom flag is gone
- **WHEN** the user runs `spg --zoom`
- **THEN** the application SHALL report that the flag is not recognised

#### Scenario: The change-mode flag is gone
- **WHEN** the user runs `spg --change-mode=active`
- **THEN** the application SHALL report that the flag is not recognised

### Requirement: Empty state when no project is selected
When no project is selected, the view SHALL say so and name the key that opens
the picker.

#### Scenario: No project selected
- **WHEN** the project view is shown with no project selected
- **THEN** it SHALL display a message stating that no project is selected and
  naming the key that opens the project picker

### Requirement: Project view fills the terminal
The project view SHALL occupy the full terminal width and height, with no
project list beside it and no panel below it.

#### Scenario: Layout
- **WHEN** the project view is displayed
- **THEN** it SHALL span the full width, with the nav bar below it and nothing
  between them

#### Scenario: Terminal too small
- **WHEN** the terminal is narrower than 60 columns or shorter than 20 rows
- **THEN** the application SHALL display a message indicating the terminal is
  too small

### Requirement: Persistent project header
The project view SHALL show the project path and its statistics above the tab
bar, whichever tab is active. Where the content comes from a store, the header
SHALL carry one mark saying so, and nothing more: the store's id, its path and
its state belong to the config tab.

#### Scenario: Header always visible
- **WHEN** any tab is active
- **THEN** the project path and the statistics line SHALL be visible above the
  tab bar

#### Scenario: Statistics line content
- **WHEN** a project is open
- **THEN** the statistics line SHALL show the spec count, the active change
  count and the archived change count

#### Scenario: A project whose content comes from a store
- **WHEN** the open project was resolved through a store declaration
- **THEN** the header SHALL show the path of the repo the user started in,
  together with a mark saying the content comes from a store

#### Scenario: The store itself is open
- **WHEN** the open project is a store, selected from the picker
- **THEN** the header SHALL name the store by its id, with the same mark

#### Scenario: The mark says nothing further
- **WHEN** the header carries the store mark
- **THEN** it SHALL NOT show the store's id, root path, remote or git state,
  which the config tab carries instead

#### Scenario: A plain project
- **WHEN** the open project holds its own specs and changes
- **THEN** the header SHALL carry no store mark, and SHALL read as it does today

### Requirement: Escape does not leave the project view
The project view is the floor of the navigation stack, and `esc` SHALL NOT exit
the application or reveal any view above it.

#### Scenario: Escape at the project view
- **WHEN** the user presses `esc` in the project view with no overlay open
- **THEN** nothing SHALL happen and the application SHALL NOT exit
