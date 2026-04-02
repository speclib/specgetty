## ADDED Requirements

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
