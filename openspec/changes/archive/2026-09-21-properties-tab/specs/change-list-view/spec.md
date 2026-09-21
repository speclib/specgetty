## MODIFIED Requirements

### Requirement: A change opens at its own navigation level
A change SHALL open at a navigation level below the change list. `enter`
descends from the list into the change, and `esc` returns to the list. An open
change SHALL name the workflow schema it records, because which artifacts the
change needs follows from it.

This capability describes the change list and the change below it. What lies
above the change list is defined by the project navigation and is deliberately
not constrained here.

#### Scenario: Descend from the change list
- **WHEN** the change list is displayed with a change under the cursor and the
  user presses `enter`
- **THEN** that change SHALL open at its own level, filling the panel, with its
  artifact sub-tabs

#### Scenario: Ascend from a change
- **WHEN** a change is open and the user presses `esc`
- **THEN** the change list SHALL be displayed with the same change still under
  the cursor

#### Scenario: An open change names its schema
- **WHEN** a change is open and its `.openspec.yaml` records a schema
- **THEN** that schema's name SHALL be visible above the artifact document,
  beside the change's own name

#### Scenario: A change with no recorded schema
- **WHEN** a change is open and it has no `.openspec.yaml`
- **THEN** no schema SHALL be named, rather than the project default being shown
  as though the change had recorded it
