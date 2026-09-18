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
