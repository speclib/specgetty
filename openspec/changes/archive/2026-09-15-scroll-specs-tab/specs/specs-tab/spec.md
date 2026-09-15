## ADDED Requirements

### Requirement: The specs tab has a focus
Either the spec list or the spec content SHALL hold the keyboard, and the user
SHALL be able to move between them.

#### Scenario: Default focus
- **WHEN** the user switches to the specs tab
- **THEN** the spec list SHALL hold the keyboard

#### Scenario: Moving focus to the content
- **WHEN** the spec list holds the keyboard and the user presses `tab`
- **THEN** the spec content SHALL hold the keyboard

#### Scenario: Moving focus back
- **WHEN** the spec content holds the keyboard, the log panel is closed, and the
  user presses `tab`
- **THEN** the spec list SHALL hold the keyboard

#### Scenario: Cycling through the log panel
- **WHEN** the log panel is open and the user presses `tab` repeatedly on the
  specs tab
- **THEN** the keyboard SHALL move through the spec list, the spec content and
  the log panel in turn

#### Scenario: Focus is visible
- **WHEN** the spec list holds the keyboard
- **THEN** the selected spec SHALL be highlighted, and when the content holds the
  keyboard instead that highlight SHALL be dimmed, so the half receiving the keys
  is never a guess

#### Scenario: Leaving and returning to the tab
- **WHEN** the user switches away from the specs tab and back
- **THEN** the spec list SHALL hold the keyboard again

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

## MODIFIED Requirements

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
