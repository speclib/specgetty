## ADDED Requirements

### Requirement: The store registry is watched where it decides resolution
A project that declares a store resolves through the registry, which lives
outside every `openspec/` tree and so changes without any watched file changing.
The system SHALL watch the registry for such a project, so that a store
registered, removed or repointed while the project is open is noticed the way an
edit to the declaration itself already is.

A project that declares no store SHALL NOT watch the registry, the registry
being unable to change what such a project resolves to.

#### Scenario: A store registered while the project is open
- **GIVEN** an open project whose store declaration could not be resolved,
  because the store was not registered
- **WHEN** the store is registered
- **THEN** the project SHALL be read again, and its specs and changes SHALL be
  shown rather than the unresolved declaration

#### Scenario: A store repointed while the project is open
- **GIVEN** an open project resolved through a store
- **WHEN** that store's registered path is changed to another directory
- **THEN** the project SHALL be read again from the new path, and the watch
  SHALL move with it

#### Scenario: A project holding its own content
- **GIVEN** an open project that declares no store
- **WHEN** the watch set is established
- **THEN** it SHALL NOT include the registry

### Requirement: A change arriving during a scan is not lost
The watcher SHALL never block on a consumer that is busy, and a signal it cannot
deliver SHALL NOT be discarded. A filesystem change that settles while a scan is
already in flight SHALL cause another scan once that one completes.

#### Scenario: A change during a scan
- **GIVEN** a scan of the open project is in flight
- **WHEN** a file under a watched tree is written and the debounce settles
- **THEN** a further scan SHALL follow the one in flight, and the display SHALL
  come to show the written content

#### Scenario: The watcher is never blocked
- **WHEN** a signal cannot be delivered because the previous one is unconsumed
- **THEN** the watcher SHALL continue reading filesystem events rather than
  waiting for the consumer

#### Scenario: A quiet project scans once
- **GIVEN** no filesystem change during a scan
- **WHEN** that scan completes
- **THEN** no further scan SHALL follow it
