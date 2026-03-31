## MODIFIED Requirements

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
