## Purpose
Lets a user narrow the change list of one project by typing, matching either the
change names or the text inside a change's artifacts and specs.

## ADDED Requirements

### Requirement: Search input mode
The user SHALL be able to open a search prompt over the change list with `/`,
and the prompt SHALL take keyboard input until dismissed.

#### Scenario: Open the prompt
- **WHEN** the change list is displayed and the user presses `/`
- **THEN** a search prompt SHALL appear and subsequent letter keys SHALL be
  entered as query text rather than triggering actions

#### Scenario: Live filtering
- **WHEN** the user types or deletes a character in the prompt
- **THEN** the change list SHALL immediately re-filter to match the new query

#### Scenario: Navigate while typing
- **WHEN** the search prompt is active and the user presses up, down, `ctrl+p`
  or `ctrl+n`
- **THEN** the list cursor SHALL move and the query text SHALL be unchanged

#### Scenario: Enter opens the highlighted change
- **WHEN** the search prompt is active and the user presses `enter`
- **THEN** the change under the cursor SHALL open at its own level and the filter SHALL
  remain applied to the list behind it

#### Scenario: Escape clears the filter
- **WHEN** the search prompt is active and the user presses `esc`
- **THEN** the query SHALL be cleared, the prompt SHALL close, and the full list
  SHALL be shown

#### Scenario: Prompt visible whenever a filter is active
- **WHEN** a non-empty query is filtering the list
- **THEN** the prompt SHALL be visible with the query in it

### Requirement: Query matching rules
A query SHALL match change names by fuzzy subsequence by default, with prefix
sigils selecting other matchers.

#### Scenario: Default fuzzy match on names
- **WHEN** the query is `expzip`
- **THEN** a change named `export-change-as-zip` SHALL match, because the query
  characters appear in that order in the name

#### Scenario: Literal match on names
- **WHEN** the query begins with `'`
- **THEN** the rest of the query SHALL be matched as a literal substring of the
  change name, and fuzzy matching SHALL NOT be used

#### Scenario: Literal match in artifact text
- **WHEN** the query begins with `:`
- **THEN** the rest of the query SHALL be matched as a literal substring against
  the text of every artifact file and every spec file in each change, and
  against the change name

#### Scenario: Smart case
- **WHEN** the query contains no uppercase characters
- **THEN** matching SHALL be case-insensitive

#### Scenario: Smart case with an uppercase character
- **WHEN** the query contains at least one uppercase character
- **THEN** matching SHALL be case-sensitive

### Requirement: Result ordering
Filtered results SHALL be ordered by match quality so the best match is under
the cursor.

#### Scenario: Ranked results
- **WHEN** a fuzzy query matches more than one change
- **THEN** the changes SHALL be ordered by match score with the strongest match
  first

#### Scenario: Unfiltered ordering
- **WHEN** the query is empty
- **THEN** the list SHALL use its normal ordering

### Requirement: Body matches are explained
When a change matches on artifact or spec text rather than on its name, the row
SHALL indicate which files matched.

#### Scenario: Match in artifact text
- **WHEN** a change matches a `:` query through the contents of proposal.md and
  tasks.md
- **THEN** the row SHALL name those artifacts

#### Scenario: Match on name only
- **WHEN** a change matches through its name
- **THEN** no file names SHALL be shown for that row

### Requirement: Empty result is distinguishable from an empty project
When a query matches nothing, the list SHALL say so and echo the query.

#### Scenario: Query matches nothing
- **WHEN** a non-empty query matches no change in the current filter mode
- **THEN** the list SHALL display a message naming the query, such as
  `No changes match "zipp"`

### Requirement: Filter lifetime
An active filter SHALL persist across navigation within a project and across
rescans, and SHALL be cleared when the selected project changes.

#### Scenario: Descending and returning
- **WHEN** a filter is active, the user opens a change and returns to the change list
- **THEN** the filter SHALL still be applied

#### Scenario: Filesystem rescan
- **WHEN** a filter is active and the filesystem watcher triggers a rescan
- **THEN** the filter SHALL still be applied to the rescanned changes

#### Scenario: Switching project
- **WHEN** a filter is active and the user selects a different project
- **THEN** the filter SHALL be cleared and the new project's full change list
  SHALL be shown

### Requirement: Cursor survives the list changing under it
The selected change SHALL be tracked by name so that filtering and rescanning do
not silently move the selection to a different change.

#### Scenario: Selected change still matches
- **WHEN** the list re-filters and the previously selected change is still
  present
- **THEN** the cursor SHALL remain on that change

#### Scenario: Selected change filtered out
- **WHEN** the list re-filters and the previously selected change is no longer
  present
- **THEN** the cursor SHALL move to the nearest valid row

#### Scenario: Selected change removed by a rescan
- **WHEN** a rescan removes the selected change, for example after it was
  archived or discarded
- **THEN** the cursor SHALL move to the nearest valid row
