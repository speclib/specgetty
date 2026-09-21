## ADDED Requirements

### Requirement: Enter opens the selected spec in detail
The specs tab SHALL treat `enter` the way the change list does: it descends into
the selected spec. The detail view it opens is described by `spec-detail-view`.

#### Scenario: Enter on a spec
- **GIVEN** the specs tab is active and its list holds the keyboard
- **WHEN** the user presses `enter`
- **THEN** the selected spec SHALL open in the spec detail view

#### Scenario: Enter while the content holds the keyboard
- **GIVEN** the specs tab is active and its content holds the keyboard
- **WHEN** the user presses `enter`
- **THEN** the selected spec SHALL open in the spec detail view, because which
  half has the keyboard does not change which spec is selected

#### Scenario: Returning from the detail view
- **WHEN** the user leaves the spec detail view with `esc`
- **THEN** the specs tab SHALL be shown again with the same spec selected and the
  spec list holding the keyboard

#### Scenario: A project with no specs
- **WHEN** the user presses `enter` on the specs tab of a project with no specs
- **THEN** nothing SHALL happen

### Requirement: A spec that cannot be parsed reports rather than opening
Not every `spec.md` follows the Purpose, requirements and scenarios shape that
the detail view reads. When the selected spec does not, `enter` SHALL leave the
cursor where it is and report why on the nav bar, in the transient one-line form
the copy keys already use. The tab SHALL go on showing the whole file as
markdown.

#### Scenario: Enter on a spec that does not fit
- **WHEN** the user presses `enter` on a spec whose file cannot be read as
  requirements and scenarios
- **THEN** the specs tab SHALL stay on screen with its cursor unmoved, and the
  nav bar SHALL report that the spec has no requirements to show

#### Scenario: The report is transient
- **WHEN** that report is on the nav bar and the user presses any key
- **THEN** the report SHALL be gone and the nav bar SHALL list the keys again

#### Scenario: The markdown view is unaffected
- **WHEN** a spec cannot be opened in the detail view
- **THEN** the specs tab SHALL still render its whole file as markdown, so no
  content becomes unreachable

### Requirement: Backticked spans are highlighted
A span between backticks in a rendered markdown document SHALL be drawn
distinctly from the prose around it, on the specs tab and in every other document
the same renderer draws.

#### Scenario: A code span in a spec
- **WHEN** a spec contains a span between backticks
- **THEN** that span SHALL be highlighted, and the backticks themselves SHALL NOT
  be shown

#### Scenario: An unclosed backtick
- **WHEN** a line contains a single backtick with no closing one
- **THEN** the line SHALL be rendered unchanged rather than swallowing the rest
  of the document
