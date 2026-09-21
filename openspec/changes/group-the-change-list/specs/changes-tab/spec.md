## MODIFIED Requirements

### Requirement: Changes tab lists active changes with task progress
The changes tab SHALL display the change list across the full panel width as a
table, with the columns selected by configuration, and with active and archived
changes in their own groups.

#### Scenario: Project with active changes
- **WHEN** the user switches to the changes tab for a project with changes
- **THEN** the full panel width SHALL show a table of changes with a column
  header row, grouped, the active group first

#### Scenario: Change with no tasks.md
- **WHEN** a change directory has no tasks.md file
- **THEN** the change SHALL be listed with an empty task progress column

#### Scenario: No active changes
- **WHEN** the project has no directories under `openspec/changes/` other than
  the archive
- **THEN** the active group SHALL be shown with a count of zero, and the
  archived group SHALL be shown beneath it
