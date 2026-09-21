## ADDED Requirements

### Requirement: A registered store is not listed as a project
A directory the registry resolves a store id to SHALL NOT be reported by the
scanner. A store is where content lives, not where work happens, and every repo
reading from it already stands for it in the list.

#### Scenario: The registered stores under a scan directory
- **WHEN** the walker reaches a directory the registry resolves a store id to
- **THEN** it SHALL NOT be reported, however much content it holds

#### Scenario: A directory whose identity file the registry does not confirm
- **WHEN** the walker reaches a directory carrying `.openspec-store/store.yaml`
  that the registry does not resolve to it
- **THEN** it SHALL be reported as an ordinary project

#### Scenario: No registry on the machine
- **WHEN** no store registry exists
- **THEN** no directory SHALL be excluded as a store, and discovery SHALL be
  what it would be with no stores in play

### Requirement: A discovered project carries the statistics of its content
The scanner SHALL report each discovered project with the specs, changes and
task counts of the root its content resolves to, which for a store-backed repo
is the store's.

#### Scenario: A store-backed repo's statistics
- **WHEN** a repo declaring a store is reported
- **THEN** its spec, change and archived counts SHALL be the store's, not zero

#### Scenario: Two repos sharing one store
- **WHEN** two reported repos declare the same store id
- **THEN** both SHALL carry that store's statistics, and the store SHALL be read
  once for the scan rather than once for each repo

#### Scenario: A declaration that cannot be followed
- **WHEN** a reported repo declares a store that cannot be resolved
- **THEN** it SHALL still be reported, carrying the problem rather than being
  dropped from the list

## MODIFIED Requirements

### Requirement: Detect OpenSpec projects by directory
The scanner SHALL identify a directory as an OpenSpec project when it contains a
direct child directory named `openspec` holding at least one of `config.yaml`,
`config.yml` or `project.md`. Content of its own is NOT required: a repo that
declares a `store:` and keeps no `specs/` or `changes/` is still a project, and
is the only place that repo's own context and rules can be read from.

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

### Requirement: List OpenSpec directory contents
The scanner SHALL read the contents of each detected project's `openspec/`
directory recursively and produce a flat list of entries, reading from the root
the project's content resolves to.

#### Scenario: Project with openspec contents
- **WHEN** an OpenSpec project is detected
- **THEN** the scanner SHALL list all files and subdirectories under `openspec/`
  with paths relative to the `openspec/` directory

#### Scenario: A store-backed repo's contents
- **WHEN** a repo declaring a store is detected
- **THEN** the listing SHALL be the store's `openspec/` tree, so that a search
  inside a project reaches the specs it actually shows

#### Scenario: Empty openspec directory
- **WHEN** an OpenSpec project's `openspec/` directory is empty
- **THEN** the project SHALL still appear in the project list with an empty file
  list

## REMOVED Requirements

### Requirement: A discovered root reports whether it is a store
**Reason**: Superseded. The scanner no longer reports stores at all, so there is
nothing for a discovered row to be a store of. What replaces it is the rule
above, that a registered store is not listed, and the identity check in
`store-resolution` that decides what counts as one.

### Requirement: A repo that only points at a store is not a row of its own
**Reason**: Reversed. Listing the store instead of the repos that read from it
hid every directory a person actually works in, and made a repo's own context
and rules unreachable from the picker, since only an origin gives the config tab
that pane. Replaced by the rule that such a repo is listed in its own right.
