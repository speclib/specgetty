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
- **WHEN** a document viewer is displayed and the user presses down, `j`, up or
  `k`
- **THEN** the document SHALL scroll by one row in that direction

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
