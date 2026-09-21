# config-tab-display Specification

## Purpose
The tab that reports what a project is and where its parts come from: the one
configuration that applies, the workflow schemas its changes use, and the store
its content lives in. On screen the tab is called `properties`.

The capability keeps its original directory name because OpenSpec's RENAMED
operation works on requirements and has no capability-level form, and removing
every requirement to recreate them elsewhere would leave a spec with an empty
requirements section that `validate --strict` rejects. The name is a wart, not
an oversight.

## Requirements

### Requirement: Config tab displays project.md as styled markdown
When an OpenSpec project has a `project.md` file, the config tab SHALL render it with basic markdown styling.

#### Scenario: Project with project.md
- **WHEN** user switches to the config tab for a project that has `openspec/project.md`
- **THEN** the content SHALL be displayed with headers styled bold and colored, list items indented, and bold/italic text styled appropriately

#### Scenario: Markdown headers
- **WHEN** a line starts with `#`, `##`, or `###`
- **THEN** it SHALL be rendered in bold with a distinct color

#### Scenario: Markdown list items
- **WHEN** a line starts with `- ` or `* `
- **THEN** it SHALL be rendered with proper indentation

### Requirement: Config tab displays config.yaml with syntax highlighting
When an OpenSpec project has a `config.yaml` file (and no project.md), the config tab SHALL render it with YAML syntax highlighting.

#### Scenario: Project with config.yaml only
- **WHEN** user switches to the config tab for a project that has `openspec/config.yaml` but no `openspec/project.md`
- **THEN** the YAML content SHALL be displayed with keys in one color, values in another, and comments dimmed

#### Scenario: YAML comments
- **WHEN** a line contains a `#` comment
- **THEN** the comment portion SHALL be rendered in a dimmed style

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

### Requirement: Config tab handles missing configuration
When no configuration file exists, the config tab SHALL show an appropriate message.

#### Scenario: No config file
- **WHEN** a project has neither `openspec/project.md` nor `openspec/config.yaml`
- **THEN** the config tab SHALL display "No project configuration found"

### Requirement: Config tab content is a document viewer
The content of the config tab, whether it is styled markdown or highlighted
YAML, SHALL be a document viewer, with the wrapping, scrolling, position
reporting and position retention that capability describes.

#### Scenario: Configuration longer than the panel
- **WHEN** the config tab shows a `project.md` or `config.yaml` longer than the
  panel
- **THEN** the user SHALL be able to reach the end of it with the keyboard, and
  the panel title SHALL report the reading position

#### Scenario: Leaving and returning to the tab
- **WHEN** the config tab has been scrolled and the user switches to another tab
  and back, without changing project
- **THEN** the content SHALL be shown at the position it was left at

#### Scenario: A different project is selected
- **WHEN** the config tab has been scrolled and the user selects a different
  project
- **THEN** that project's configuration SHALL be shown from its first row

#### Scenario: The file source indicator stays put
- **WHEN** the config tab content is scrolled
- **THEN** the dimmed line naming the file SHALL remain visible at the top of
  the tab

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

### Requirement: The properties tab is a list beside its content
The tab SHALL be drawn as a narrow list of sections and a wide content pane,
each in its own border, following the split the specs tab already uses. The
keyboard SHALL move between the two halves, and the lit border SHALL say which
half holds it.

#### Scenario: Moving between the halves
- **WHEN** the properties tab is active and the user presses the key that moves
  focus
- **THEN** the keyboard SHALL move between the list and the content, as it does
  on the specs tab

#### Scenario: Moving down the list
- **WHEN** the list holds the keyboard
- **THEN** the vertical keys SHALL move the selected section

#### Scenario: Scrolling the content
- **WHEN** the content holds the keyboard
- **THEN** the vertical keys SHALL scroll the document

#### Scenario: Which border is lit
- **WHEN** either half holds the keyboard
- **THEN** that half's border SHALL be lit and the other's dim

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

### Requirement: The properties rows can be paged and jumped through
The properties row list SHALL move its selection by a page, by half a page, and
to either end, while it holds the keyboard.

Its list is short, so the jumps matter more than the pages; both are there
because every list in the application answers to the same keys.

#### Scenario: Jumping the row list
- **WHEN** the row list holds the keyboard and `gg` or `G` is pressed
- **THEN** the selection SHALL move to the first or last row, and the content
  SHALL NOT scroll

#### Scenario: Paging a list shorter than a page
- **WHEN** a page key is pressed and the list is shorter than a page
- **THEN** the selection SHALL move to that end and stop

#### Scenario: The content still pages when it holds the keyboard
- **WHEN** the content holds the keyboard and a page key is pressed
- **THEN** it SHALL scroll

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
