## ADDED Requirements

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

## REMOVED Requirements

### Requirement: The specs tab has a focus
**Reason**: One of its scenarios describes cycling through a log panel that no
longer exists, and OpenSpec cannot drop a scenario from a modified requirement.
Replaced by "The specs tab has two halves and one focus", which says the same
thing about the two halves that remain.
