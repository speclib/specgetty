# fs-watch Specification

## Purpose
TBD - created by archiving change inotify-auto-rescan. Update Purpose after archive.

## Requirements

### Requirement: Debounce rapid filesystem changes
The system SHALL debounce rapid filesystem events so that multiple changes within a short window (200ms) result in a single rescan.

#### Scenario: Burst of file writes
- **WHEN** multiple files are written within 200ms (e.g., editor save operation)
- **THEN** the system SHALL perform exactly one rescan after the burst settles

### Requirement: Watch new subdirectories dynamically
The system SHALL detect newly created subdirectories within the watched `openspec/` tree and add them to the watch set.

#### Scenario: New change directory created
- **WHEN** a new subdirectory is created under `openspec/changes/` while watching
- **THEN** the system SHALL add the new subdirectory to the watch set so that files within it are also monitored

### Requirement: Cross-platform filesystem watching
The system SHALL support filesystem watching on both Linux (via inotify) and macOS (via kqueue) using the fsnotify library.

#### Scenario: Running on Linux
- **WHEN** the application runs on Linux
- **THEN** filesystem watching SHALL use inotify and function correctly

#### Scenario: Running on macOS
- **WHEN** the application runs on macOS
- **THEN** filesystem watching SHALL use kqueue and function correctly

### Requirement: Auto-rescan on filesystem changes in the open project
The system SHALL automatically rescan the open project when any file or
directory within its `openspec/` directory tree is created, modified, or
deleted.

#### Scenario: File modified in the open project
- **WHEN** a project is open AND a file within its `openspec/` directory is
  modified
- **THEN** the system SHALL rescan that project and refresh the display

#### Scenario: File created in the open project
- **WHEN** a project is open AND a new file is created within its `openspec/`
  directory
- **THEN** the system SHALL rescan that project and refresh the display

#### Scenario: File deleted in the open project
- **WHEN** a project is open AND a file is removed from its `openspec/`
  directory
- **THEN** the system SHALL rescan that project and refresh the display

#### Scenario: An active filter survives the rescan
- **WHEN** a change list filter is active and a rescan is triggered
- **THEN** the filter SHALL remain applied afterwards

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
