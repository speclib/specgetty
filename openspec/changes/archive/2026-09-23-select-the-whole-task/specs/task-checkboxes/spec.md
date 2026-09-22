## MODIFIED Requirements

### Requirement: The tasks pane has a cursor
When the tasks artifact of a change is shown, one task SHALL be selected, and
the selection SHALL be visible in full.

A task is not a source line. A task item is a `- [ ] ` or `- [x] ` line at
column zero together with the indented lines that follow it, and it ends at the
first line that is unindented, blank, a heading, or the next checkbox. The two
prefixes are the ones the scanner counts, so an indented checkbox is a
continuation line rather than a task of its own, and the boxes on screen cannot
disagree with the totals beside them.

#### Scenario: A line that fits on one row
- **WHEN** the selected task occupies a single source line that fits on one
  screen row
- **THEN** that row SHALL be highlighted

#### Scenario: A line that wraps
- **WHEN** the selected task's source line wraps onto several screen rows
- **THEN** every one of those rows SHALL be highlighted

#### Scenario: A task with continuation lines
- **WHEN** the selected task is a checkbox line followed by indented
  continuation lines
- **THEN** every screen row of the checkbox line and of its continuation lines
  SHALL be highlighted, as one unbroken band

#### Scenario: The highlight stops at the task
- **WHEN** a task is selected
- **THEN** no row belonging to another task, to a heading or to a blank line
  SHALL be highlighted

#### Scenario: Moving the cursor
- **WHEN** the user presses `j`, `k`, up or down in the tasks pane
- **THEN** the cursor SHALL move to the neighbouring task

#### Scenario: Headings and blank lines are passed over
- **WHEN** a heading or a blank line lies between two tasks and the user moves
  the cursor across it
- **THEN** the cursor SHALL move from one task to the other without stopping in
  between

#### Scenario: The view follows the cursor
- **WHEN** the cursor moves to a task that is not fully visible
- **THEN** the pane SHALL scroll so that the whole of that task is shown, and a
  task taller than the pane SHALL be shown from its first row

#### Scenario: Bounds
- **WHEN** the cursor is on the first or last task and the user moves further in
  that direction
- **THEN** the cursor SHALL NOT move

#### Scenario: A tasks file with no tasks
- **WHEN** the tasks artifact contains no checkbox at column zero
- **THEN** the pane SHALL have no cursor, and its keys SHALL scroll it by rows
  as any document without a cursor is scrolled

#### Scenario: Other documents are unaffected
- **WHEN** any document other than a change's tasks is shown
- **THEN** it SHALL have no cursor and its keys SHALL behave as they did before

### Requirement: Space toggles the selected task
Pressing `space` SHALL change the state of the checkbox on the selected task and
persist it immediately. The file written SHALL be the one under the root the
project's content was resolved to.

#### Scenario: Ticking a task
- **WHEN** the selected task is unchecked and the user presses `space`
- **THEN** that task SHALL become checked on disk, and the pane SHALL show it
  checked

#### Scenario: Unticking a task
- **WHEN** the selected task is checked and the user presses `space`
- **THEN** that task SHALL become unchecked on disk

#### Scenario: The checkbox line is the one written
- **WHEN** a task with continuation lines is toggled
- **THEN** the checkbox line SHALL be the line rewritten, and the continuation
  lines SHALL be left exactly as they are

#### Scenario: A line that is not a task
- **WHEN** `space` is pressed on a pane whose cursor is not on a task, which
  now means a pane with no cursor at all, there being no task in it to select
- **THEN** nothing SHALL be written and nothing SHALL change

#### Scenario: The counts follow
- **WHEN** a task's state changes
- **THEN** the task counts shown elsewhere SHALL come to agree with the file,
  without the user rescanning

#### Scenario: The reader is not interrupted
- **WHEN** a task is toggled
- **THEN** no modal SHALL be raised by the rescan that follows, and no key
  pressed while it runs SHALL be discarded

#### Scenario: A task in a store-backed project
- **GIVEN** a change open in a project that declares a store
- **WHEN** a task is toggled
- **THEN** the `tasks.md` under the store SHALL be written, rather than a path
  under the repository the user started in that holds no such file

## ADDED Requirements

### Requirement: The tasks pane advertises its own keys
The nav bar SHALL describe the tasks pane as it behaves, and SHALL do so exactly
where the pane has a cursor, by the rule it already follows for `E` and for the
card view arrows: a key is listed where it does something.

#### Scenario: The toggle is advertised
- **WHEN** a pane with a task cursor holds the keyboard
- **THEN** the nav bar SHALL list a hint for `space`

#### Scenario: The movement keys are described correctly
- **WHEN** a pane with a task cursor holds the keyboard
- **THEN** the nav bar SHALL describe `j`, `k` and the arrow keys as navigating
  rather than as scrolling

#### Scenario: A document without a cursor
- **WHEN** an artifact other than tasks holds the keyboard
- **THEN** the nav bar SHALL NOT list `space`, and SHALL describe `j`, `k` and
  the arrow keys as scrolling

#### Scenario: A narrow terminal
- **WHEN** the nav bar has too little width for every hint
- **THEN** the hint for `space` SHALL be kept in preference to the hints for the
  paging keys, the toggle being the action of the pane
