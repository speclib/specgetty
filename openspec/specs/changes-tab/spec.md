# changes-tab Specification

## Purpose
TBD - created by archiving change changes-tab-with-artifacts. Update Purpose after archive.

## Requirements

### Requirement: Changes tab lists active changes with task progress
The changes tab SHALL display the change list across the full panel width as a
table, with the columns selected by configuration.

#### Scenario: Project with active changes
- **WHEN** the user switches to the changes tab for a project with changes
- **THEN** the full panel width SHALL show a table of changes sorted
  alphabetically, with a column header row

#### Scenario: Change with no tasks.md
- **WHEN** a change directory has no tasks.md file
- **THEN** the change SHALL be listed with an empty task progress column

#### Scenario: No active changes
- **WHEN** the project has no directories under `openspec/changes/` and the
  filter is in active-only mode
- **THEN** the changes tab SHALL display "No active changes"

### Requirement: Changes tab shows artifact sub-tabs
When a change is opened with `enter`, the change SHALL fill the panel and show
sub-tabs for each discovered artifact file plus a specs sub-tab.

#### Scenario: Change with standard artifacts
- **WHEN** a change is opened that contains proposal.md, design.md, and tasks.md
- **THEN** sub-tabs SHALL show: `proposal`, `design`, `tasks`, `specs`

#### Scenario: Change with custom artifacts
- **WHEN** a change contains additional .md files (e.g. notes.md, research.md)
- **THEN** those files SHALL appear as additional sub-tabs

#### Scenario: Sub-tab navigation
- **WHEN** a change is open and the user presses left/right arrows
- **THEN** the artifact sub-tab SHALL change and the content SHALL update,
  without affecting the project tab bar

### Requirement: Artifact content rendered as markdown
Selected artifact content SHALL be rendered with markdown styling.

#### Scenario: Viewing proposal.md
- **WHEN** the proposal sub-tab is selected
- **THEN** the content of proposal.md SHALL be rendered with headers, lists, bold, and italic styled

### Requirement: Tasks artifact shows checkbox stats
The tasks sub-tab SHALL display task completion statistics prominently.

#### Scenario: Tasks with checkboxes
- **WHEN** tasks.md contains `- [x]` and `- [ ]` lines
- **THEN** the header SHALL show "Tasks: N/M complete" and the content SHALL render the full markdown

### Requirement: Specs sub-tab shows change specs
The specs sub-tab SHALL show specs within the change's `specs/` subdirectory.

#### Scenario: Change with specs
- **WHEN** the specs sub-tab is selected for a change that has `specs/` with subdirectories
- **THEN** the spec names SHALL be listed and the first spec's content SHALL be displayed

#### Scenario: Change without specs
- **WHEN** the specs sub-tab is selected for a change with no `specs/` directory
- **THEN** the content SHALL display "No specs in this change"

### Requirement: Change list navigation
The user SHALL be able to navigate the change list with j/k keys when the
changes tab is active and no search prompt is open.

#### Scenario: Navigate changes with j/k
- **WHEN** the changes tab is active and the user presses j or k
- **THEN** the change cursor SHALL move to the next or previous row of the table

#### Scenario: j and k while the search prompt is open
- **WHEN** the search prompt is open and the user presses j or k
- **THEN** those characters SHALL be entered into the query
