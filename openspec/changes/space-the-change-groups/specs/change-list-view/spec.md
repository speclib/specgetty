## MODIFIED Requirements

### Requirement: Active and archived changes are grouped in one list
The change list SHALL show active and archived changes together, always, as one
table with one column header row, divided by group headers naming each group and
counting the rows in it. Active changes SHALL come first, because they are what
the developer is working on.

One blank line SHALL separate each group from the one above it, so that the
boundary reads as a break rather than as another row. The line is drawn
unconditionally, whether or not either group has rows in it, because a rule with
no exception is one less thing on screen to explain.

The group a row sits in is what says whether it is active or archived. No filter
selects between them: the list is not in a mode, and there is no state to be in
without knowing it.

#### Scenario: Both groups shown
- **WHEN** the change list is displayed for a project with active and archived
  changes
- **THEN** it SHALL show one table with a header for each group, the active
  group first, and every change under its own group

#### Scenario: The groups are separated
- **WHEN** both group headers are on screen
- **THEN** one blank line SHALL sit between the last line of the group above and
  the header of the group below

#### Scenario: No blank line above the first group
- **WHEN** the change list is displayed
- **THEN** the first group's header SHALL follow the column header row directly,
  with no blank line between them

#### Scenario: A group with nothing in it
- **WHEN** a group contains no changes
- **THEN** its header SHALL still be shown with a count of zero, so that having
  nothing in flight is stated rather than left as an absence, and the blank line
  above it SHALL still be drawn

#### Scenario: A project with no changes at all
- **WHEN** a project has neither active nor archived changes
- **THEN** both group headers SHALL be shown with counts of zero

#### Scenario: Group headers under a filter
- **WHEN** a search filter is narrowing the list
- **THEN** each group header SHALL count the rows the filter left in that group,
  and SHALL still be shown when that count is zero

#### Scenario: The cursor never rests on a header
- **WHEN** the user moves the cursor through the list
- **THEN** it SHALL move between changes only, passing over group headers and the
  blank lines between groups without stopping on them

#### Scenario: Scrolling past a header
- **WHEN** the list is longer than the panel and the cursor is moved beyond it
- **THEN** the view SHALL scroll so the selected change stays visible, counting
  group headers and the blank lines between groups as the rows they occupy

#### Scenario: The pane never opens on a blank line
- **WHEN** scrolling would put the blank line between two groups at the top of
  the pane
- **THEN** the pane SHALL start at the group header below it instead, and the
  selected change SHALL remain visible
