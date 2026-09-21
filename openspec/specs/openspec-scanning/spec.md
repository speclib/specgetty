# openspec-scanning Specification

## Purpose
TBD - created by archiving change replace-git-with-openspec-scanning. Update Purpose after archive.

## Requirements

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

### Requirement: Display projects in TUI
The TUI SHALL present discovered OpenSpec projects in the project picker, sorted
alphabetically by path.

#### Scenario: Projects found
- **WHEN** projects have been discovered and the picker is opened
- **THEN** the picker SHALL list them, and selecting one SHALL open it in the
  project view

#### Scenario: No projects found
- **WHEN** no OpenSpec projects have been discovered
- **THEN** the picker SHALL display a message indicating none were found

### Requirement: Remove go-git dependency
The application SHALL NOT depend on `github.com/go-git/go-git/v5` or any git-specific libraries.

#### Scenario: Clean build without go-git
- **WHEN** the application is built
- **THEN** `go.mod` SHALL NOT contain `github.com/go-git/go-git/v5` as a direct or indirect dependency

### Requirement: A scan directory that cannot be read never ends the process
The scanner SHALL handle an unreadable scan directory the same way whether it is
reached by expanding a glob include or by walking a plain include. When directory
errors are ignored, which is the default, the scanner SHALL log the error, skip
that include, and carry on with the rest. When they are not ignored, the scanner
SHALL return the error to its caller. In neither case SHALL the scanner terminate
the process.

#### Scenario: A glob include points into a directory that does not exist
- **WHEN** an include ends in `*`, its parent directory cannot be read, and
  directory errors are ignored
- **THEN** the scanner SHALL log the error, skip that include, and report the
  projects found under the remaining includes

#### Scenario: The same glob with directory errors not ignored
- **WHEN** an include ends in `*`, its parent directory cannot be read, and
  directory errors are not ignored
- **THEN** the scanner SHALL return the error, close its results channel, and
  leave the process running

### Requirement: An unusable scandirs entry is ignored rather than fatal
The scanner SHALL ignore an empty entry in `scandirs.include` or
`scandirs.exclude`. An empty entry SHALL match nothing and SHALL NOT panic.

#### Scenario: An empty exclude entry
- **WHEN** `scandirs.exclude` contains an empty string, which is what a YAML dash
  with nothing after it parses to
- **THEN** the entry SHALL exclude nothing and the scan SHALL complete

#### Scenario: An empty include entry
- **WHEN** `scandirs.include` contains an empty string
- **THEN** the entry SHALL be skipped and the scan SHALL complete over the
  remaining includes

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

### Requirement: A glob include contributes only what it expands to
A `scandirs.include` entry ending in `*` SHALL contribute the directories it
matches and SHALL NOT be walked as a path in its own right.

The pattern is not a directory. Walking it fails on every scan, which was
logged and swallowed while a log panel existed to swallow it, and which fails
the whole scan when directory errors are not ignored.

#### Scenario: A glob that matches directories
- **WHEN** an include ends in `*` and matches directories beside it
- **THEN** those directories SHALL be walked and the pattern itself SHALL NOT be

#### Scenario: A glob that matches nothing
- **WHEN** an include ends in `*` and matches no directory
- **THEN** nothing SHALL be walked for that include, and no error SHALL be
  reported for the pattern

#### Scenario: A plain include is unaffected
- **WHEN** an include does not end in `*`
- **THEN** it SHALL be walked as the path it is

#### Scenario: A glob with directory errors not ignored
- **WHEN** an include ends in `*`, it matches directories, and directory errors
  are not ignored
- **THEN** the scan SHALL complete rather than failing on the pattern, because
  the pattern is no longer walked
