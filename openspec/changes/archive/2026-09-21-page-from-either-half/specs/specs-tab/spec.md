## ADDED Requirements

### Requirement: The specs tab pages from either half
The spec content SHALL respond to the paging and jump keys whether the spec list
or the spec content holds the keyboard. This is the same rule the properties tab
follows, and it is what `document-viewer` has required all along.

#### Scenario: Paging while the spec list holds the keyboard
- **GIVEN** the specs tab is active and its list holds the keyboard
- **WHEN** the user presses `pgdown` or `G`
- **THEN** the selected spec's content SHALL scroll

#### Scenario: The reading position is unchanged
- **GIVEN** the spec list holds the keyboard and the content is scrolled
- **THEN** the panel title SHALL still report no position, which is what this
  capability already requires and what nobody asked to change
