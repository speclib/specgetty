## ADDED Requirements

### Requirement: A border is lit when the keyboard is inside it
Every border in the frame SHALL be drawn in the active colour when the region
holding the keyboard lies inside it, and in the dim colour when it does not.
Borders nest, so more than one may be lit at once, and the innermost lit border
is the region actually receiving the keys.

#### Scenario: A single-pane tab with the log panel closed
- **WHEN** the changes or config tab is active and the log panel is closed
- **THEN** the panel border and the content border SHALL both be lit, because
  the keyboard is inside both

#### Scenario: The log panel takes the keyboard
- **WHEN** the log panel is open and holds the keyboard
- **THEN** the log panel's border SHALL be lit, and the panel border and the
  content border SHALL both be dim

#### Scenario: Two panes side by side
- **WHEN** the content is split into panes and one of them holds the keyboard
- **THEN** that pane's border SHALL be lit, every other pane's border SHALL be
  dim, and the panel border enclosing them SHALL be lit
