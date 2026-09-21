## MODIFIED Requirements

### Requirement: Detect OpenSpec projects by directory
The scanner SHALL identify a directory as an OpenSpec project when it contains a
direct child directory named `openspec` holding at least one of `config.yaml`,
`config.yml` or `project.md`. Content of its own is NOT required: a repo that
declares a `store:` and keeps no `specs/` or `changes/` is still a project,
because it is the directory a person works in and looks for.

Such a repo's configuration is read by OpenSpec for one key, `store:`. Its
`schema`, `context`, `rules` and `operations` are inert, and the resolved root's
configuration is what applies. An earlier reading of this requirement said the
repo was the only place its own context and rules could be read from, which is
true and misleading: they can be read there and they are never used.

#### Scenario: Valid openspec directory with config.yaml and specs/
- **WHEN** the walker encounters a directory containing `openspec/` with `config.yaml` and `specs/` inside
- **THEN** the parent directory SHALL be reported as an OpenSpec project

#### Scenario: Valid openspec directory with project.md and archive/
- **WHEN** the walker encounters a directory containing `openspec/` with `project.md` and `archive/` inside
- **THEN** the parent directory SHALL be reported as an OpenSpec project

#### Scenario: A repo that declares a store and keeps no content
- **WHEN** the walker encounters a directory whose `openspec/` holds only a
  configuration declaring `store: <id>`
- **THEN** the parent directory SHALL be reported as an OpenSpec project

#### Scenario: Two repos sharing one store
- **WHEN** two directories under the scan directories declare the same store id
- **THEN** both SHALL be reported, each under its own directory name

#### Scenario: Directory named openspec without required contents
- **WHEN** the walker encounters a directory named `openspec` holding none of
  `config.yaml`, `config.yml` or `project.md`
- **THEN** the directory SHALL NOT be reported as an OpenSpec project, whatever
  else it contains

#### Scenario: openspec directory inside .git
- **WHEN** a directory named `openspec` exists inside a `.git/` subtree
- **THEN** it SHALL NOT be reported as an OpenSpec project (it will fail validation)

#### Scenario: openspec directory in test fixtures
- **WHEN** a directory named `openspec` exists inside a test fixture directory without proper structure
- **THEN** it SHALL NOT be reported as an OpenSpec project
