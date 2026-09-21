## ADDED Requirements

### Requirement: The properties rows can be paged and jumped through
The properties row list SHALL move its selection by a page, by half a page, and
to either end, while it holds the keyboard.

Its list is short, so the jumps matter more than the pages; both are there
because every list in the application answers to the same keys.

#### Scenario: Jumping the row list
- **WHEN** the row list holds the keyboard and `gg` or `G` is pressed
- **THEN** the selection SHALL move to the first or last row, and the content
  SHALL NOT scroll

#### Scenario: Paging a list shorter than a page
- **WHEN** a page key is pressed and the list is shorter than a page
- **THEN** the selection SHALL move to that end and stop

#### Scenario: The content still pages when it holds the keyboard
- **WHEN** the content holds the keyboard and a page key is pressed
- **THEN** it SHALL scroll

## REMOVED Requirements

### Requirement: The properties tab pages from either half
**Reason**: Added one commit ago and wrong, for the reason given in `specs-tab`:
with the row list pageable there are two things a key could move, and this
required moving the one the keyboard is not on. The original complaint it was
meant to answer was about lists, not about this tab's document.
