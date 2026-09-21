## ADDED Requirements

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

## REMOVED Requirements

### Requirement: A store is identified by its metadata file
**Reason**: Superseded by "A store is a directory the registry resolves to".
Trusting the metadata file alone let an unregistered copy wear the registered
store's name: `gh.nivis-project/ospecs` carries `id: nivis` while the registered
`nivis` is `nivis-openspec-stores/nivis`, and both were shown as `nivis`. The
registry is what decides, which is also how the declaration path already worked.
