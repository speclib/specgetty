## MODIFIED Requirements

### Requirement: A document taller than its pane scrolls
When a document occupies more rows than its pane, the user SHALL be able to
reach every row with the keyboard.

A document with a cursor moves its cursor and lets the view follow; a document
without one moves the view directly. Every movement key follows that one split,
rather than `j` and `k` following it and the paging keys ignoring it.

#### Scenario: Move one row
- **WHEN** a document viewer without a cursor is displayed and the user presses
  down, `j`, up or `k`
- **THEN** the document SHALL scroll by one row in that direction

#### Scenario: Move one line, in a document with a cursor
- **WHEN** a document viewer with a cursor is displayed and the user presses
  down, `j`, up or `k`
- **THEN** the cursor SHALL move by one item and the view SHALL follow it, which
  may move the view by more than one row because an item can occupy several

#### Scenario: Move one page
- **WHEN** a document viewer without a cursor is displayed and the user presses
  `pgdown`, `ctrl+f`, `pgup` or `ctrl+b`
- **THEN** the document SHALL scroll by the height of the pane in that direction

#### Scenario: Move one page, in a document with a cursor
- **WHEN** a document viewer with a cursor is displayed and the user presses
  `pgdown`, `ctrl+f`, `pgup` or `ctrl+b`
- **THEN** the cursor SHALL move by as many items as fill one pane of rows, and
  the view SHALL follow it, so that the cursor is never left off screen

#### Scenario: Move half a page
- **WHEN** a document viewer is displayed and the user presses `ctrl+d` or
  `ctrl+u`
- **THEN** it SHALL move half as far as a page key moves it, by whichever of the
  two rules above applies to it

#### Scenario: Jump to the ends
- **WHEN** a document viewer without a cursor is displayed and the user presses
  `gg` or `G`
- **THEN** the view SHALL move to the first or the last row of the document
  respectively

#### Scenario: Jump to the ends, in a document with a cursor
- **WHEN** a document viewer with a cursor is displayed and the user presses
  `gg` or `G`
- **THEN** the cursor SHALL move to the first or the last item respectively, and
  the view SHALL follow it

#### Scenario: Bounds
- **WHEN** the first row of the document is at the top of the pane and the user
  scrolls up, or the last row is at the bottom and the user scrolls down
- **THEN** the view SHALL NOT move

#### Scenario: Paging past an end, in a document with a cursor
- **WHEN** a page key would move the cursor beyond the first or the last item
- **THEN** it SHALL stop on that item rather than wrapping or going out of range

#### Scenario: Document shorter than the pane
- **WHEN** a document occupies fewer rows than its pane and the user presses any
  scroll key
- **THEN** the view SHALL NOT move

### Requirement: The paging keys act on whichever list or document holds the keyboard
`pgdown`, `pgup`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` SHALL act
on whichever list or document holds the keyboard, by the same rule `j` and `k`
already follow. Where that surface has a cursor, they SHALL move the cursor;
where it has none, they SHALL move the view.

One rule rather than one per surface. A list that cannot be paged is walked a
row at a time however long it is, and a document that pages while a list holds
the keyboard moves something the user is not looking at. A document whose page
key moved its rows but not its cursor was a third case of the same mistake: the
cursor was left where the reader was not, and the keys that act on it then acted
out of sight.

#### Scenario: A document that holds the keyboard
- **WHEN** a document holds the keyboard and one of these keys is pressed
- **THEN** it SHALL move by that amount: its cursor where it has one, and its
  rows where it has none

#### Scenario: A document with a cursor that holds the keyboard
- **WHEN** a document with a cursor holds the keyboard and one of these keys is
  pressed
- **THEN** the cursor SHALL move, and the view SHALL follow it

#### Scenario: A list that holds the keyboard
- **WHEN** a list holds the keyboard and one of these keys is pressed
- **THEN** the list's cursor SHALL move, and no document SHALL scroll

#### Scenario: An overlay that holds the keyboard
- **WHEN** the project picker is open and holds the keyboard
- **THEN** its cursor SHALL move, and nothing behind it SHALL
