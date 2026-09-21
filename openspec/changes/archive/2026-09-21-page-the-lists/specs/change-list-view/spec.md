## ADDED Requirements

### Requirement: The change list can be paged and jumped through
The change list SHALL move its cursor by a page, by half a page, and to either
end, with the same keys every other surface uses.

A project's archive grows without bound, so a list that only moves a row at a
time gets slower to use for as long as the project lives.

#### Scenario: Paging down the list
- **WHEN** the change list holds the keyboard and the user presses `pgdown` or
  `ctrl+f`
- **THEN** the selected change SHALL move down by about the number of rows on
  screen

#### Scenario: Paging up
- **WHEN** the user presses `pgup` or `ctrl+b`
- **THEN** the selection SHALL move up by the same amount

#### Scenario: Half a page
- **WHEN** the user presses `ctrl+d` or `ctrl+u`
- **THEN** the selection SHALL move by half that

#### Scenario: To the ends
- **WHEN** the user presses `gg` or `G`
- **THEN** the selection SHALL move to the first or the last change

#### Scenario: Paging across a group boundary
- **WHEN** a page move would cross from the active group into the archived one
- **THEN** it SHALL land on a change rather than on a group header

#### Scenario: The selection is remembered
- **WHEN** the selection is moved by any of these keys
- **THEN** it SHALL be remembered the way a single-row move already is, so a
  rescan or a filter keeps it

#### Scenario: A filtered list
- **WHEN** a search is narrowing the list
- **THEN** these keys SHALL move within what the filter left
