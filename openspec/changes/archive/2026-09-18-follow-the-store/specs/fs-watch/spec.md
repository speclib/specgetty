## MODIFIED Requirements

### Requirement: Watcher lifecycle follows the open project
The system SHALL watch the open project's content tree, and SHALL watch exactly
one project at a time. Where the content was resolved through a store
declaration, the system SHALL watch two trees: the store's `openspec/`, where
the content lives, and the originating repo's `openspec/`, where the declaration
that points at it lives.

#### Scenario: Project opened at startup
- **WHEN** the application resolves a project at startup
- **THEN** it SHALL begin watching that project's `openspec/` directory tree

#### Scenario: A store-backed project opened at startup
- **WHEN** the application resolves a project through a store declaration
- **THEN** it SHALL watch both the store's `openspec/` tree and the originating
  repo's `openspec/` tree

#### Scenario: The declaration changes
- **WHEN** the originating repo's configuration file is edited to name a
  different store
- **THEN** the change SHALL be noticed, the root SHALL be resolved again, and the
  watched trees SHALL follow the new resolution

#### Scenario: Switching project in the picker
- **WHEN** the user selects a different project in the picker
- **THEN** watching SHALL stop for every tree of the previous project before it
  starts for the newly selected one

#### Scenario: One tree when the two coincide
- **WHEN** the open project holds its own specs and changes
- **THEN** exactly one tree SHALL be watched, as today

#### Scenario: Application exit
- **WHEN** the application exits
- **THEN** every watcher SHALL be closed
