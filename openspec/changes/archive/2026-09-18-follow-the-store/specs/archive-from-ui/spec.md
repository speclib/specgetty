## MODIFIED Requirements

### Requirement: Archive execution
The system SHALL execute `openspec archive <name> -y` from the resolved root and
capture the output. Where the open project was resolved through a store
declaration, the resolved root is the store, not the repo the user started in.

#### Scenario: Successful archive
- **WHEN** the archive command exits with code 0
- **THEN** the system SHALL display a success modal with the command output and rescan the project

#### Scenario: Failed archive
- **WHEN** the archive command exits with a non-zero code
- **THEN** the system SHALL display a failure modal with the error output

#### Scenario: Archiving in a store-backed project
- **WHEN** a change is archived in a project resolved through a store
  declaration
- **THEN** the command SHALL run with the store as its working directory
