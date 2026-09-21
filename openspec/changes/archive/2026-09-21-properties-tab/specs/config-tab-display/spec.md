## MODIFIED Requirements

### Requirement: Config tab shows file source indicator
The tab SHALL show which of its properties is being displayed, as a vertical
list beside the content rather than as a row above it. The list SHALL always
offer the same sections, whatever shape the project has, so that the tab does
not change its own layout from project to project.

#### Scenario: File indicator
- **WHEN** the properties tab is displayed
- **THEN** a list SHALL name the sections and mark the active one, and the
  content beside it SHALL be that section's

#### Scenario: The sections offered
- **WHEN** the properties tab is displayed for any project
- **THEN** the list SHALL offer `project`, one row for each schema the project
  uses, and `store`

#### Scenario: Sub-tabs for a store-backed project
- **WHEN** the properties tab is displayed for a project resolved through a
  store declaration
- **THEN** the list SHALL be the same as for any other project, because the
  configuration that applies is a single file either way

#### Scenario: The indicator is chrome
- **WHEN** the content of a section is scrolled
- **THEN** the list SHALL remain in place beside it

### Requirement: The config tab shows every configuration the project has
The tab SHALL show the one configuration that applies, which is the resolved
root's. A project that declares a store has exactly one configuration in force,
the store's; the declaring file is reported under `store` rather than shown as
if it were in effect.

#### Scenario: The repo's configuration
- **GIVEN** a project resolved through a store declaration
- **WHEN** the `project` section is active
- **THEN** the store's configuration SHALL be shown, because that is the file
  OpenSpec reads

#### Scenario: The store's configuration
- **GIVEN** the same project
- **WHEN** the `project` section is active
- **THEN** the file's path SHALL be named, so that it is clear the content came
  from the store and not from the directory the user is standing in

#### Scenario: A project with one configuration
- **GIVEN** a project holding its own specs and changes
- **WHEN** the `project` section is active
- **THEN** its own configuration SHALL be shown, with the highlighting it has
  today

#### Scenario: The store opened directly
- **GIVEN** a store opened by path rather than from the picker
- **WHEN** the `project` section is active
- **THEN** the store's configuration SHALL be shown, and the `store` section
  SHALL report no declaring file, because there is no repo that pointed here

### Requirement: The config tab carries the store details
The `store` section SHALL report what is known about where the content comes
from, reading only local files, and SHALL report the declaring file together
with what that file declares in vain.

#### Scenario: What the details report
- **GIVEN** a project whose content comes from a store
- **WHEN** the `store` section is active
- **THEN** it SHALL report the store's id, its root path, the file that declared
  it, and the registry's recorded remote and branch when present

#### Scenario: A project with no store
- **GIVEN** a project holding its own specs and changes
- **WHEN** the `store` section is active
- **THEN** it SHALL report that the content is local, rather than being absent
  or empty

#### Scenario: Declarations that have no effect
- **GIVEN** a declaring file that also carries `schema`, `context`, `rules` or
  `operations`
- **WHEN** the `store` section is active
- **THEN** those keys SHALL be named as having no effect, and the report SHALL
  say that the resolved root's configuration is what applies

#### Scenario: A declaring file with nothing inert in it
- **GIVEN** a declaring file carrying only the store declaration
- **WHEN** the `store` section is active
- **THEN** no such warning SHALL be shown

#### Scenario: The store's local git state
- **WHEN** the `store` section is active and the store's root is a git working
  copy
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
- **WHEN** the `store` section is active
- **THEN** it SHALL report the declared id, the file that declared it, and why
  the resolution failed

### Requirement: Each configuration sub-tab is a document of its own
The content beside the list SHALL be a document viewer, with the wrapping,
scrolling and position reporting that capability describes. Selecting a
different section is moving to a different document, so it SHALL start at the
top, by the same rule the spec list already follows.

#### Scenario: Selecting another sub-tab
- **WHEN** a section has been scrolled and another is selected
- **THEN** the second SHALL be shown from its first row

#### Scenario: Leaving the config tab and returning
- **WHEN** the properties tab has been scrolled, another tab is selected and the
  properties tab is selected again, without changing project
- **THEN** the active section SHALL be shown at the position it was left at

#### Scenario: A different project is selected
- **WHEN** a different project is selected
- **THEN** the properties tab SHALL open on its first section, shown from its
  first row

## ADDED Requirements

### Requirement: The properties tab is a list beside its content
The tab SHALL be drawn as a narrow list of sections and a wide content pane,
each in its own border, following the split the specs tab already uses. The
keyboard SHALL move between the two halves, and the lit border SHALL say which
half holds it.

#### Scenario: Moving between the halves
- **WHEN** the properties tab is active and the user presses the key that moves
  focus
- **THEN** the keyboard SHALL move between the list and the content, and on
  through the log panel when it is open, as it does on the specs tab

#### Scenario: Moving down the list
- **WHEN** the list holds the keyboard
- **THEN** the vertical keys SHALL move the selected section

#### Scenario: Scrolling the content
- **WHEN** the content holds the keyboard
- **THEN** the vertical keys SHALL scroll the document

#### Scenario: Which border is lit
- **WHEN** either half holds the keyboard
- **THEN** that half's border SHALL be lit and the other's dim, and both SHALL
  be dim when the log panel holds the keyboard

#### Scenario: The list is sized to its labels
- **WHEN** the tab is drawn
- **THEN** the list SHALL take the width its labels need rather than a fixed
  share of the panel, so that the content keeps the room its paths require

### Requirement: The properties tab reports a schema per row
Each schema the project uses SHALL have its own row, and its content SHALL
report what that schema is and how much of the project runs on it.

#### Scenario: What a schema row reports
- **WHEN** a schema's row is active and its details have been read
- **THEN** the content SHALL report the schema's source, its path, its ordered
  artifacts with what each generates and requires, its apply rule, and the
  number of changes using it

#### Scenario: The project default
- **WHEN** a schema is the one the project's configuration names
- **THEN** its row SHALL say so

#### Scenario: Changes carrying no recorded schema
- **GIVEN** changes with no `.openspec.yaml`
- **WHEN** the schema rows are shown
- **THEN** their number SHALL be reported, rather than being folded into a
  schema they were never recorded as using

#### Scenario: Details not yet read
- **WHEN** a schema's row is active and its details have not arrived
- **THEN** the content SHALL say that it is being read, and the interface SHALL
  stay responsive

#### Scenario: Details that could not be read
- **WHEN** a schema's details could not be read
- **THEN** the content SHALL report the schema by name and the reason it could
  not be read, and the remaining rows SHALL be unaffected
