## MODIFIED Requirements

### Requirement: What is copied
The name and the path SHALL each be the form that is useful where it is pasted.
A path SHALL be built from the root the project's content was resolved to, not
from the directory the user started in.

#### Scenario: The name is the one commands take
- **WHEN** the name of a change is copied
- **THEN** it SHALL be the name the change list displays, which is the name
  `openspec` commands accept

#### Scenario: The path of an active change
- **WHEN** the path of an active change is copied
- **THEN** it SHALL be the absolute path of
  `<resolved root>/openspec/changes/<directory>`

#### Scenario: The path of an archived change
- **WHEN** the path of an archived change is copied
- **THEN** it SHALL be the absolute path of the directory as it exists on disk,
  including the `YYYY-MM-DD-` prefix that the displayed name does not carry

#### Scenario: The path in a store-backed project
- **GIVEN** a project that declares a store
- **WHEN** the path of one of its changes is copied
- **THEN** it SHALL be under the store, which is where the change is, and not
  under the repository the user started in

#### Scenario: The copied path exists
- **WHEN** any change's path is copied
- **THEN** that path SHALL resolve to a directory on disk
