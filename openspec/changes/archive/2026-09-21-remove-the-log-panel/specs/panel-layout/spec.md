## ADDED Requirements

### Requirement: A border is lit when the keyboard is in the region it encloses
Every border in the frame SHALL be drawn in the active colour when the region
holding the keyboard lies inside it, and in the dim colour when it does not.
Borders nest, so more than one may be lit at once, and the innermost lit border
is the region actually receiving the keys.

#### Scenario: A single-pane tab
- **WHEN** the changes tab is active
- **THEN** the panel border and the content border SHALL both be lit, because
  the keyboard is inside both

#### Scenario: Two panes side by side
- **WHEN** the content is split into panes and one of them holds the keyboard
- **THEN** that pane's border SHALL be lit, every other pane's border SHALL be
  dim, and the panel border enclosing them SHALL be lit

#### Scenario: An open change
- **WHEN** a change is open
- **THEN** the panel border and the artifact's border SHALL both be lit

#### Scenario: A modal
- **WHEN** a modal is awaiting an answer
- **THEN** it replaces the frame rather than being drawn over it, so no border
  underneath it is on screen to be lit or dimmed

## MODIFIED Requirements

### Requirement: Panel content is inset from its border
Every view drawn inside the main panel SHALL begin one column right of the left
border and end one column left of the right border. The rule belongs to the
panel, not to the views, so a view added later inherits it without doing
anything.

This applies to both tab bars, the change table and its search prompt, the specs
list and its document, the properties list and its content, the markdown
documents, the change header line, and the placeholder for an unimplemented tab.

#### Scenario: A list row under the cursor
- **WHEN** a row of the change list is under the cursor
- **THEN** its highlight SHALL begin and end one column inside the border, in
  the same way the project picker already draws its selected row

#### Scenario: A wrapped document line
- **WHEN** a line of a markdown document is long enough to wrap
- **THEN** it SHALL wrap one column earlier than the border, leaving a blank
  column between the text and the frame on both sides

#### Scenario: A view added later
- **WHEN** a new tab or pane is drawn inside the main panel
- **THEN** it SHALL be inset without that view rendering any padding of its own

## REMOVED Requirements

### Requirement: A border is lit when the keyboard is inside it
**Reason**: Two of its scenarios describe a log panel that no longer exists, and
OpenSpec cannot drop a scenario from a modified requirement. Replaced by "A
border is lit when the keyboard is in the region it encloses", which is the same
rule with the log cases removed and the remaining ones unchanged.
