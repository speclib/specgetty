## ADDED Requirements

### Requirement: The picker can be paged
The project picker SHALL move its cursor by a page and by half a page, beside
the jumps to either end it already has.

#### Scenario: Paging the picker
- **WHEN** the picker holds the keyboard and the user presses `pgdown`, `pgup`,
  `ctrl+f`, `ctrl+b`, `ctrl+d` or `ctrl+u`
- **THEN** the selected project SHALL move by a page or half a page

#### Scenario: The jumps are unchanged
- **WHEN** the user presses `g` or `G`
- **THEN** the selection SHALL move to the first or last project, as it already
  does

#### Scenario: Within the filter
- **WHEN** a filter is narrowing the picker
- **THEN** these keys SHALL move within what the filter left
