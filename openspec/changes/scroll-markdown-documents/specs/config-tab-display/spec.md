## ADDED Requirements

### Requirement: Config tab content is a document viewer
The content of the config tab, whether it is styled markdown or highlighted
YAML, SHALL be a document viewer, with the wrapping, scrolling, position
reporting and position retention that capability describes.

#### Scenario: Configuration longer than the panel
- **WHEN** the config tab shows a `project.md` or `config.yaml` longer than the
  panel
- **THEN** the user SHALL be able to reach the end of it with the keyboard, and
  the panel title SHALL report the reading position

#### Scenario: Leaving and returning to the tab
- **WHEN** the config tab has been scrolled and the user switches to another tab
  and back, without changing project
- **THEN** the content SHALL be shown at the position it was left at

#### Scenario: A different project is selected
- **WHEN** the config tab has been scrolled and the user selects a different
  project
- **THEN** that project's configuration SHALL be shown from its first row

#### Scenario: The file source indicator stays put
- **WHEN** the config tab content is scrolled
- **THEN** the dimmed line naming the file SHALL remain visible at the top of
  the tab
