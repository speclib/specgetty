## MODIFIED Requirements

### Requirement: Config tab shows file source indicator
The config tab SHALL show which configuration is being displayed. Where there is
more than one, the indicator SHALL be a row of sub-tabs, one per configuration,
with the active one marked.

#### Scenario: File indicator
- **WHEN** the config tab is displayed for a project with one configuration
- **THEN** a dimmed line at the top SHALL indicate the file path (e.g.
  "openspec/project.md" or "openspec/config.yaml")

#### Scenario: Sub-tabs for a store-backed project
- **WHEN** the config tab is displayed for a project resolved through a store
  declaration
- **THEN** the top of the tab SHALL show sub-tabs naming the repo's
  configuration, the store's configuration and the store's details, with the
  active one marked

#### Scenario: The indicator is chrome
- **WHEN** the config tab content is scrolled
- **THEN** the indicator, whether a single line or a row of sub-tabs, SHALL
  remain visible above the content border

## ADDED Requirements

### Requirement: The config tab shows every configuration the project has
A project resolved through a store declaration reads from two configuration
files: the originating repo's, which carries the declaration together with that
repo's own context and rules, and the store's, which carries what everything
using the store shares. The config tab SHALL make both reachable, and SHALL NOT
hide either.

#### Scenario: The repo's configuration
- **GIVEN** a project resolved through a store declaration
- **WHEN** the repo sub-tab is active
- **THEN** the originating repo's `openspec/config.yaml` SHALL be shown, with
  the YAML highlighting an ordinary configuration gets

#### Scenario: The store's configuration
- **GIVEN** the same project
- **WHEN** the store sub-tab is active
- **THEN** the store's `openspec/config.yaml` SHALL be shown

#### Scenario: A project with one configuration
- **GIVEN** a project holding its own specs and changes
- **WHEN** the config tab is opened
- **THEN** no sub-tabs SHALL be shown and the tab SHALL read as it does today

#### Scenario: The store opened directly
- **GIVEN** a store opened from the picker, with no originating repo
- **WHEN** the config tab is opened
- **THEN** the store's configuration and its details SHALL be reachable, and no
  repo sub-tab SHALL be offered

### Requirement: The config tab carries the store details
Where a project's content comes from a store, the config tab SHALL carry a
sub-tab reporting what is known about that store. This is where the header's
mark is explained, and it SHALL report only what can be read from local files.

#### Scenario: What the details report
- **GIVEN** a project whose content comes from a store
- **WHEN** the store details sub-tab is active
- **THEN** it SHALL report the store's id, its root path, the originating repo
  when there is one, and the registry's recorded remote and branch when present

#### Scenario: The store's local git state
- **WHEN** the store details sub-tab is active and the store's root is a git
  working copy
- **THEN** it SHALL report whether there are uncommitted changes, and how far
  ahead or behind its upstream ref the working copy is

#### Scenario: Git state read locally only
- **WHEN** the store's git state is reported
- **THEN** it SHALL be read from local refs alone, with no fetch and no network
  access, and the report SHALL say that the comparison is against the last known
  upstream ref

#### Scenario: A store that is not a git working copy
- **WHEN** the store's root is not a git working copy
- **THEN** the details SHALL report the store without git state, rather than
  reporting an error

#### Scenario: A declaration that could not be resolved
- **GIVEN** a project whose store declaration could not be resolved
- **WHEN** the config tab is opened
- **THEN** the store details sub-tab SHALL report the declared id, the file that
  declared it, and why the resolution failed

### Requirement: Each configuration sub-tab is a document of its own
The content behind each sub-tab SHALL be a document viewer, with the wrapping,
scrolling and position reporting that capability describes. Selecting a
different sub-tab is moving to a different document, so it SHALL start at the
top, by the same rule the spec list and a change's artifact sub-tabs already
follow.

#### Scenario: Selecting another sub-tab
- **WHEN** a sub-tab has been scrolled and another is selected
- **THEN** the second SHALL be shown from its first row

#### Scenario: Leaving the config tab and returning
- **WHEN** the config tab has been scrolled, another tab is selected and the
  config tab is selected again, without changing project
- **THEN** the active sub-tab SHALL be shown at the position it was left at,
  which is the retention the single-configuration tab already has

#### Scenario: A different project is selected
- **WHEN** a different project is selected
- **THEN** the config tab SHALL open on its first sub-tab, shown from its first
  row
