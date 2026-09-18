## MODIFIED Requirements

### Requirement: Picker lists projects with their statistics
The picker SHALL show each project as a table row with its name and statistics.
A row for a store SHALL be named by the store's id and SHALL be marked as a
store, so that it is not read as a project that happens to be called after its
folder.

#### Scenario: Row content
- **WHEN** the picker is open
- **THEN** each row SHALL show the project name, its spec count, its active
  change count, its archived change count and its aggregate task progress

#### Scenario: A store row
- **WHEN** a row stands for a store
- **THEN** it SHALL be named by the id in that store's `.openspec-store/store.yaml`
  and SHALL carry a mark distinguishing it from a plain project

#### Scenario: A store whose id differs from its folder
- **WHEN** a store's directory name and its declared id differ
- **THEN** the row SHALL show the declared id

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

#### Scenario: A store id colliding with a project name
- **WHEN** a store's id equals another row's name
- **THEN** the two SHALL be disambiguated the same way duplicate basenames are

#### Scenario: No projects
- **WHEN** no projects have been discovered
- **THEN** the picker SHALL say so and name the key that refreshes it

## ADDED Requirements

### Requirement: The picker opens on the row holding the open project's content
When the picker opens, the cursor SHALL rest on the row whose root holds the
content currently on screen, which for a store-backed project is the store's
row rather than any row for the repo the user started in.

#### Scenario: Opening the picker from a store-backed project
- **GIVEN** a project opened from a repo that points at a store
- **WHEN** the picker is opened
- **THEN** the cursor SHALL rest on that store's row

#### Scenario: Dismissing without choosing
- **GIVEN** the picker opened from a store-backed project
- **WHEN** the user presses `esc`
- **THEN** the view underneath SHALL be unchanged, origin included
