## ADDED Requirements

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
The system SHALL watch the open project's `openspec/` directory tree, and SHALL
watch exactly one project at a time.

#### Scenario: Project opened at startup
- **WHEN** the application resolves a project at startup
- **THEN** it SHALL begin watching that project's `openspec/` directory tree

#### Scenario: Switching project in the picker
- **WHEN** the user selects a different project in the picker
- **THEN** watching SHALL stop for the previous project before it starts for the
  newly selected one

#### Scenario: Application exit
- **WHEN** the application exits
- **THEN** the watcher SHALL be closed

## REMOVED Requirements

### Requirement: Auto-rescan on filesystem changes in zoom mode
**Reason**: Rephrased without the zoom vocabulary, and extended with the filter
survival that `change-search` relies on. Replaced by "Auto-rescan on filesystem
changes in the open project".

### Requirement: Watcher lifecycle tied to zoom mode
**Reason**: There is no zoom mode to enter or leave. The lifecycle now follows
which project is open, including switching through the picker, which is a path
that did not previously exist. Replaced by "Watcher lifecycle follows the open
project".
