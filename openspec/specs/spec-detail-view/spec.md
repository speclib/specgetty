# spec-detail-view Specification

## Purpose
Describes reading one specification as the structure it already has: an outline
of its Purpose, its requirements and their scenarios, beside a card that shows
whichever of them the reader is on, and what happens when a file does not fit
that structure.

## Requirements

### Requirement: A spec opens in a view of its own
Pressing `enter` on a spec in the specs tab SHALL open that spec in a view one
level below the project view, and `esc` SHALL return to the specs tab with the
same spec selected. This is the rule the change list already follows, applied to
a second list.

#### Scenario: Opening a spec
- **GIVEN** the specs tab is active and its list holds the keyboard
- **WHEN** the user presses `enter`
- **THEN** the selected spec SHALL open in the spec detail view

#### Scenario: Leaving the view
- **WHEN** the user presses `esc` in the spec detail view
- **THEN** the specs tab SHALL be shown again with the same spec selected

#### Scenario: Enter goes no deeper
- **WHEN** the user presses `enter` in the spec detail view
- **THEN** nothing SHALL happen, because the view is the floor of that branch

#### Scenario: The view names the spec
- **WHEN** the spec detail view is open
- **THEN** the name of the spec SHALL be on screen, so the view is never
  anonymous

### Requirement: The view is an outline beside a card
The spec detail view SHALL draw two halves, an outline of the spec on the left
and a card on the right, each in its own border, and `tab` SHALL move the
keyboard between them. The half holding the keyboard SHALL be the lit one. This
is the focus model the specs tab and the properties tab already use.

#### Scenario: Default focus
- **WHEN** the spec detail view opens
- **THEN** the outline SHALL hold the keyboard

#### Scenario: Moving the focus
- **WHEN** the outline holds the keyboard and the user presses `tab`
- **THEN** the card SHALL hold the keyboard, and its border SHALL be the lit one

#### Scenario: Moving it back
- **WHEN** the card holds the keyboard and the user presses `tab`
- **THEN** the outline SHALL hold it again, there being only two places for it
  to be

#### Scenario: Leaving and returning
- **WHEN** the user leaves the view and opens the same spec again
- **THEN** the outline SHALL hold the keyboard again

### Requirement: The outline lists what the spec contains
The outline SHALL list the spec's Purpose, then each requirement in the order the
file gives them, with that requirement's scenarios under it in file order. A
requirement SHALL be distinguishable from a scenario by how it is drawn and not
by its indentation alone, because a wrapped label occupies the indentation the
level below it would use.

#### Scenario: The order is the file's order
- **WHEN** a spec with several requirements is opened
- **THEN** the outline SHALL list Purpose first and then the requirements in the
  order they appear in the file, each followed by its own scenarios in file order

#### Scenario: The levels are told apart
- **WHEN** the outline is drawn
- **THEN** a requirement SHALL be drawn differently from a scenario, and that
  difference SHALL survive a label wrapping onto a second row

### Requirement: The outline wraps a long label rather than cutting it
A label wider than the outline half SHALL be wrapped onto as many rows as it
needs. A node that occupies several rows SHALL be selected as one node: the
cursor SHALL cover every row it produced, and moving the cursor onto it SHALL
bring all of its rows into view.

#### Scenario: A label wider than the half
- **WHEN** a requirement title is wider than the outline half
- **THEN** it SHALL be wrapped onto as many rows as it needs, and no part of it
  SHALL be cut off

#### Scenario: The cursor covers the whole node
- **WHEN** the cursor is on a node that wrapped onto several rows
- **THEN** every one of those rows SHALL be highlighted

#### Scenario: A wrapped node at the edge of the pane
- **WHEN** the cursor moves onto a node whose last row is below the bottom of the
  pane
- **THEN** the outline SHALL scroll far enough that the whole node is visible

#### Scenario: One step is one node
- **WHEN** the user presses `j` or `k` in the outline
- **THEN** the cursor SHALL move one node, whatever number of rows that node and
  its neighbour occupy

### Requirement: Every node has a card
The card SHALL show the node the cursor is on, and what it shows SHALL depend on
the kind of node.

#### Scenario: Purpose
- **WHEN** the cursor is on Purpose
- **THEN** the card SHALL show the spec's purpose prose

