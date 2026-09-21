# change-list-view Specification

## Purpose
Gives the change list a full-width view with selectable columns and an
open/archived/both filter, and puts a single change on its own navigation level
below it so the change's artifact sub-tabs and the project tab bar never share
keys.

## Requirements

### Requirement: A change opens at its own navigation level
A change SHALL open at a navigation level below the change list. `enter`
descends from the list into the change, and `esc` returns to the list. An open
change SHALL name the workflow schema it records, because which artifacts the
change needs follows from it.

This capability describes the change list and the change below it. What lies
above the change list is defined by the project navigation and is deliberately
not constrained here.

#### Scenario: Descend from the change list
- **WHEN** the change list is displayed with a change under the cursor and the
  user presses `enter`
- **THEN** that change SHALL open at its own level, filling the panel, with its
  artifact sub-tabs

#### Scenario: Ascend from a change
- **WHEN** a change is open and the user presses `esc`
- **THEN** the change list SHALL be displayed with the same change still under
  the cursor

#### Scenario: An open change names its schema
- **WHEN** a change is open and its `.openspec.yaml` records a schema
- **THEN** that schema's name SHALL be visible above the artifact document,
  beside the change's own name

#### Scenario: A change with no recorded schema
- **WHEN** a change is open and it has no `.openspec.yaml`
- **THEN** no schema SHALL be named, rather than the project default being shown
  as though the change had recorded it

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

### Requirement: Change list occupies the full panel width
The change list SHALL be rendered as a table across the full width of the
detail panel, not beside an artifact viewer.

#### Scenario: Change list rendering
- **WHEN** the changes tab is active
- **THEN** the change list SHALL span the full panel width with a column header
  row, and no artifact content SHALL be shown beside it

#### Scenario: Cursor row
- **WHEN** the change list is displayed
- **THEN** the row under the cursor SHALL be visually highlighted

### Requirement: Change list columns are selectable
The set of columns shown SHALL be selectable through the config file and a
command-line option, with the option taking priority.

#### Scenario: Default columns
- **WHEN** no column selection is configured
- **THEN** the list SHALL show the change name, task progress as `done/total`,
  and the number of specs in the change

#### Scenario: Columns from config
- **WHEN** the config file names a list of fields
- **THEN** those fields SHALL be shown as columns in the order given

#### Scenario: Columns from the command line
- **WHEN** the application is started with `--change-fields=` naming a list of
  fields
- **AND** the config file also names a list of fields
- **THEN** the command-line list SHALL be used and the config list SHALL be
  ignored

#### Scenario: Unknown field name
- **WHEN** a configured or supplied field name is not recognised
- **THEN** the application SHALL report the unknown name and the valid names
  rather than failing silently

#### Scenario: Change with no tasks.md
- **WHEN** a change directory has no tasks.md file
- **THEN** the task progress column SHALL be empty for that row

### Requirement: Active and archived changes share one list
Active changes and archived changes SHALL be shown in a single list with a
filter selecting which are included. The three states SHALL be named `active`,
`archived` and `active+archived` wherever the user meets them, whether on screen
or in configuration.

#### Scenario: Filter modes
- **WHEN** the change list is displayed
- **THEN** the filter SHALL be in one of three modes: active only, archived
  only, or both, with the current mode visible

#### Scenario: Default mode
- **WHEN** the change list is first displayed for a project
- **THEN** the filter SHALL be in the configured default mode, which is active
  only when nothing is configured

#### Scenario: Default mode from configuration
- **WHEN** a default mode is set in the configuration file
- **THEN** the change list SHALL open in that mode

#### Scenario: Default mode from the command line
- **WHEN** a default mode is given on the command line
- **AND** one is also set in the configuration file
- **THEN** the command line value SHALL be used

#### Scenario: An unrecognised mode
- **WHEN** a configured or supplied mode is not one of the three names
- **THEN** the application SHALL report the unknown value together with the
  valid ones, rather than falling back silently

#### Scenario: Returning to the default
- **WHEN** the user has changed the mode and then selects a different project
- **THEN** the filter SHALL return to the configured default

#### Scenario: Both mode
- **WHEN** the filter is in both mode
- **THEN** active and archived changes SHALL appear in the same list, and each
  row SHALL indicate whether it is active or archived

#### Scenario: No changes in the selected mode
- **WHEN** the filter mode selects no changes
- **THEN** the list SHALL state which mode is active and that it contains no
  changes

### Requirement: Actions operate on the row under the cursor
The archive, discard and export actions SHALL act on the change currently under
the cursor in the change list, whatever filter is applied.

#### Scenario: Action with a filter applied
- **WHEN** a filter is narrowing the change list and the user triggers an action
- **THEN** the action SHALL apply to the change under the cursor as displayed

#### Scenario: Archive action on an archived change
- **WHEN** the change under the cursor is archived and the user triggers archive
- **THEN** the system SHALL do nothing

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

### Requirement: Table columns are separated by two blank columns
Adjacent columns of the change list SHALL be separated by two blank columns
rather than one, so that a value filling its column is still clearly apart from
the value beside it.

#### Scenario: A value that fills its column
- **WHEN** a change name is as wide as the name column allows
- **THEN** two blank columns SHALL separate it from the task progress beside it

#### Scenario: A value that is cut
- **WHEN** a change name is wider than the name column allows
- **THEN** it SHALL be cut with an ellipsis, and two blank columns SHALL
  separate the ellipsis from the value beside it

#### Scenario: Columns dropped at a narrow width
- **WHEN** the panel is too narrow for the configured columns at the wider
  separation
- **THEN** columns SHALL be dropped from the right by the existing rule, rather
  than the separation being given up to keep them
