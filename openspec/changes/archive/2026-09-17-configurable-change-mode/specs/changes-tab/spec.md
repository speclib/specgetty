## MODIFIED Requirements

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
