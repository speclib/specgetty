# panel-layout Specification

## Purpose
TBD - created by archiving change pad-every-view. Update Purpose after archive.

## Requirements

### Requirement: Panel content is inset from its border
Every view drawn inside the main panel SHALL begin one column right of the left
border and end one column left of the right border. The rule belongs to the
panel, not to the views, so a view added later inherits it without doing
anything.

This applies to both tab bars, the change table and its search prompt, the specs
list and its document, the config pane, the markdown documents, the change
header line, the placeholder for an unimplemented tab, and the log panel.

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

### Requirement: Nothing compensates for the inset on its own
No view SHALL add horizontal padding of its own on top of the panel's inset.
A view that indents itself as well would sit two columns in while its neighbours
sit one, which is the misalignment this change exists to remove.

#### Scenario: The project header
- **WHEN** the project header is drawn above the tab bar
- **THEN** its first character SHALL sit in the same column as the first
  character of the tab bar, of the table header, and of every row below it

#### Scenario: The empty state
- **WHEN** no project is selected and the panel shows the prompt to open the
  picker
- **THEN** that text SHALL begin in the same column as any other panel content

### Requirement: The nav bar is inset without breaking its fill
The nav bar SHALL indent its text by one column on each side while its
background continues to span the full terminal width.

#### Scenario: The nav bar at any width
- **WHEN** the nav bar is drawn
- **THEN** its first key hint SHALL begin one column in, the version at the
  right SHALL end one column in, and the background SHALL still be painted from
  the first terminal column to the last

#### Scenario: A status message replaces the hints
- **WHEN** a status message takes the nav bar's row
- **THEN** it SHALL be inset and filled by the same rule as the key hints
