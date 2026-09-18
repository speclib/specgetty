## MODIFIED Requirements

### Requirement: The specs tab has a focus
Either the spec list or the spec content SHALL hold the keyboard, and the user
SHALL be able to move between them. Each SHALL be drawn in its own border, and
the one holding the keyboard SHALL be the lit one.

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
- **THEN** the spec list's border SHALL be lit and the spec content's border
  SHALL be dim, the selected spec SHALL be highlighted, and when the content
  holds the keyboard instead the two borders SHALL swap and that highlight
  SHALL be dimmed

#### Scenario: Neither half has the keyboard
- **WHEN** the log panel holds the keyboard while the specs tab is active
- **THEN** both halves' borders SHALL be dim

#### Scenario: Leaving and returning to the tab
- **WHEN** the user switches away from the specs tab and back
- **THEN** the spec list SHALL hold the keyboard again
