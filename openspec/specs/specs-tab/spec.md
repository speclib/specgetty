# specs-tab Specification

## Purpose
TBD - created by archiving change specs-tab-with-listing. Update Purpose after archive.
## Requirements
### Requirement: Specs tab displays a list of specs
The specs tab SHALL display a vertical list of spec names on the left side of the detail panel.

#### Scenario: Project with specs
- **WHEN** the user switches to the specs tab for a project that has specs
- **THEN** the left side SHALL list all spec directory names sorted alphabetically

#### Scenario: Project with no specs
- **WHEN** the user switches to the specs tab for a project with no specs
- **THEN** the tab SHALL display "No specs found"

### Requirement: Specs tab displays selected spec content
The specs tab SHALL display the content of the selected spec's `spec.md` file on the right side, rendered as styled markdown.

#### Scenario: Spec selected
- **WHEN** a spec is highlighted in the list
- **THEN** the right side SHALL show the spec.md content rendered with markdown styling (headers, lists, bold, italic)

#### Scenario: Spec without spec.md
- **WHEN** a spec directory exists but has no spec.md file
- **THEN** the right side SHALL display "No spec.md found"

### Requirement: Spec list navigation
The user SHALL be able to navigate the spec list with j/k keys when the specs tab is active.

#### Scenario: Navigate specs with j/k
- **WHEN** the detail panel is focused, the specs tab is active, and the user presses j or k
- **THEN** the spec cursor SHALL move down or up and the content panel SHALL update to show the newly selected spec

#### Scenario: Cursor bounds
- **WHEN** the spec cursor is at the first or last spec
- **THEN** pressing k or j respectively SHALL not move the cursor beyond the bounds

