## ADDED Requirements

### Requirement: Enter descends whether or not the file fits the grammar
`enter` SHALL descend into the selected spec in every case. Whether the file can
be structured is answered by the view it opens and reported there, not on the
nav bar, because a file usually has several faults and the nav bar holds one
line. The tab SHALL go on showing the whole file as markdown, so the reader who
descends into a report is one `esc` from the content.

#### Scenario: Enter on a spec that does not fit
- **WHEN** the user presses `enter` on a spec whose file cannot be read as
  requirements and scenarios
- **THEN** the spec view SHALL open and report why, rather than the cursor
  staying put and the nav bar carrying the reason

#### Scenario: The nav bar is not used for the reason
- **WHEN** a spec cannot be structured
- **THEN** no transient report SHALL be put on the nav bar, the reasons having a
  place of their own

#### Scenario: The markdown view is unaffected
- **WHEN** a spec cannot be opened as an outline
- **THEN** the specs tab SHALL still render its whole file as markdown, so no
  content becomes unreachable

## REMOVED Requirements

### Requirement: A spec that cannot be parsed reports rather than opening
**Reason**: Replaced by "Enter descends whether or not the file fits the
grammar". A one-line transient report was the right weight for one fault and the
wrong weight for the several a non-conforming file usually has: 133 files in the
local corpus lack a Purpose, 127 lack a `## Requirements` section and 102 carry
a delta header, and most carry all three at once. The reasons need room, a line
number each, and the key that opens an editor, none of which fits on the nav bar.

**Migration**: None for a user. The guarantee that the whole file stays readable
as markdown is carried by the scenario "The markdown view is unaffected" in the
requirement that replaces it.
