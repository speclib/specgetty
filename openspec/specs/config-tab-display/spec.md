# config-tab-display Specification

## Purpose
TBD - created by archiving change config-tab-project-display. Update Purpose after archive.

## Requirements

### Requirement: Config tab displays project.md as styled markdown
When an OpenSpec project has a `project.md` file, the config tab SHALL render it with basic markdown styling.

#### Scenario: Project with project.md
- **WHEN** user switches to the config tab for a project that has `openspec/project.md`
- **THEN** the content SHALL be displayed with headers styled bold and colored, list items indented, and bold/italic text styled appropriately

#### Scenario: Markdown headers
- **WHEN** a line starts with `#`, `##`, or `###`
- **THEN** it SHALL be rendered in bold with a distinct color

#### Scenario: Markdown list items
- **WHEN** a line starts with `- ` or `* `
- **THEN** it SHALL be rendered with proper indentation

### Requirement: Config tab displays config.yaml with syntax highlighting
When an OpenSpec project has a `config.yaml` file (and no project.md), the config tab SHALL render it with YAML syntax highlighting.

#### Scenario: Project with config.yaml only
- **WHEN** user switches to the config tab for a project that has `openspec/config.yaml` but no `openspec/project.md`
- **THEN** the YAML content SHALL be displayed with keys in one color, values in another, and comments dimmed

#### Scenario: YAML comments
- **WHEN** a line contains a `#` comment
- **THEN** the comment portion SHALL be rendered in a dimmed style

### Requirement: Config tab shows file source indicator
The config tab SHALL show which file is being displayed.

#### Scenario: File indicator
- **WHEN** the config tab is displayed
- **THEN** a dimmed line at the top SHALL indicate the file path (e.g. "openspec/project.md" or "openspec/config.yaml")

### Requirement: Config tab handles missing configuration
When no configuration file exists, the config tab SHALL show an appropriate message.

#### Scenario: No config file
- **WHEN** a project has neither `openspec/project.md` nor `openspec/config.yaml`
- **THEN** the config tab SHALL display "No project configuration found"

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
