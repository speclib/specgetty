## ADDED Requirements

### Requirement: Auto-rescan on filesystem changes in zoom mode
The system SHALL automatically rescan the zoomed project when any file or directory within its `openspec/` directory tree is created, modified, or deleted.

#### Scenario: File modified in openspec directory
- **WHEN** the user is in zoom mode AND a file within the project's `openspec/` directory is modified
- **THEN** the TUI SHALL trigger a single-project rescan and update the displayed content

#### Scenario: New file created in openspec directory
- **WHEN** the user is in zoom mode AND a new file is created within the project's `openspec/` directory
- **THEN** the TUI SHALL trigger a rescan and reflect the new file in the view

#### Scenario: File deleted in openspec directory
- **WHEN** the user is in zoom mode AND a file is removed from the project's `openspec/` directory
- **THEN** the TUI SHALL trigger a rescan and reflect the deletion in the view

### Requirement: Debounce rapid filesystem changes
The system SHALL debounce rapid filesystem events so that multiple changes within a short window (200ms) result in a single rescan.

#### Scenario: Burst of file writes
- **WHEN** multiple files are written within 200ms (e.g., editor save operation)
- **THEN** the system SHALL perform exactly one rescan after the burst settles

### Requirement: Watcher lifecycle tied to zoom mode
The system SHALL start filesystem watching when entering zoom mode and stop watching when exiting zoom mode.

#### Scenario: Enter zoom mode
- **WHEN** the user enters zoom mode (via Enter key or --zoom CLI flag)
- **THEN** the system SHALL begin watching the zoomed project's `openspec/` directory tree for changes

#### Scenario: Exit zoom mode
- **WHEN** the user exits zoom mode (via Escape or Enter key)
- **THEN** the system SHALL stop the filesystem watcher and release all watch handles

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
