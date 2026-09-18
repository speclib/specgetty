## ADDED Requirements

### Requirement: Table columns are separated by two blank columns
Adjacent columns of the change list SHALL be separated by two blank columns
rather than one, so that a value filling its column is still clearly apart from
the value beside it.

#### Scenario: A value that fills its column
- **WHEN** a change name is as wide as the name column allows
- **THEN** two blank columns SHALL separate it from the task progress beside it

#### Scenario: A value that is cut
- **WHEN** a change name is wider than the name column allows
- **THEN** it SHALL be cut with an ellipsis, and two blank columns SHALL
  separate the ellipsis from the value beside it

#### Scenario: Columns dropped at a narrow width
- **WHEN** the panel is too narrow for the configured columns at the wider
  separation
- **THEN** columns SHALL be dropped from the right by the existing rule, rather
  than the separation being given up to keep them
