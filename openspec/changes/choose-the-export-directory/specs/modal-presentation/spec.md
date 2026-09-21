## ADDED Requirements

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
