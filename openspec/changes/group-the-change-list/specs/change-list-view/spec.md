## ADDED Requirements

### Requirement: Active and archived changes are grouped in one list
The change list SHALL show active and archived changes together, always, as one
table with one column header row, divided by group headers naming each group and
counting the rows in it. Active changes SHALL come first, because they are what
the developer is working on.

The group a row sits in is what says whether it is active or archived. No filter
selects between them: the list is not in a mode, and there is no state to be in
without knowing it.

#### Scenario: Both groups shown
- **WHEN** the change list is displayed for a project with active and archived
  changes
- **THEN** it SHALL show one table with a header for each group, the active
  group first, and every change under its own group

#### Scenario: A group with nothing in it
- **WHEN** a group contains no changes
- **THEN** its header SHALL still be shown with a count of zero, so that having
  nothing in flight is stated rather than left as an absence

#### Scenario: A project with no changes at all
- **WHEN** a project has neither active nor archived changes
- **THEN** both group headers SHALL be shown with counts of zero

#### Scenario: Group headers under a filter
- **WHEN** a search filter is narrowing the list
- **THEN** each group header SHALL count the rows the filter left in that group,
  and SHALL still be shown when that count is zero

#### Scenario: The cursor never rests on a header
- **WHEN** the user moves the cursor through the list
- **THEN** it SHALL move between changes only, passing over group headers
  without stopping on them

#### Scenario: Scrolling past a header
- **WHEN** the list is longer than the panel and the cursor is moved beyond it
- **THEN** the view SHALL scroll so the selected change stays visible, counting
  group headers as the rows they occupy

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

## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: Active and archived changes share one list
**Reason**: Replaced by "Active and archived changes are grouped in one list".
Every one of its scenarios described a filter mode, and the modes are gone: what
they selected between is now said by which group a row sits in. The mode was
also the least visible state in the application, indicated only by a label in
the nav bar, so being in one without realising it looked exactly like a project
whose active work had disappeared.
