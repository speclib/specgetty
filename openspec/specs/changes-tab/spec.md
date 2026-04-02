# changes-tab Specification

## Purpose
TBD - created by archiving change changes-tab-with-artifacts. Update Purpose after archive.
## Requirements
### Requirement: Changes tab lists active changes with task progress
The changes tab SHALL display a list of active changes on the left side with task completion stats.

#### Scenario: Project with active changes
- **WHEN** the user switches to the changes tab for a project with changes
- **THEN** the left side SHALL list all change directory names sorted alphabetically, each showing task progress as `(done/total)`

#### Scenario: Change with no tasks.md
- **WHEN** a change directory has no tasks.md file
- **THEN** the change SHALL be listed without task stats

#### Scenario: No active changes
- **WHEN** the project has no directories under `openspec/changes/`
- **THEN** the changes tab SHALL display "No active changes"

### Requirement: Changes tab shows artifact sub-tabs
When a change is selected, the right side SHALL show sub-tabs for each discovered artifact file plus a specs sub-tab.

#### Scenario: Change with standard artifacts
- **WHEN** a change is selected that contains proposal.md, design.md, and tasks.md
- **THEN** sub-tabs SHALL show: `proposal`, `design`, `tasks`, `specs`

#### Scenario: Change with custom artifacts
- **WHEN** a change contains additional .md files (e.g. notes.md, research.md)
- **THEN** those files SHALL appear as additional sub-tabs

#### Scenario: Sub-tab navigation
- **WHEN** the changes tab is active and the user presses left/right arrows
- **THEN** the artifact sub-tab SHALL change and the content SHALL update

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
The user SHALL be able to navigate the change list with j/k keys when the changes tab is active.

#### Scenario: Navigate changes with j/k
- **WHEN** the detail panel is focused, the changes tab is active, and the user presses j or k
- **THEN** the change cursor SHALL move and the right side SHALL update to show the newly selected change's artifacts

