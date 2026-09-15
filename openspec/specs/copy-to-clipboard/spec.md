# copy-to-clipboard Specification

## Purpose
Puts a reference to the selected change on the system clipboard: either its
name, for pasting into a command, or its absolute path, for pasting into a shell
or an editor.

## Requirements

### Requirement: Copy keybindings
The change list SHALL offer two keys that copy a reference to the selected
change to the system clipboard.

#### Scenario: Copy the name
- **WHEN** a change is selected in the change list and the user presses `y`
- **THEN** the change's name SHALL be placed on the system clipboard

#### Scenario: Copy the path
- **WHEN** a change is selected in the change list and the user presses `Y`
- **THEN** the absolute path of the change's directory SHALL be placed on the
  system clipboard

#### Scenario: No change selected
- **WHEN** the user presses `y` or `Y` with no change under the cursor
- **THEN** nothing SHALL be copied and nothing SHALL be reported

#### Scenario: Another view is active
- **WHEN** the user presses `y` or `Y` while a tab other than changes is active,
  or while an overlay holds the keyboard
- **THEN** the change list SHALL NOT act on the key

### Requirement: What is copied
The name and the path SHALL each be the form that is useful where it is pasted.

#### Scenario: The name is the one commands take
- **WHEN** the name of a change is copied
- **THEN** it SHALL be the name the change list displays, which is the name
  `openspec` commands accept

#### Scenario: The path of an active change
- **WHEN** the path of an active change is copied
- **THEN** it SHALL be the absolute path of
  `<project>/openspec/changes/<directory>`

#### Scenario: The path of an archived change
- **WHEN** the path of an archived change is copied
- **THEN** it SHALL be the absolute path of the directory as it exists on disk,
  including the `YYYY-MM-DD-` prefix that the displayed name does not carry

#### Scenario: The copied path exists
- **WHEN** any change's path is copied
- **THEN** that path SHALL resolve to a directory on disk

### Requirement: The result is reported without interrupting
Copying SHALL report what happened in a way that requires no acknowledgement,
because the action is instant and reversible.

#### Scenario: A successful copy
- **WHEN** a copy succeeds
- **THEN** a status message SHALL name what was copied, and SHALL NOT require a
  keypress to dismiss

#### Scenario: The message clears itself
- **WHEN** a status message is shown and the user presses any key
- **THEN** the message SHALL be replaced by the ordinary display

#### Scenario: The clipboard is unavailable
- **WHEN** the system has no usable clipboard mechanism, or writing to it fails
- **THEN** a status message SHALL say so, so that a copy which did not happen is
  distinguishable from one that did

### Requirement: Nav bar hint
The change list SHALL advertise the copy keys.

#### Scenario: A change is selected
- **WHEN** the change list has a change under the cursor
- **THEN** the nav bar SHALL show hints for the copy keys

#### Scenario: The list is empty
- **WHEN** the change list has no entries
- **THEN** the nav bar SHALL NOT show the copy hints
