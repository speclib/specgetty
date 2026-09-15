## ADDED Requirements

### Requirement: The artifact pane of an open change is a document viewer
The content shown under an artifact sub-tab of an open change SHALL be a
document viewer, with the wrapping, scrolling, position reporting and position
retention that capability describes.

#### Scenario: Long artifact
- **WHEN** a change is open on an artifact longer than the panel
- **THEN** the user SHALL be able to reach the end of that artifact with the
  keyboard, and the panel title SHALL report the reading position

#### Scenario: Vertical keys at the change level
- **WHEN** a change is open and the user presses down, `j`, up, `k`, a page key
  or `gg` or `G`
- **THEN** the artifact content SHALL scroll, and no cursor belonging to a list
  above this level SHALL move

#### Scenario: Specs sub-tab
- **WHEN** a change is open on its specs sub-tab, showing several spec deltas in
  one pane
- **THEN** that pane SHALL scroll as a single document

## MODIFIED Requirements

### Requirement: Artifact sub-navigation is scoped to the open change
Left and right arrows while a change is open SHALL move only between that
change's artifact sub-tabs, and SHALL NOT change the project tab bar.

#### Scenario: Right arrow at the last artifact sub-tab
- **WHEN** a change is open with the last artifact sub-tab active and the user
  presses right
- **THEN** the active sub-tab SHALL NOT change and the project tab bar SHALL NOT
  change

#### Scenario: Left arrow at the first artifact sub-tab
- **WHEN** a change is open with the first artifact sub-tab active and the user
  presses left
- **THEN** the active sub-tab SHALL NOT change and the project tab bar SHALL NOT
  change

#### Scenario: Tab bar keys belong to the change list
- **WHEN** the change list is displayed and the user presses left or right
- **THEN** the active project tab SHALL change

#### Scenario: A new sub-tab starts at the top
- **WHEN** a change is open, its artifact has been scrolled, and the user moves
  to a different artifact sub-tab
- **THEN** the new artifact SHALL be shown from its first row
