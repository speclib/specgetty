## ADDED Requirements

### Requirement: Picker columns are separated by two blank columns
The project picker draws its rows with the same table renderer as the change
list, and its columns SHALL be separated by the same two blank columns.

#### Scenario: A project name that fills its column
- **WHEN** a project name is as wide as its column allows
- **THEN** two blank columns SHALL separate it from the spec count beside it
