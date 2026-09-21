## ADDED Requirements

### Requirement: The paging keys do not wait for focus
The keys that move by a page, by half a page, or to an end SHALL act on the
document whenever one is displayed, whichever half of a split tab holds the
keyboard.

A list is never paged: three rows or twenty-four, it is walked with the line
keys. So these keys are unambiguously a gesture at the document, and asking the
user to move the keyboard first buys nothing.

The line keys are the exception, and stay with the keyboard. `j` and `k` have to
choose between moving a list and scrolling a document, and where the keyboard is
is what decides that.

#### Scenario: Paging while the list holds the keyboard
- **GIVEN** a split tab whose list holds the keyboard
- **WHEN** the user presses `pgdown`, `ctrl+f`, `pgup`, `ctrl+b`, `ctrl+d` or
  `ctrl+u`
- **THEN** the document beside the list SHALL scroll

#### Scenario: Jumping while the list holds the keyboard
- **GIVEN** a split tab whose list holds the keyboard
- **WHEN** the user presses `gg` or `G`
- **THEN** the document SHALL move to its first or last row

#### Scenario: The line keys stay with the keyboard
- **GIVEN** a split tab whose list holds the keyboard
- **WHEN** the user presses `j` or `k`
- **THEN** the list selection SHALL move and the document SHALL NOT scroll

#### Scenario: The log panel still takes precedence
- **GIVEN** the log panel holds the keyboard
- **WHEN** any of these keys is pressed
- **THEN** the log SHALL move and the document SHALL NOT
