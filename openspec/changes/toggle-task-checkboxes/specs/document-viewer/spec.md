## MODIFIED Requirements

### Requirement: A document taller than its pane scrolls
When a document occupies more rows than its pane, the user SHALL be able to reach
every row with the keyboard.

#### Scenario: Move one row
- **WHEN** a document viewer without a cursor is displayed and the user presses
  down, `j`, up or `k`
- **THEN** the document SHALL scroll by one row in that direction

#### Scenario: Move one line, in a document with a cursor
- **WHEN** a document viewer with a cursor is displayed and the user presses
  down, `j`, up or `k`
- **THEN** the cursor SHALL move by one source line and the view SHALL follow
  it, which may move the view by more than one row because a source line can
  wrap

#### Scenario: Move one page
- **WHEN** a document viewer is displayed and the user presses `pgdown`,
  `ctrl+f`, `pgup` or `ctrl+b`
- **THEN** the document SHALL scroll by the height of the pane in that direction

#### Scenario: Move half a page
- **WHEN** a document viewer is displayed and the user presses `ctrl+d` or
  `ctrl+u`
- **THEN** the document SHALL scroll by half the height of the pane in that
  direction

#### Scenario: Jump to the ends
- **WHEN** a document viewer is displayed and the user presses `gg` or `G`
- **THEN** the view SHALL move to the first or the last row of the document
  respectively

#### Scenario: Bounds
- **WHEN** the first row of the document is at the top of the pane and the user
  scrolls up, or the last row is at the bottom and the user scrolls down
- **THEN** the view SHALL NOT move

#### Scenario: Document shorter than the pane
- **WHEN** a document occupies fewer rows than its pane and the user presses any
  scroll key
- **THEN** the view SHALL NOT move
