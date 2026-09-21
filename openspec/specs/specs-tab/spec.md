# specs-tab Specification

## Purpose
TBD - created by archiving change specs-tab-with-listing. Update Purpose after archive.

## Requirements

### Requirement: Specs tab displays a list of specs
The specs tab SHALL display a vertical list of spec names on the left side of the detail panel.

#### Scenario: Project with specs
- **WHEN** the user switches to the specs tab for a project that has specs
- **THEN** the left side SHALL list all spec directory names sorted alphabetically

#### Scenario: Project with no specs
- **WHEN** the user switches to the specs tab for a project with no specs
- **THEN** the tab SHALL display "No specs found"

### Requirement: Specs tab displays selected spec content
The specs tab SHALL display the content of the selected spec's `spec.md` file on the right side, rendered as styled markdown.

#### Scenario: Spec selected
- **WHEN** a spec is highlighted in the list
- **THEN** the right side SHALL show the spec.md content rendered with markdown styling (headers, lists, bold, italic)

#### Scenario: Spec without spec.md
- **WHEN** a spec directory exists but has no spec.md file
- **THEN** the right side SHALL display "No spec.md found"

### Requirement: Spec list navigation
The user SHALL be able to navigate the spec list with j/k keys when the specs tab
is active and the spec list holds the keyboard.

#### Scenario: Navigate specs with j/k
- **WHEN** the specs tab is active, the spec list holds the keyboard, and the
  user presses j or k
- **THEN** the spec cursor SHALL move down or up and the content panel SHALL
  update to show the newly selected spec

#### Scenario: Cursor bounds
- **WHEN** the spec cursor is at the first or last spec
- **THEN** pressing k or j respectively SHALL not move the cursor beyond the
  bounds

#### Scenario: j and k while the content holds the keyboard
- **WHEN** the spec content holds the keyboard and the user presses j or k
- **THEN** the content SHALL scroll and the spec cursor SHALL NOT move

### Requirement: The spec content is a document viewer
When the spec content holds the keyboard, it SHALL behave as a document viewer,
with the wrapping, scrolling, position reporting and position retention that
capability describes.

#### Scenario: Spec longer than the pane
- **WHEN** the content of a spec occupies more rows than the pane and the content
  holds the keyboard
- **THEN** the user SHALL be able to reach its last row with the keyboard, and
  the panel title SHALL report the reading position

#### Scenario: Position reported only when focused
- **WHEN** the spec list holds the keyboard
- **THEN** the panel title SHALL show no scroll position

#### Scenario: Selecting a different spec
- **WHEN** a spec has been scrolled and the user selects a different spec
- **THEN** the newly selected spec SHALL be shown from its first row

#### Scenario: Content wrapped to its own half
- **WHEN** a spec contains a line wider than the content half of the tab
- **THEN** that line SHALL be wrapped to the width of that half, not to the width
  of the whole panel

### Requirement: The spec list can be paged and jumped through
The spec list SHALL move its cursor by a page, by half a page, and to either
end, while it holds the keyboard.

#### Scenario: Paging the spec list
- **WHEN** the spec list holds the keyboard and a page or half-page key is
  pressed
- **THEN** the selected spec SHALL move, and the spec content SHALL NOT scroll

#### Scenario: Jumping the spec list
- **WHEN** the spec list holds the keyboard and `gg` or `G` is pressed
- **THEN** the selection SHALL move to the first or last spec

#### Scenario: The content still pages when it holds the keyboard
- **WHEN** the spec content holds the keyboard and a page key is pressed
- **THEN** the content SHALL scroll, as it always has

### Requirement: The specs tab has two halves and one focus
Either the spec list or the spec content SHALL hold the keyboard, and the user
SHALL be able to move it between them. Of the two borders drawn, the one holding
the keyboard SHALL be the lit one.

#### Scenario: Default focus
- **WHEN** the specs tab becomes active
- **THEN** the spec list SHALL hold the keyboard

#### Scenario: Moving the focus
- **WHEN** the user presses `tab` on the specs tab
- **THEN** the keyboard SHALL move to the other half

#### Scenario: Moving it back
- **WHEN** the user presses `tab` again
- **THEN** the keyboard SHALL return to the half it came from, there being only
  two places for it to be

#### Scenario: The lit border
- **WHEN** either half holds the keyboard
- **THEN** that half's border SHALL be lit and the other's dim

### Requirement: Enter opens the selected spec in detail
The specs tab SHALL treat `enter` the way the change list does: it descends into
the selected spec. The detail view it opens is described by `spec-detail-view`.

#### Scenario: Enter on a spec
- **GIVEN** the specs tab is active and its list holds the keyboard
- **WHEN** the user presses `enter`
- **THEN** the selected spec SHALL open in the spec detail view

#### Scenario: Enter while the content holds the keyboard
- **GIVEN** the specs tab is active and its content holds the keyboard
- **WHEN** the user presses `enter`
- **THEN** the selected spec SHALL open in the spec detail view, because which
  half has the keyboard does not change which spec is selected

#### Scenario: Returning from the detail view
- **WHEN** the user leaves the spec detail view with `esc`
- **THEN** the specs tab SHALL be shown again with the same spec selected and the
  spec list holding the keyboard

#### Scenario: A project with no specs
- **WHEN** the user presses `enter` on the specs tab of a project with no specs
- **THEN** nothing SHALL happen

### Requirement: Backticked spans are highlighted
A span between backticks in a rendered markdown document SHALL be drawn
distinctly from the prose around it, on the specs tab and in every other document
the same renderer draws.

#### Scenario: A code span in a spec
- **WHEN** a spec contains a span between backticks
- **THEN** that span SHALL be highlighted, and the backticks themselves SHALL NOT
  be shown

#### Scenario: An unclosed backtick
- **WHEN** a line contains a single backtick with no closing one
- **THEN** the line SHALL be rendered unchanged rather than swallowing the rest
  of the document

### Requirement: Enter descends whether or not the file fits the grammar
`enter` SHALL descend into the selected spec in every case. Whether the file can
be structured is answered by the view it opens and reported there, not on the
nav bar, because a file usually has several faults and the nav bar holds one
line. The tab SHALL go on showing the whole file as markdown, so the reader who
descends into a report is one `esc` from the content.

#### Scenario: Enter on a spec that does not fit
- **WHEN** the user presses `enter` on a spec whose file cannot be read as
  requirements and scenarios
- **THEN** the spec view SHALL open and report why, rather than the cursor
  staying put and the nav bar carrying the reason

#### Scenario: The nav bar is not used for the reason
- **WHEN** a spec cannot be structured
- **THEN** no transient report SHALL be put on the nav bar, the reasons having a
  place of their own

#### Scenario: The markdown view is unaffected
- **WHEN** a spec cannot be opened as an outline
- **THEN** the specs tab SHALL still render its whole file as markdown, so no
  content becomes unreachable
