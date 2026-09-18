## MODIFIED Requirements

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

### Requirement: Rescan is scoped to the open project
Pressing `s` SHALL rescan only the open project, reading from the root its
content was resolved to.

#### Scenario: Rescan
- **WHEN** the user presses `s` in the project view
- **THEN** only the open project SHALL be reread, and the configured scan
  directories SHALL NOT be walked

#### Scenario: Rescanning a store-backed project
- **WHEN** the user presses `s` in a project resolved through a store
  declaration
- **THEN** the declaration SHALL be resolved again before reading, so that a
  repointed `store:` key is picked up
