## MODIFIED Requirements

### Requirement: The properties tab is a list beside its content
The tab SHALL be drawn as a narrow list of sections and a wide content pane,
each in its own border, following the split the specs tab already uses. The
keyboard SHALL move between the two halves, and the lit border SHALL say which
half holds it.

#### Scenario: Moving between the halves
- **WHEN** the properties tab is active and the user presses the key that moves
  focus
- **THEN** the keyboard SHALL move between the list and the content, as it does
  on the specs tab

#### Scenario: Moving down the list
- **WHEN** the list holds the keyboard
- **THEN** the vertical keys SHALL move the selected section

#### Scenario: Scrolling the content
- **WHEN** the content holds the keyboard
- **THEN** the vertical keys SHALL scroll the document

#### Scenario: Which border is lit
- **WHEN** either half holds the keyboard
- **THEN** that half's border SHALL be lit and the other's dim

#### Scenario: The list is sized to its labels
- **WHEN** the tab is drawn
- **THEN** the list SHALL take the width its labels need rather than a fixed
  share of the panel, so that the content keeps the room its paths require
