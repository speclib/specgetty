## ADDED Requirements

### Requirement: The arrow keys stay inside the spec detail view
`left` and `right` SHALL do nothing in the spec detail view. The view shows one
spec and has no sibling view to move between, and its two halves are reached
with `tab`.

Neither key SHALL change the project tab bar, and neither SHALL move the
keyboard between the outline and the card. This is the rule
`change-list-view` already states for an open change, applied to the level below
it: a level's keys belong to that level and do not spill into the one above.

#### Scenario: Right arrow in the spec detail view
- **WHEN** the spec detail view is open and the user presses `right`
- **THEN** nothing SHALL happen, and the tab that will be active on returning to
  the specs tab SHALL be unchanged

#### Scenario: Left arrow in the spec detail view
- **WHEN** the spec detail view is open and the user presses `left`
- **THEN** nothing SHALL happen, and the tab that will be active on returning to
  the specs tab SHALL be unchanged

#### Scenario: The keyboard stays where it was
- **GIVEN** the card holds the keyboard in the spec detail view
- **WHEN** the user presses `left` or `right`
- **THEN** the card SHALL still hold the keyboard, and the lit border SHALL not
  move

#### Scenario: Returning to the specs tab
- **GIVEN** the user has pressed `left` or `right` any number of times in the
  spec detail view
- **WHEN** the user presses `esc`
- **THEN** the specs tab SHALL be the active tab, as it is when the keys were
  never pressed
