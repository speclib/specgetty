## ADDED Requirements

### Requirement: The match hint is separated from the table
The column naming which files a contents search matched SHALL be separated from
the table by the same gap that separates every other column.

It is currently appended flush, and looked right only because the last default
column held values shorter than its width. The archive date fills its column
exactly, so the hint lands against it.

#### Scenario: A row whose last column is full
- **WHEN** a row's last column holds a value as wide as the column allows and a
  match hint is shown beside it
- **THEN** the gap between them SHALL be the same as between any two columns

#### Scenario: The column header
- **WHEN** the header row is drawn with a match hint column
- **THEN** its gap SHALL be the same, rather than depending on the last header's
  length

#### Scenario: The picker's table
- **WHEN** a match hint is shown in the project picker
- **THEN** the same gap SHALL separate it, because both tables share the
  measurement
