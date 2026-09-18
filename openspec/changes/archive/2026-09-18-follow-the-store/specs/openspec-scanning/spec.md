## ADDED Requirements

### Requirement: A discovered root reports whether it is a store
The scanner SHALL report, for each discovered root, whether that root is a store
and under which id, by reading `.openspec-store/store.yaml` beside its
`openspec/` directory.

#### Scenario: A store among the scanned directories
- **WHEN** a discovered root holds `.openspec-store/store.yaml`
- **THEN** the scanner SHALL report it as a store carrying the `id` from that
  file

#### Scenario: An ordinary project
- **WHEN** a discovered root holds no `.openspec-store/` directory
- **THEN** the scanner SHALL report it as a plain project

#### Scenario: Store metadata that cannot be read
- **WHEN** `.openspec-store/store.yaml` exists but cannot be parsed, or names no
  id
- **THEN** the root SHALL still be reported, as a plain project, so that
  unreadable metadata hides nothing

### Requirement: A repo that only points at a store is not a row of its own
A directory whose `openspec/` holds a configuration file and neither `specs/`
nor `changes/` SHALL NOT be reported as a project, whether or not it declares a
store. Such a directory is an entry point to a root, not a root.

#### Scenario: A store-backed repo under a scan directory
- **WHEN** the walker reaches a repo whose `openspec/` holds only
  `config.yaml` declaring `store: <id>`
- **THEN** that repo SHALL NOT be reported, and the store it names SHALL be
  reported on its own account when it lies under a scan directory

#### Scenario: Two repos sharing one store
- **WHEN** two repos under the scan directories declare the same store id
- **THEN** neither SHALL be reported, and the store SHALL appear once
