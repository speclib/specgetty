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
command-line option, with the option taking priority. The default set SHALL
include the archive date, which is blank for an active change.

#### Scenario: Default columns
- **WHEN** no column selection is configured
- **THEN** the list SHALL show the change name, task progress as `done/total`,
  the number of specs in the change, and the archive date

#### Scenario: The date column on an active change
- **WHEN** an active change is listed
- **THEN** its archive date SHALL be blank, because it has none

#### Scenario: The state column is no longer a default
- **WHEN** no column selection is configured
- **THEN** the column naming a change active or archived SHALL NOT be shown,
  because the group header says it, and it SHALL remain available to configure

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

### Requirement: Actions operate on the row under the cursor
The archive, discard and export actions SHALL act on the change currently under
the cursor in the change list, whichever group it sits in and whatever filter is
applied.

#### Scenario: Action with a filter applied
- **WHEN** a filter is narrowing the change list and the user triggers an action
- **THEN** the action SHALL apply to the change under the cursor as displayed

#### Scenario: Archive action on an archived change
- **WHEN** the change under the cursor is archived and the user triggers archive
- **THEN** the system SHALL do nothing

#### Scenario: Archiving moves a change between groups
- **WHEN** an active change is archived
- **THEN** it SHALL appear in the archived group on the next read, at the top of
  it, and SHALL no longer appear in the active group

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

### Requirement: Each group has a fixed order
Active changes SHALL be ordered by name. Archived changes SHALL be ordered by
archive date, most recent first, so that the work most recently finished is at
the top of its group. Neither order is selectable.

#### Scenario: Active order
- **WHEN** the active group is displayed
- **THEN** its changes SHALL be in alphabetical order by name

#### Scenario: Archived order
- **WHEN** the archived group is displayed
- **THEN** its changes SHALL be in order of archive date, most recent first

#### Scenario: Archived changes sharing a date
- **WHEN** two archived changes carry the same date
- **THEN** their order relative to each other SHALL be stable from one render to
  the next

#### Scenario: An archived change with no date in its directory name
- **WHEN** an archived change's directory carries no parsable date prefix
- **THEN** it SHALL still be listed, ordered after those that do, rather than
  being dropped or sorted unpredictably

### Requirement: A configuration setting that no longer applies is reported
A `change_mode` key in the configuration file selected among filter modes that
no longer exist. YAML ignores keys a program does not know, so such a key would
otherwise sit in a configuration doing nothing. The application SHALL report it.

#### Scenario: A configuration still naming a change mode
- **WHEN** the configuration file carries a `change_mode` key
- **THEN** the application SHALL report that the key no longer has an effect,
  and SHALL start normally

#### Scenario: A configuration without it
- **WHEN** the configuration file carries no `change_mode` key
- **THEN** nothing SHALL be reported

### Requirement: The change list can be paged and jumped through
The change list SHALL move its cursor by a page, by half a page, and to either
end, with the same keys every other surface uses.

A project's archive grows without bound, so a list that only moves a row at a
time gets slower to use for as long as the project lives.

#### Scenario: Paging down the list
- **WHEN** the change list holds the keyboard and the user presses `pgdown` or
  `ctrl+f`
- **THEN** the selected change SHALL move down by about the number of rows on
  screen

#### Scenario: Paging up
- **WHEN** the user presses `pgup` or `ctrl+b`
- **THEN** the selection SHALL move up by the same amount

#### Scenario: Half a page
- **WHEN** the user presses `ctrl+d` or `ctrl+u`
- **THEN** the selection SHALL move by half that

#### Scenario: To the ends
- **WHEN** the user presses `gg` or `G`
- **THEN** the selection SHALL move to the first or the last change

#### Scenario: Paging across a group boundary
- **WHEN** a page move would cross from the active group into the archived one
- **THEN** it SHALL land on a change rather than on a group header

#### Scenario: The selection is remembered
- **WHEN** the selection is moved by any of these keys
- **THEN** it SHALL be remembered the way a single-row move already is, so a
  rescan or a filter keeps it

#### Scenario: A filtered list
- **WHEN** a search is narrowing the list
- **THEN** these keys SHALL move within what the filter left
