# task-checkboxes Specification

## Purpose
Reading and changing a change's task list from inside specgetty: how checkbox
lines look, how one is selected, and what happens to the file when it is ticked.

## Requirements

### Requirement: Checkbox lines render as boxes
A markdown checkbox SHALL be drawn as a box rather than as its source
punctuation.

#### Scenario: An unchecked task
- **WHEN** a line begins a markdown task that is not done
- **THEN** it SHALL be drawn with an empty box in place of the `- [ ]` prefix

#### Scenario: A completed task
- **WHEN** a line begins a markdown task that is done
- **THEN** it SHALL be drawn with a filled box in place of the `- [x]` prefix

#### Scenario: The two states are distinguishable without colour
- **WHEN** the two boxes are compared
- **THEN** they SHALL differ in shape, so the state is readable on a highlighted
  row and by a reader who cannot rely on colour

#### Scenario: Both boxes occupy one cell
- **WHEN** either box is drawn
- **THEN** it SHALL occupy exactly one terminal cell, so that the width the
  renderer counts is the width the terminal draws

### Requirement: The tasks pane has a cursor
When the tasks artifact of a change is shown, one source line SHALL be selected,
and the selection SHALL be visible.

#### Scenario: A line that fits on one row
- **WHEN** the selected source line occupies a single screen row
- **THEN** that row SHALL be highlighted

#### Scenario: A line that wraps
- **WHEN** the selected source line wraps onto several screen rows
- **THEN** every one of those rows SHALL be highlighted, and no row belonging to
  another source line SHALL be

#### Scenario: Moving the cursor
- **WHEN** the user presses `j`, `k`, up or down in the tasks pane
- **THEN** the cursor SHALL move to the neighbouring source line

#### Scenario: The view follows the cursor
- **WHEN** the cursor moves to a line that is not fully visible
- **THEN** the pane SHALL scroll so that the whole of that line is shown

#### Scenario: Bounds
- **WHEN** the cursor is on the first or last source line and the user moves
  further in that direction
- **THEN** the cursor SHALL NOT move

#### Scenario: Other documents are unaffected
- **WHEN** any document other than a change's tasks is shown
- **THEN** it SHALL have no cursor and its keys SHALL behave as they did before

### Requirement: Space toggles the selected task
Pressing `space` SHALL change the state of the checkbox on the selected line and
persist it immediately.

#### Scenario: Ticking a task
- **WHEN** the selected line is an unchecked task and the user presses `space`
- **THEN** that task SHALL become checked on disk, and the pane SHALL show it
  checked

#### Scenario: Unticking a task
- **WHEN** the selected line is a checked task and the user presses `space`
- **THEN** that task SHALL become unchecked on disk

#### Scenario: A line that is not a task
- **WHEN** the selected line carries no checkbox and the user presses `space`
- **THEN** nothing SHALL be written and nothing SHALL change

#### Scenario: The counts follow
- **WHEN** a task's state changes
- **THEN** the task counts shown elsewhere SHALL come to agree with the file,
  without the user rescanning

### Requirement: A toggle never discards another writer's work
The saved file SHALL be built from the contents of `tasks.md` as they are on
disk at the moment of the toggle, not from any copy read earlier.

#### Scenario: The file changed since it was read
- **WHEN** `tasks.md` has been edited and saved elsewhere since specgetty read
  it, and the user toggles a task
- **THEN** those edits SHALL survive, and only the toggled checkbox SHALL differ
  from what is on disk

#### Scenario: The selected line is no longer present
- **WHEN** the selected line cannot be found in the file as it now stands
- **THEN** nothing SHALL be written, and the user SHALL be told that the file
  changed

#### Scenario: The selected line appears more than once
- **WHEN** the selected line appears more than once in the file, so the intended
  one cannot be identified
- **THEN** nothing SHALL be written, and the user SHALL be told

### Requirement: The save is atomic
A reader of `tasks.md` SHALL never observe a partially written file.

#### Scenario: Replacing the file
- **WHEN** a toggle is saved
- **THEN** the file SHALL be replaced in a single step, so that any reader sees
  either the previous contents or the new contents

#### Scenario: Permissions are preserved
- **WHEN** a toggle is saved
- **THEN** the file's permissions SHALL be what they were before the save

#### Scenario: The save fails
- **WHEN** the file cannot be written, for instance because it is read only
- **THEN** the user SHALL be told, and `tasks.md` SHALL be left as it was
