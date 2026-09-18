## ADDED Requirements

### Requirement: Tab content is drawn in its own border
The content of the active tab SHALL be drawn inside a border of its own, placed
directly below the tab bar with no blank row between them. The border SHALL
carry no title.

This border is in addition to the panel's, so the content is inset from the
panel border and again from its own.

#### Scenario: The change list
- **WHEN** the changes tab is active
- **THEN** the table and its header row SHALL be drawn inside the border, and
  the project header and the tab bar SHALL be outside it

#### Scenario: An open change
- **WHEN** a change is open on one of its artifacts
- **THEN** the artifact SHALL be drawn inside the border, and the change header
  line and the artifact sub-tab row SHALL be outside it

#### Scenario: The border is not a signal
- **WHEN** any tab is active
- **THEN** the panel's own border SHALL keep the colour and the meaning it has
  today, and the content border SHALL NOT change what it says about where the
  keyboard is

### Requirement: A label that names the content sits outside the border
A line that says what the content is, rather than being the content, SHALL be
drawn above the border alongside the tab bar. A line that reports on the content
SHALL be drawn inside it.

#### Scenario: The config tab's source file
- **WHEN** the config tab shows a file
- **THEN** the dimmed line naming that file SHALL sit above the border, with the
  tab bar

#### Scenario: The change list's search prompt
- **WHEN** a filter is being typed or is applied to the change list
- **THEN** the prompt and its count of matching rows SHALL sit inside the
  border, under the table it describes

### Requirement: The border costs rows, not correctness
The content region SHALL shrink by two rows and four columns to make room for
the border and its inset, and every view SHALL still fit the terminal exactly.

#### Scenario: At the smallest supported terminal
- **WHEN** the terminal is 60 by 20, the smallest specgetty draws at
- **THEN** the frame SHALL still be exactly 20 lines and 60 columns, with the
  border drawn and fewer content rows inside it
