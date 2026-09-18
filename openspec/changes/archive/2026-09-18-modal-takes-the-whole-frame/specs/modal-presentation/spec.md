## ADDED Requirements

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
