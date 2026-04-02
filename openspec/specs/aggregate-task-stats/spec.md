# aggregate-task-stats Specification

## Purpose
TBD - created by archiving change open-tasks-in-header. Update Purpose after archive.
## Requirements
### Requirement: Project-level task aggregation
`ProjectInfo` SHALL expose `TasksTotal` and `TasksDone` fields representing the sum of all task checkboxes across active changes.

#### Scenario: Project with tasks across multiple changes
- **WHEN** a project has 2 active changes, one with 3/5 tasks done and another with 2/4 tasks done
- **THEN** `ProjectInfo.TasksTotal` SHALL be 9 and `ProjectInfo.TasksDone` SHALL be 5

#### Scenario: Project with no tasks
- **WHEN** a project has no active changes, or no changes contain tasks.md
- **THEN** `ProjectInfo.TasksTotal` SHALL be 0 and `ProjectInfo.TasksDone` SHALL be 0

### Requirement: Tasks displayed in persistent header
The persistent header stats line SHALL include the aggregate task count when tasks exist.

#### Scenario: Project has open tasks
- **WHEN** `ProjectInfo.TasksTotal` is greater than 0
- **THEN** the stats line SHALL append `Tasks: <done>/<total>` after the existing stats

#### Scenario: Project has no tasks
- **WHEN** `ProjectInfo.TasksTotal` is 0
- **THEN** the stats line SHALL NOT include any task information

