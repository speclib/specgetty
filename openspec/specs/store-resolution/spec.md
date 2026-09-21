# store-resolution Specification

## Purpose
TBD - created by archiving change follow-the-store. Update Purpose after archive.

## Requirements

### Requirement: The root is resolved from a starting directory
The application SHALL resolve, from a starting directory, the single root whose
`openspec/specs/` and `openspec/changes/` hold the content to display. The
resolution SHALL walk up from the starting directory to the nearest qualifying
`openspec/` directory, where a directory qualifies when it holds `specs/` or
`changes/`, or holds a configuration file.

#### Scenario: The nearest directory holds specs or changes
- **GIVEN** a directory whose `openspec/` holds `specs/` or `changes/`
- **WHEN** the root is resolved from it or from any directory beneath it
- **THEN** that directory SHALL be the root

#### Scenario: A pointer beside specs or changes is ignored
- **GIVEN** a directory whose `openspec/` holds `specs/` or `changes/` and whose
  configuration also declares a store
- **WHEN** the root is resolved from it
- **THEN** that directory SHALL be the root and the declaration SHALL NOT be
  followed, because content present locally outranks a pointer

#### Scenario: A configuration with no content and no pointer
- **GIVEN** a directory whose `openspec/` holds a configuration file but neither
  `specs/` nor `changes/` nor a store declaration
- **WHEN** the root is resolved from it
- **THEN** that directory SHALL be the root, which is the shape of a freshly
  initialised project

#### Scenario: A directory named openspec with nothing in it
- **GIVEN** a directory holding an `openspec/` directory with no configuration
  file, no `specs/` and no `changes/`
- **WHEN** the root is resolved from a directory beneath it
- **THEN** that directory SHALL NOT qualify and the walk SHALL continue upward

### Requirement: A store declaration resolves through the registry
When the nearest qualifying `openspec/` directory holds no `specs/` and no
`changes/` and its configuration declares `store: <id>`, the application SHALL
resolve the root to the local path the store registry records for that id.

#### Scenario: A repo that points at a store
- **GIVEN** a repo whose `openspec/config.yaml` holds `store: nivis-tunnel` and
  no `specs/` or `changes/` directory
- **AND** a registry entry for `nivis-tunnel` naming a local path
- **WHEN** the root is resolved from that repo
- **THEN** the root SHALL be that local path, and the repo SHALL be recorded as
  the origin the resolution started from

#### Scenario: Two repos pointing at one store
- **GIVEN** two repos whose configurations both declare the same store id
- **WHEN** the root is resolved from each
- **THEN** both SHALL resolve to the same root, and each SHALL keep its own
  origin

### Requirement: The store registry is read from the OpenSpec data directory
The application SHALL read the store registry from
`<data dir>/openspec/stores/registry.yaml`, where the data directory is
`$XDG_DATA_HOME` when that variable is set, and `~/.local/share` otherwise.

#### Scenario: XDG_DATA_HOME is set
- **GIVEN** `XDG_DATA_HOME` is set
- **WHEN** the registry is read
- **THEN** it SHALL be read from beneath that directory

#### Scenario: XDG_DATA_HOME is not set
- **GIVEN** `XDG_DATA_HOME` is not set
- **WHEN** the registry is read
- **THEN** it SHALL be read from `~/.local/share/openspec/stores/registry.yaml`

#### Scenario: The registry entry names where the store is
- **WHEN** a registry entry is read
- **THEN** the store's root SHALL be taken from the entry's `backend.local_path`

### Requirement: Resolution reads only local files
Resolving a root SHALL read only the local filesystem. The application SHALL NOT
clone, fetch or contact any network host, and SHALL NOT invoke the `openspec`
binary to resolve a root.

#### Scenario: A registry entry recording a remote
- **GIVEN** a registry entry carrying a `remote` alongside its `local_path`
- **WHEN** the root is resolved
- **THEN** the `local_path` SHALL be used and the `remote` SHALL NOT be contacted

#### Scenario: No openspec binary installed
- **GIVEN** no `openspec` binary on the path
- **WHEN** a store-backed project is opened
- **THEN** its content SHALL be displayed, because resolution does not depend on
  that binary

### Requirement: Both configuration file spellings are read
Wherever the application reads `openspec/config.yaml` it SHALL also accept
`openspec/config.yml`, preferring `.yaml` when both are present.

#### Scenario: A project using the yml spelling
- **GIVEN** a project with `openspec/config.yml` and no `openspec/config.yaml`
- **WHEN** it is scanned
- **THEN** it SHALL be recognised, and any store declaration in that file SHALL
  be followed

