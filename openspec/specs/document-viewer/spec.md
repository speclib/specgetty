# document-viewer Specification

## Purpose
Describes a pane that shows one markdown document: how the document is fitted to
the pane width, how the reader moves through it, how the pane reports where the
reader is, and when that position is kept or reset.

## Requirements

### Requirement: A document is wrapped to the pane width
A document viewer SHALL wrap every line of the rendered document to the width of
its content area, so that the number of rows the viewer produces is the number of
rows the terminal displays.

#### Scenario: Line longer than the pane
- **WHEN** a document contains a line wider than the content area
- **THEN** the line SHALL be wrapped onto as many rows as it needs, and no part
  of it SHALL be cut off or run past the right edge

#### Scenario: Styling survives a wrap
- **WHEN** a wrapped line contains styled text such as a header or a bold span
- **THEN** the styling SHALL be preserved on every row the line occupies

#### Scenario: Wrapping does not displace following content
- **WHEN** a document contains a line that wraps onto several rows
- **THEN** the lines after it SHALL still be reachable, and SHALL NOT be dropped
  to make room

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

### Requirement: The pane reports the reading position
A document viewer SHALL show where the reader is in a document that does not
fit, and SHALL show nothing when it does fit.

#### Scenario: Position while scrolling
- **WHEN** a document taller than its pane is displayed
- **THEN** the panel title SHALL show the position as a percentage, and that
  percentage SHALL change as the document scrolls

#### Scenario: Position at the end
- **WHEN** the last row of the document is visible at the bottom of the pane
- **THEN** the panel title SHALL report 100%

#### Scenario: Document that fits
- **WHEN** a document occupies fewer rows than its pane
- **THEN** the panel title SHALL show no position indicator

### Requirement: Scroll position follows the document
The reading position SHALL belong to the document being read, not to the pane.

#### Scenario: A different document is opened
- **WHEN** the pane switches to a different document
- **THEN** the new document SHALL be shown from its first row

#### Scenario: The same document is redisplayed
- **WHEN** the pane returns to a document that was scrolled, without any other
  document having been opened in between
- **THEN** the position SHALL be the position it was left at

#### Scenario: The document changes on disk
- **WHEN** a rescan or a file system event replaces the content of the document
  currently displayed
- **THEN** the reading position SHALL be kept, clamped so that it does not point
  past the end of the new content

#### Scenario: The terminal is resized
- **WHEN** the terminal size changes while a document is displayed
- **THEN** the document SHALL be re-wrapped to the new width and the position
  SHALL be clamped to the end of the re-wrapped document

### Requirement: A page is what the surface can show
Moving by a page SHALL move by the number of rows the surface is currently
showing, and by half that for the half-page keys, so that the keys mean the same
thing at any terminal size. Where a surface draws items whose rows are not all
one row tall, a page SHALL move the cursor by as many items as fill that number
of rows, so the rule stays a rule about rows and the cursor stays on an item.

#### Scenario: A page in a taller pane
- **GIVEN** two terminal heights
- **WHEN** a page key is pressed in each
- **THEN** the taller one SHALL move further

#### Scenario: Jumping to the ends
- **WHEN** `gg` or `G` is pressed
- **THEN** the surface SHALL go to its first or last row

#### Scenario: Paging past an end
- **WHEN** a page key would move beyond the first or last row
- **THEN** it SHALL stop there rather than wrapping or going out of range

#### Scenario: A page over items of unequal height
- **GIVEN** a list in which some items occupy more rows than others
- **WHEN** a page key is pressed
- **THEN** the cursor SHALL move by as many items as occupy one pane of rows,
  which is fewer items where the items are taller

#### Scenario: Items that are all one row tall
- **GIVEN** a list in which every item occupies one row
- **WHEN** a page key is pressed
- **THEN** the cursor SHALL move by the number of rows the list is showing, which
  is what it does today

### Requirement: The paging keys act on whichever list or document holds the keyboard
`pgdown`, `pgup`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` SHALL act
on whichever list or document holds the keyboard, by the same rule `j` and `k`
already follow.

One rule rather than one per surface. A list that cannot be paged is walked a
row at a time however long it is, and a document that pages while a list holds
the keyboard moves something the user is not looking at.

#### Scenario: A document that holds the keyboard
- **WHEN** a document holds the keyboard and one of these keys is pressed
- **THEN** the document SHALL scroll or jump

#### Scenario: A list that holds the keyboard
- **WHEN** a list holds the keyboard and one of these keys is pressed
- **THEN** the list's cursor SHALL move, and no document SHALL scroll

#### Scenario: An overlay that holds the keyboard
- **WHEN** the project picker is open and holds the keyboard
- **THEN** its cursor SHALL move, and nothing behind it SHALL
