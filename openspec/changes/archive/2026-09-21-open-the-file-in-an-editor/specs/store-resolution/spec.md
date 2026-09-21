## MODIFIED Requirements

### Requirement: Every filesystem operation targets the resolved root
Every operation that reads or writes a project's content SHALL act on the
resolved root. None of them SHALL act on the starting directory when the two
differ.

This is stated over every operation rather than over a list of them. A list has
to be joined by each operation added later, and two were not: the task toggle and
the key that copies a change's path both built their paths from the starting
directory, and both were wrong in a store-backed project for as long as that
list stood.

#### Scenario: Reading a store-backed project
- **GIVEN** a project resolved to a store
- **WHEN** its specs and changes are read
- **THEN** they SHALL be read from the store's `openspec/` directory

#### Scenario: Discarding a change in a store-backed project
- **GIVEN** a change open in a project resolved to a store
- **WHEN** the change is discarded
- **THEN** it SHALL be moved within the store's `openspec/changes/`, and no
  directory SHALL be created under the pointing repo

#### Scenario: Exporting a change in a store-backed project
- **GIVEN** a change in a project resolved to a store
- **WHEN** it is exported
- **THEN** the change directory SHALL be read from the store

#### Scenario: Writing a task back in a store-backed project
- **GIVEN** a change open in a project resolved to a store
- **WHEN** a task is toggled
- **THEN** the `tasks.md` under the store SHALL be written, and the toggle SHALL
  succeed

#### Scenario: Reporting a path in a store-backed project
- **GIVEN** a change selected in a project resolved to a store
- **WHEN** a path to it is put on the clipboard or handed to another program
- **THEN** that path SHALL be under the store, and SHALL resolve to something
  that exists

#### Scenario: An operation added later
- **WHEN** an operation that reads or writes project content is added to the
  application
- **THEN** it SHALL act on the resolved root, this requirement covering it
  without being amended to name it
