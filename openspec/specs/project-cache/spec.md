# project-cache Specification

## Purpose
Remembers which project paths a scan discovered, so that only the first run pays
for walking the configured directories.

## Requirements

### Requirement: Discovered paths are cached on disk
The set of discovered project paths SHALL be written to a cache file and reused
on later runs.

#### Scenario: First run
- **WHEN** no cache file exists and the projects are needed
- **THEN** the configured scan directories SHALL be walked and the discovered
  paths SHALL be written to the cache

#### Scenario: Later runs
- **WHEN** a valid cache file exists and the projects are needed
- **THEN** the paths SHALL be read from the cache and the configured scan
  directories SHALL NOT be walked

#### Scenario: Cache location
- **WHEN** the cache is written
- **THEN** it SHALL be stored under the user's cache directory, following the
  XDG Base Directory Specification

### Requirement: Only paths are cached
The cache SHALL hold discovered paths, and SHALL NOT hold project statistics.

#### Scenario: Statistics are not served from the cache
- **WHEN** projects are loaded from the cache
- **THEN** each project's specs, changes and task counts SHALL be read from disk
  rather than from the cache

### Requirement: Cache is invalidated when the question changes
The cache SHALL record the scan directories it was built from, and SHALL be
discarded when those no longer match.

#### Scenario: Scan directories edited
- **WHEN** the configured include or exclude directories differ from those
  recorded in the cache
- **THEN** the cache SHALL be treated as absent and a full walk SHALL be
  performed

#### Scenario: Unreadable or malformed cache
- **WHEN** the cache file cannot be read or does not parse
- **THEN** it SHALL be treated as absent rather than failing the application

### Requirement: Vanished projects are dropped without a walk
Cached paths that no longer exist SHALL be omitted.

#### Scenario: A cached project was deleted
- **WHEN** the cache is loaded and one of its paths no longer exists on disk
- **THEN** that path SHALL be omitted from the list, without walking the scan
  directories

### Requirement: New projects appear on refresh
A project created since the last walk SHALL appear after an explicit refresh.

#### Scenario: A project was created
- **WHEN** a new OpenSpec project exists under the configured directories and the
  cache predates it
- **THEN** it SHALL NOT be listed until the user refreshes

#### Scenario: After refreshing
- **WHEN** the user refreshes
- **THEN** the newly created project SHALL be listed and the cache SHALL be
  rewritten
