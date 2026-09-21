# schema-inspection Specification

## Purpose
TBD - created by archiving change properties-tab. Update Purpose after archive.

## Requirements

### Requirement: The schemas a project uses are read from its changes
A change records its workflow schema in its own `.openspec.yaml`. The
application SHALL read that file for every change, active and archived, and
SHALL report the distinct schemas in use together with how many changes are on
each.

#### Scenario: A project running more than one schema
- **GIVEN** a project whose changes name two different schemas
- **WHEN** the schemas in use are reported
- **THEN** both SHALL be reported, each with the number of changes naming it

#### Scenario: A change with no metadata file
- **GIVEN** a change directory with no `.openspec.yaml`, which is what changes
  created before that file existed look like
- **WHEN** the schemas in use are reported
- **THEN** the change SHALL be counted separately as carrying no recorded
  schema, rather than being assumed to use the project default

#### Scenario: The project default is always in the set
- **GIVEN** a project whose configuration names a schema no change uses
- **WHEN** the schemas in use are reported
- **THEN** that schema SHALL still be reported, marked as the project default

#### Scenario: No schema key in the configuration
- **GIVEN** a project configuration with no `schema` key
- **WHEN** the project default is reported
- **THEN** it SHALL be `spec-driven`, shown as the default rather than as a
  value that was written down

#### Scenario: The configuration that decides
- **GIVEN** a project resolved through a store declaration
- **WHEN** the project default is read
- **THEN** it SHALL be read from the resolved root's configuration, which is the
  one OpenSpec itself reads

### Requirement: A schema definition is located by the CLI and read locally
Built-in schemas live inside the OpenSpec installation, at a path the
application cannot compute. The application SHALL obtain that path by running
`openspec schema which <name> --json`, and SHALL then read
`<path>/schema.yaml` with its own parser.

#### Scenario: A built-in schema
- **GIVEN** a schema whose definition ships with OpenSpec
- **WHEN** its details are read
- **THEN** the path SHALL come from the CLI and the definition SHALL be parsed
  from that path

#### Scenario: A schema defined by the project
- **GIVEN** a schema under the resolved root's `openspec/schemas/`
- **WHEN** its details are read
- **THEN** it SHALL be reported as belonging to the project, with its path

#### Scenario: What a schema reports
- **WHEN** a schema's definition has been read
- **THEN** the report SHALL carry its description, the ordered artifacts with
  what each generates and requires, and the apply rule with what it requires and
  tracks

#### Scenario: A schema that overrides a built-in
- **GIVEN** a project schema with the same name as one that ships with OpenSpec
- **WHEN** its details are read
- **THEN** the report SHALL say that it overrides, and name what it overrides

#### Scenario: The command is run where the content lives
- **WHEN** the CLI is invoked
- **THEN** its working directory SHALL be the resolved root, so that a
  store-backed project resolves the store's schemas

### Requirement: Locating a schema is bounded and never fatal
Running another program from inside a terminal interface can fail or hang. The
application SHALL bound the wait, SHALL survive every failure, and SHALL say
which failure occurred.

#### Scenario: The CLI is not installed
- **WHEN** no `openspec` binary is on the path
- **THEN** each schema SHALL still be reported by name and by how many changes
  use it, with its details reported as unavailable and the reason given

#### Scenario: The schema does not resolve
- **GIVEN** a configuration or a change naming a schema that does not exist
- **WHEN** its details are requested
- **THEN** the report SHALL name the schema that failed and list the schemas
  that do exist, which the CLI supplies with the failure

#### Scenario: The command does not return
- **WHEN** the CLI has not answered within the bounded wait
- **THEN** the wait SHALL be abandoned, the interface SHALL remain responsive,
  and the schema SHALL be reported as having timed out

#### Scenario: The definition cannot be parsed
- **WHEN** a located `schema.yaml` cannot be read as YAML
- **THEN** the schema SHALL be reported as unreadable, naming the path

### Requirement: Schema details are read once per project
Locating a schema costs a subprocess of roughly a second. The application SHALL
do it at most once per schema per project opened, SHALL keep the result for as
long as that project is open, and SHALL request several schemas concurrently so
that the wait does not grow with their number.

#### Scenario: Returning to the tab
- **GIVEN** schema details already read for the open project
- **WHEN** the user leaves the properties tab and returns to it
- **THEN** the details SHALL be shown from what was already read, with no
  further subprocess

#### Scenario: Opening a different project
- **WHEN** a different project is opened, whether at startup or from the picker
- **THEN** the schema details SHALL be read again for that project

#### Scenario: Several schemas at once
- **GIVEN** a project using more than one schema
- **WHEN** their details are read
- **THEN** they SHALL be requested concurrently, so the wait is that of the
  slowest rather than the sum

#### Scenario: A schema file edited while the project is open
- **WHEN** a schema definition changes on disk while the project is open
- **THEN** the details already read SHALL NOT be reread. Schemas are workflow
  definitions rather than content, and the watcher exists for content that moves
  while the user works

#### Scenario: Nothing is read until the tab is opened
- **WHEN** a project is opened and the properties tab is never selected
- **THEN** no subprocess SHALL be run
