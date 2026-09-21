## ADDED Requirements

### Requirement: The properties tab is a grouped list beside its content
The tab SHALL be drawn as a narrow list and a wide content pane, each in its own
border. The keyboard SHALL move between the two halves, and the lit border SHALL
say which half holds it.

The list SHALL take a width that gives its labels room without taking room the
content needs, rather than only the width its labels happen to need. It SHALL
NOT take a fixed share of the panel: its labels are short, and a share would hand
a wide terminal columns that nothing fills.

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

#### Scenario: The list keeps its width as the terminal grows
- **GIVEN** two terminals of different widths, both wide enough for the tab
- **WHEN** the tab is drawn in each
- **THEN** the list SHALL be the same width in both, and every extra column
  SHALL go to the content

#### Scenario: The list gives way on a narrow terminal
- **WHEN** the terminal is narrow enough that the list's width would crowd the
  content
- **THEN** the list SHALL be reduced rather than the content, and neither half
  SHALL be drawn wider than the panel holding them

### Requirement: The properties list is grouped into sections
The list SHALL be drawn in groups, each under a header naming what the rows
below it are. A header SHALL NOT be selectable: the cursor SHALL move from one
row to the next as though headers were not there, and no key SHALL leave the
cursor on one.

The groups SHALL be the project itself, holding the configuration and the store,
and the schemas, holding one row per schema the project uses. What a row reports
is unchanged by which group it sits in.

#### Scenario: What the groups hold
- **WHEN** the properties tab is drawn for a project
- **THEN** the configuration row and the store row SHALL be under the project
  header, and every schema row under the schemas header

#### Scenario: The cursor steps over a header
- **GIVEN** the cursor is on the last row of a group
- **WHEN** the user presses the key that moves down
- **THEN** the cursor SHALL land on the first row of the next group, and never
  on the header between them

#### Scenario: Jumping to either end
- **WHEN** `gg` or `G` is pressed while the list holds the keyboard
- **THEN** the cursor SHALL land on the first or the last row, and never on a
  header

#### Scenario: A group that holds nothing
- **GIVEN** a project whose changes record no schema and whose configuration
  names no default
- **WHEN** the properties tab is drawn
- **THEN** the schemas header SHALL still be drawn, because a group that says it
  is empty teaches the concept where an absent one teaches nothing

#### Scenario: The configuration row says what it is
- **WHEN** the list is drawn
- **THEN** the row showing the project's configuration SHALL be labelled for the
  configuration it shows

#### Scenario: A list too long for the pane
- **GIVEN** a project with more schema rows than the pane has lines
- **WHEN** the cursor is moved to a row below the fold
- **THEN** the list SHALL scroll far enough to show that row, counting the lines
  its headers occupy, and SHALL NOT open on a blank line

## REMOVED Requirements

### Requirement: The properties tab is a list beside its content
**Reason**: Replaced by "The properties tab is a grouped list beside its
content", which keeps every guarantee it made except one. Its scenario "The list
is sized to its labels" required the list to take only the width its labels need,
"so that the content keeps the room its paths require". Measured, that reason
does not hold: the longest line the content shows is a store root of 60 columns,
which wraps at an 80-column terminal and fits at a 100-column one under both the
old policy and the new one. The policy never decides whether a path wraps, so the
columns the section headers need are not taken from anything.

**Migration**: None for a user. The four scenarios about the two halves, the
keyboard and the lit border are carried unchanged into the requirement that
replaces it.
