# modal-presentation Specification

## Purpose
TBD - created by archiving change modal-takes-the-whole-frame. Update Purpose after archive.

## Requirements

### Requirement: A modal occupies the whole frame
When a modal is up, the frame SHALL contain the modal centred on blank space and
nothing else. The view underneath is not drawn. This covers the project picker,
the scan indicator, the error box, and every confirmation prompt.

Only one modal is ever up at a time, because the component holding the keyboard
refuses the keys that would raise another. The order the modals are drawn in
therefore decides which one wins when the model somehow holds two states at
once, not what is layered over what.

#### Scenario: A confirmation is raised from the change list
- **WHEN** the user presses a key that raises a confirmation prompt
- **THEN** the frame SHALL show the prompt centred on blank space, without the
  change list behind it

#### Scenario: The project picker is open
- **WHEN** the project picker is open
- **THEN** the frame SHALL show the picker centred on blank space, without the
  view it was opened from behind it

#### Scenario: A modal never pushes the frame out of shape
- **WHEN** any modal is up
- **THEN** the frame SHALL still be exactly as wide and as tall as the terminal

### Requirement: A modal may accept typing
A modal MAY hold an editable field. While it does, every printable key SHALL
belong to that field, and only the keys the modal names SHALL do anything else.

Every modal until now answered `y`, `n` or any key at all. A modal that accepts
typing has to take those keys away from the actions they otherwise trigger,
which is the same rule the inline search prompt already follows.

#### Scenario: Printable keys reach the field
- **WHEN** a modal with an editable field is up and a printable key is pressed
- **THEN** it SHALL be entered into the field, and SHALL NOT trigger an action
  bound to that key elsewhere

#### Scenario: The modal names its keys
- **WHEN** a modal with an editable field is up
- **THEN** it SHALL show which keys submit, cancel and complete, because they
  cannot be guessed from a field

#### Scenario: Still one modal at a time
- **WHEN** a modal with an editable field is up
- **THEN** no key SHALL raise another modal, as with every other modal

#### Scenario: Still the whole frame
- **WHEN** a modal with an editable field is up
- **THEN** the frame SHALL show it centred on blank space and nothing else, and
  SHALL be exactly as wide and as tall as the terminal