### Requirement: A declaration that cannot be resolved is reported
When a store declaration cannot be followed, the application SHALL say which
declaration failed and why, rather than presenting the project as empty.

#### Scenario: The declared store is not registered
- **GIVEN** a configuration declaring a store id with no entry in the registry
- **WHEN** the project is opened
- **THEN** the application SHALL report the unresolved id and the file that
  declared it

#### Scenario: No registry on the machine
- **GIVEN** a configuration declaring a store and no registry file at all
- **WHEN** the project is opened
- **THEN** the application SHALL report that no stores are registered

#### Scenario: The declaration is not a single id
- **GIVEN** a configuration whose `store` key holds a list or a mapping rather
  than a string
- **WHEN** the project is opened
- **THEN** the application SHALL report the declaration as malformed

#### Scenario: The registered path is gone
- **GIVEN** a registry entry whose `local_path` no longer exists
- **WHEN** the project is opened
- **THEN** the application SHALL report the missing path rather than showing
  zero specs and zero changes

### Requirement: Every filesystem operation targets the resolved root
Every operation that reads or writes a project's content SHALL act on the
resolved root. None of them SHALL act on the starting directory when the two
differ.

This is stated over every operation rather than over a list of them. A list has
to be joined by each operation added later, and two were not: the task toggle and
the key that copies a change's path both built their paths from the starting
directory, and both were wrong in a store-backed project for as long as that
list stood.

#### Scenario: Reading a store-backed project
- **GIVEN** a project resolved to a store
- **WHEN** its specs and changes are read
- **THEN** they SHALL be read from the store's `openspec/` directory

#### Scenario: Discarding a change in a store-backed project
- **GIVEN** a change open in a project resolved to a store
- **WHEN** the change is discarded
- **THEN** it SHALL be moved within the store's `openspec/changes/`, and no
  directory SHALL be created under the pointing repo

#### Scenario: Exporting a change in a store-backed project
- **GIVEN** a change in a project resolved to a store
- **WHEN** it is exported
- **THEN** the change directory SHALL be read from the store

#### Scenario: Writing a task back in a store-backed project
- **GIVEN** a change open in a project resolved to a store
- **WHEN** a task is toggled
- **THEN** the `tasks.md` under the store SHALL be written, and the toggle SHALL
  succeed

#### Scenario: Reporting a path in a store-backed project
- **GIVEN** a change selected in a project resolved to a store
- **WHEN** a path to it is put on the clipboard or handed to another program
- **THEN** that path SHALL be under the store, and SHALL resolve to something
  that exists

#### Scenario: An operation added later
- **WHEN** an operation that reads or writes project content is added to the
  application
- **THEN** it SHALL act on the resolved root, this requirement covering it
  without being amended to name it

### Requirement: A store is a directory the registry resolves to
A directory SHALL be treated as a store only when it holds
`.openspec-store/store.yaml` AND the store registry resolves the id in that file
to that same directory. An identity file alone SHALL NOT make a directory a
store: a clone, or a predecessor the registry has moved on from, keeps its copy
of that file and would otherwise wear the registered store's name.

#### Scenario: A registered store
- **GIVEN** a directory holding `.openspec-store/store.yaml` with `id: alpha`
- **AND** a registry entry for `alpha` naming that directory
- **THEN** it SHALL be reported as a store named `alpha`, whatever its directory
  is called

#### Scenario: An identity file the registry points elsewhere
- **GIVEN** a directory holding `.openspec-store/store.yaml` with `id: alpha`
- **AND** a registry whose `alpha` names a different directory
- **THEN** it SHALL NOT be reported as a store, and SHALL be named by its own
  directory rather than by the id it claims

#### Scenario: An identity file with no registry entry at all
- **GIVEN** a directory holding `.openspec-store/store.yaml` and a registry that
  does not mention its id
- **THEN** it SHALL NOT be reported as a store

#### Scenario: No registry on the machine
- **GIVEN** a directory holding `.openspec-store/store.yaml` and no registry file
- **THEN** it SHALL NOT be reported as a store

#### Scenario: A root with no store metadata
- **GIVEN** a discovered root with no `.openspec-store/` directory
- **WHEN** it is reported
- **THEN** it SHALL be reported as a plain project, named as projects are named

#### Scenario: Following a declaration is unaffected
- **GIVEN** a repo declaring a store the registry does resolve
- **WHEN** the root is resolved from it
- **THEN** it SHALL resolve to that store, whose identity the registry and the
  metadata file have already been made to agree on
