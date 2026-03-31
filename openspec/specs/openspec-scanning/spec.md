# openspec-scanning Specification

## Purpose
TBD - created by archiving change replace-git-with-openspec-scanning. Update Purpose after archive.
## Requirements
### Requirement: Detect OpenSpec projects by directory
The scanner SHALL identify a directory as an OpenSpec project when it contains a direct child directory named `openspec` that passes validation. A valid `openspec/` directory MUST contain at least one of `config.yaml` or `project.md`, AND at least one of the subdirectories `specs/` or `archive/`.

#### Scenario: Valid openspec directory with config.yaml and specs/
- **WHEN** the walker encounters a directory containing `openspec/` with `config.yaml` and `specs/` inside
- **THEN** the parent directory SHALL be reported as an OpenSpec project

#### Scenario: Valid openspec directory with project.md and archive/
- **WHEN** the walker encounters a directory containing `openspec/` with `project.md` and `archive/` inside
- **THEN** the parent directory SHALL be reported as an OpenSpec project

#### Scenario: Directory named openspec without required contents
- **WHEN** the walker encounters a directory named `openspec` that does not contain (`config.yaml` or `project.md`) AND (`specs/` or `archive/`)
- **THEN** the directory SHALL NOT be reported as an OpenSpec project

#### Scenario: openspec directory inside .git
- **WHEN** a directory named `openspec` exists inside a `.git/` subtree
- **THEN** it SHALL NOT be reported as an OpenSpec project (it will fail validation)

#### Scenario: openspec directory in test fixtures
- **WHEN** a directory named `openspec` exists inside a test fixture directory without proper structure
- **THEN** it SHALL NOT be reported as an OpenSpec project

### Requirement: List OpenSpec directory contents
The scanner SHALL read the contents of each detected project's `openspec/` directory recursively and produce a flat list of entries.

#### Scenario: Project with openspec contents
- **WHEN** an OpenSpec project is detected
- **THEN** the scanner SHALL list all files and subdirectories under `openspec/` with paths relative to the `openspec/` directory

#### Scenario: Empty openspec directory
- **WHEN** an OpenSpec project's `openspec/` directory is empty
- **THEN** the project SHALL still appear in the project list with an empty file list

### Requirement: Display projects in TUI
The TUI SHALL display detected OpenSpec projects in the left panel, sorted alphabetically by path.

#### Scenario: Projects found
- **WHEN** the scan completes with one or more OpenSpec projects found
- **THEN** the left panel SHALL list all project paths, and selecting a project SHALL show its openspec contents in the right panel

#### Scenario: No projects found
- **WHEN** the scan completes with zero OpenSpec projects found
- **THEN** the left panel SHALL display a message indicating no projects were found

### Requirement: Display OpenSpec contents without git status
The detail panel SHALL display the contents of the selected project's `openspec/` directory as a plain file/directory listing without git status codes or diff views.

#### Scenario: Viewing project contents
- **WHEN** a project is selected in the left panel
- **THEN** the detail panel SHALL show each entry with a directory indicator (`d`) or file indicator (`f`) followed by the relative path

#### Scenario: No diff panel
- **WHEN** a project is selected and files are listed
- **THEN** there SHALL be no diff viewport or git diff functionality available

### Requirement: Remove go-git dependency
The application SHALL NOT depend on `github.com/go-git/go-git/v5` or any git-specific libraries.

#### Scenario: Clean build without go-git
- **WHEN** the application is built
- **THEN** `go.mod` SHALL NOT contain `github.com/go-git/go-git/v5` as a direct or indirect dependency

