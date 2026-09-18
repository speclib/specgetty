## ADDED Requirements

### Requirement: A scan directory that cannot be read never ends the process
The scanner SHALL handle an unreadable scan directory the same way whether it is
reached by expanding a glob include or by walking a plain include. When directory
errors are ignored, which is the default, the scanner SHALL log the error, skip
that include, and carry on with the rest. When they are not ignored, the scanner
SHALL return the error to its caller. In neither case SHALL the scanner terminate
the process.

#### Scenario: A glob include points into a directory that does not exist
- **WHEN** an include ends in `*`, its parent directory cannot be read, and
  directory errors are ignored
- **THEN** the scanner SHALL log the error, skip that include, and report the
  projects found under the remaining includes

#### Scenario: The same glob with directory errors not ignored
- **WHEN** an include ends in `*`, its parent directory cannot be read, and
  directory errors are not ignored
- **THEN** the scanner SHALL return the error, close its results channel, and
  leave the process running

### Requirement: An unusable scandirs entry is ignored rather than fatal
The scanner SHALL ignore an empty entry in `scandirs.include` or
`scandirs.exclude`. An empty entry SHALL match nothing and SHALL NOT panic.

#### Scenario: An empty exclude entry
- **WHEN** `scandirs.exclude` contains an empty string, which is what a YAML dash
  with nothing after it parses to
- **THEN** the entry SHALL exclude nothing and the scan SHALL complete

#### Scenario: An empty include entry
- **WHEN** `scandirs.include` contains an empty string
- **THEN** the entry SHALL be skipped and the scan SHALL complete over the
  remaining includes