#### Scenario: A requirement
- **WHEN** the cursor is on a requirement
- **THEN** the card SHALL show the text between that requirement's heading and
  its first scenario, and SHALL NOT show its scenarios

#### Scenario: A scenario
- **WHEN** the cursor is on a scenario
- **THEN** the card SHALL show that scenario's title and clauses

### Requirement: A scenario's card lays its clauses out
A scenario's clause is one source bullet, `- **WHEN** the user...`. The card
SHALL put the keyword on a row of its own with the clause indented under it,
rather than rendering the bullet as it stands. The title SHALL be bold, the
keywords SHALL be coloured, a backticked span SHALL be highlighted, and the prose
SHALL be left legible.

#### Scenario: A clause is laid out
- **WHEN** a scenario clause is drawn
- **THEN** its keyword SHALL be on a row of its own and the clause text SHALL be
  indented beneath it

#### Scenario: A clause longer than the card
- **WHEN** a clause is longer than the card can show on one row
- **THEN** it SHALL be wrapped and every row after the first SHALL keep the
  clause indentation

#### Scenario: The keywords
- **WHEN** a clause opens with `GIVEN`, `WHEN`, `THEN` or `AND`
- **THEN** that keyword SHALL be drawn in its own colour

#### Scenario: Code spans
- **WHEN** a clause contains a span between backticks
- **THEN** that span SHALL be highlighted

### Requirement: The card fits its pane at any width
The card SHALL be legible at every terminal width the application supports. Its
decorative padding SHALL shrink as the pane narrows, so that the space the text
itself needs is never given away to decoration.

#### Scenario: A narrow terminal
- **GIVEN** two terminal widths
- **WHEN** the same scenario is shown in each
- **THEN** the narrower one SHALL give the clause text proportionally more of the
  pane than the wider one does, and no row SHALL run past the border

### Requirement: The card is a document viewer
The card SHALL behave as a document viewer, with the wrapping, scrolling,
position reporting and position retention that capability describes. Moving the
outline cursor to another node SHALL start that node's card at its first row.

#### Scenario: A node taller than the card
- **WHEN** a node's card occupies more rows than the pane and the card holds the
  keyboard
- **THEN** the user SHALL be able to reach its last row, and the panel title
  SHALL report the reading position

#### Scenario: Moving to another node
- **WHEN** a card has been scrolled and the user selects another node
- **THEN** the new node's card SHALL be shown from its first row

#### Scenario: Position reported only when the card has the keyboard
- **WHEN** the outline holds the keyboard
- **THEN** the panel title SHALL show no scroll position

### Requirement: The spec is parsed when it is opened
Nothing SHALL be parsed until a spec is opened, and the result SHALL be held for
as long as that spec is open rather than derived again for each frame.

#### Scenario: A session that never opens a spec
- **WHEN** the application runs without the spec detail view being opened
- **THEN** no spec SHALL have been parsed

#### Scenario: Redrawing costs nothing
- **WHEN** the view is redrawn without the open spec having changed on disk
- **THEN** the spec SHALL NOT be parsed again

### Requirement: The reading position survives a rescan
The spec on disk can be rewritten while it is being read. The outline cursor
SHALL be kept by which node it is on rather than by its position in the list, so
that an edit above the cursor does not move the reader.

#### Scenario: An edit above the cursor
- **GIVEN** the cursor is on a scenario and the spec is rewritten with a new
  requirement added above it
- **WHEN** the view is redrawn from the new content
- **THEN** the cursor SHALL still be on the same scenario

#### Scenario: The node is gone
- **WHEN** the node the cursor was on no longer exists in the rewritten spec
- **THEN** the cursor SHALL land on a node that does exist rather than out of
  range

### Requirement: A spec that does not fit the structure does not open
A spec file that cannot be read as a Purpose, requirements and scenarios SHALL
NOT open in the detail view. The specs tab SHALL keep its cursor where it is and
report why, and the spec SHALL remain readable in full as markdown on the tab
itself.

#### Scenario: A file without requirements
- **WHEN** the user presses `enter` on a spec whose file has no requirements
- **THEN** the detail view SHALL NOT open, and the reason SHALL be reported

#### Scenario: Nothing is hidden
- **WHEN** a spec cannot be opened in the detail view
- **THEN** its whole file SHALL still be shown as markdown on the specs tab, as
  it is today
