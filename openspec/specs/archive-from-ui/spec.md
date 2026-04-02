# archive-from-ui Specification

## Purpose
TBD - created by archiving change archive-change-from-ui. Update Purpose after archive.
## Requirements
### Requirement: Archive keybinding
The system SHALL provide an `a` keybinding that initiates archiving when the changes tab is active and a change is selected.

#### Scenario: Press a on changes tab with a change selected
- **WHEN** the user presses `a` while the changes tab is focused and a change is highlighted
- **THEN** the system SHALL show a confirmation modal for archiving that change

#### Scenario: Press a when no changes exist
- **WHEN** the user presses `a` on the changes tab with no active changes
- **THEN** the system SHALL do nothing

#### Scenario: Press a on a different tab
- **WHEN** the user presses `a` while a tab other than changes is active
- **THEN** the system SHALL do nothing

### Requirement: CLI availability check
The system SHALL verify that the `openspec` CLI is available on PATH before attempting to archive.

#### Scenario: openspec not installed
- **WHEN** the user initiates an archive and `openspec` is not found on PATH
- **THEN** the system SHALL display a modal with the message that the openspec CLI is not installed

### Requirement: Incomplete task warning
The system SHALL warn the user when archiving a change that has incomplete tasks.

#### Scenario: Change has incomplete tasks
- **WHEN** the user initiates an archive for a change where `TasksDone < TasksTotal`
- **THEN** the confirmation modal SHALL display the count of incomplete tasks and ask "Archive anyway? (y/n)"

#### Scenario: Change has all tasks complete
- **WHEN** the user initiates an archive for a change where `TasksDone == TasksTotal` or `TasksTotal == 0`
- **THEN** the confirmation modal SHALL display "Archive <name>? (y/n)" without a task warning

### Requirement: Confirmation modal
The system SHALL require explicit confirmation before executing the archive command.

#### Scenario: User confirms with y
- **WHEN** the user presses `y` on the confirmation modal
- **THEN** the system SHALL execute the archive command

#### Scenario: User cancels with n or escape
- **WHEN** the user presses `n` or `escape` on the confirmation modal
- **THEN** the system SHALL dismiss the modal and return to normal state

### Requirement: Archive execution
The system SHALL execute `openspec archive <name> -y` from the project directory and capture the output.

#### Scenario: Successful archive
- **WHEN** the archive command exits with code 0
- **THEN** the system SHALL display a success modal with the command output and rescan the project

#### Scenario: Failed archive
- **WHEN** the archive command exits with a non-zero code
- **THEN** the system SHALL display a failure modal with the error output

### Requirement: Result feedback
The system SHALL display the archive result in a modal that dismisses on any keypress.

#### Scenario: Dismiss result modal
- **WHEN** the user presses any key while the result modal is shown
- **THEN** the modal SHALL be dismissed and the UI SHALL return to normal state

### Requirement: Nav bar hint
The nav bar SHALL show an `a archive` hint when the changes tab is active and changes exist.

#### Scenario: Changes tab active with changes
- **WHEN** the changes tab is active and the project has active changes
- **THEN** the nav bar SHALL include the `a archive` keybinding hint

#### Scenario: Changes tab not active
- **WHEN** a tab other than changes is active
- **THEN** the nav bar SHALL NOT include the `a archive` hint

