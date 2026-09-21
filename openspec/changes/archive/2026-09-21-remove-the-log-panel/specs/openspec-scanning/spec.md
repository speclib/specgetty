## ADDED Requirements

### Requirement: A glob include contributes only what it expands to
A `scandirs.include` entry ending in `*` SHALL contribute the directories it
matches and SHALL NOT be walked as a path in its own right.

The pattern is not a directory. Walking it fails on every scan, which was
logged and swallowed while a log panel existed to swallow it, and which fails
the whole scan when directory errors are not ignored.

#### Scenario: A glob that matches directories
- **WHEN** an include ends in `*` and matches directories beside it
- **THEN** those directories SHALL be walked and the pattern itself SHALL NOT be

#### Scenario: A glob that matches nothing
- **WHEN** an include ends in `*` and matches no directory
- **THEN** nothing SHALL be walked for that include, and no error SHALL be
  reported for the pattern

#### Scenario: A plain include is unaffected
- **WHEN** an include does not end in `*`
- **THEN** it SHALL be walked as the path it is

#### Scenario: A glob with directory errors not ignored
- **WHEN** an include ends in `*`, it matches directories, and directory errors
  are not ignored
- **THEN** the scan SHALL complete rather than failing on the pattern, because
  the pattern is no longer walked
