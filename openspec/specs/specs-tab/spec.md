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
