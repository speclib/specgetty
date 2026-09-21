## MODIFIED Requirements

### Requirement: Space toggles the selected task
Pressing `space` SHALL change the state of the checkbox on the selected line and
persist it immediately. The file written SHALL be the one under the root the
project's content was resolved to.

#### Scenario: Ticking a task
- **WHEN** the selected line is an unchecked task and the user presses `space`
- **THEN** that task SHALL become checked on disk, and the pane SHALL show it
  checked

#### Scenario: Unticking a task
- **WHEN** the selected line is a checked task and the user presses `space`
- **THEN** that task SHALL become unchecked on disk

#### Scenario: A line that is not a task
- **WHEN** the selected line carries no checkbox and the user presses `space`
- **THEN** nothing SHALL be written and nothing SHALL change

#### Scenario: The counts follow
- **WHEN** a task's state changes
- **THEN** the task counts shown elsewhere SHALL come to agree with the file,
  without the user rescanning

#### Scenario: A task in a store-backed project
- **GIVEN** a change open in a project that declares a store
- **WHEN** a task is toggled
- **THEN** the `tasks.md` under the store SHALL be written, rather than a path
  under the repository the user started in that holds no such file
