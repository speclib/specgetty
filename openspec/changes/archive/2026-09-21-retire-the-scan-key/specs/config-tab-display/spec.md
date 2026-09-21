## ADDED Requirements

### Requirement: The store's git state is read on entering the properties tab
A store's git state is read from its `.git/` directory, which lies outside every
watched `openspec/` tree, so a commit, fetch or checkout made elsewhere changes
it without any watched file changing. The system SHALL read it when the
properties tab is entered, rather than only when the project is scanned.

This is the rule `schema-inspection` already states for schema definitions:
nothing is read until the tab is asked for, and a session that never opens the
tab never pays for it.

#### Scenario: Entering the tab after a commit elsewhere
- **GIVEN** a project resolved through a store, whose properties tab has been
  viewed
- **WHEN** a commit is made in the store from another terminal and the user
  enters the properties tab again
- **THEN** the reported git state SHALL reflect that commit

#### Scenario: A session that never opens the tab
- **WHEN** a project is open and the properties tab is never entered
- **THEN** no git state SHALL be read

#### Scenario: Still local only
- **WHEN** git state is read on entering the tab
- **THEN** it SHALL be read from local refs alone, with no fetch and no network
  access, exactly as it is when read during a scan

#### Scenario: A project with no store
- **WHEN** the properties tab is entered for a project holding its own content
- **THEN** no git state SHALL be read, there being no store to report on
