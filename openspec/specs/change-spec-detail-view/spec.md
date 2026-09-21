# change-spec-detail-view Specification

## Purpose
Describes reading the spec deltas of one change as the structure they have: an
outline marking what each requirement does to its capability, beside a card that
shows the node the reader is on and, where there is something to compare it
with, what the change actually alters.

## Requirements

### Requirement: A change's specs open in a view of their own
Pressing `enter` on a change's specs sub-tab SHALL open that change's spec
deltas in a view below the change view, and `esc` SHALL return to the sub-tab.
This is the rule the specs tab already follows for a live spec.

#### Scenario: Opening the deltas
- **GIVEN** a change is open with its specs sub-tab active
- **WHEN** the user presses `enter`
- **THEN** the change's spec deltas SHALL open in the change spec detail view

#### Scenario: Leaving the view
- **WHEN** the user presses `esc` in the change spec detail view
- **THEN** the change SHALL be shown again with its specs sub-tab active

#### Scenario: A change with no spec deltas
- **WHEN** a change carries no spec deltas, so it has no specs sub-tab
- **THEN** there SHALL be nothing to open

#### Scenario: The view names what it is showing
- **WHEN** the change spec detail view is open
- **THEN** the change's name SHALL be on screen, so the view is never anonymous

### Requirement: The outline covers every delta the change carries
A change may carry several spec files, one per capability it touches. The
outline SHALL list each of them, in a stable order, with that file's
requirements under it and each requirement's scenarios under that.

#### Scenario: A change touching several capabilities
- **WHEN** a change carries deltas for more than one capability
- **THEN** the outline SHALL show each capability as a node, with its own
  requirements beneath it

#### Scenario: The order is stable
- **WHEN** the same change is opened twice
- **THEN** the capabilities SHALL be in the same order both times, and each
  capability's requirements SHALL be in the order its file gives them

#### Scenario: A capability node's card
- **WHEN** the cursor is on a capability node
- **THEN** the card SHALL name the capability and report what the change does to
  it

### Requirement: The operation is a mark on the requirement, not a level
Every requirement in a delta belongs to exactly one operation, so the operation
SHALL be shown as a mark on the requirement's own row rather than as a heading
the reader must scroll to. The outline SHALL NOT gain a level for it.

#### Scenario: An added requirement
- **WHEN** a requirement sits under `## ADDED Requirements`
- **THEN** its outline row SHALL be marked as added

#### Scenario: A modified requirement
- **WHEN** a requirement sits under `## MODIFIED Requirements`
- **THEN** its outline row SHALL be marked as modified

#### Scenario: A removed requirement
- **WHEN** a requirement sits under `## REMOVED Requirements`
- **THEN** its outline row SHALL be marked as removed

#### Scenario: An operation the reader has not seen before
- **WHEN** a delta carries an operation heading that is none of the known ones
- **THEN** the requirements under it SHALL still be listed, marked with what that
  heading said, rather than being hidden

#### Scenario: The marks survive a wrap
- **WHEN** a requirement's title wraps onto more than one row
- **THEN** the mark SHALL be readable on the row the title starts on, and the
  continuation rows SHALL NOT repeat it

### Requirement: A scenario inside a modified requirement says whether it changed
A `MODIFIED` requirement restates every scenario it keeps, and more than half of
what it restates is unchanged. Where the requirement being modified can be
found, each of its scenarios SHALL be marked as unchanged, edited or added, and
an unchanged scenario SHALL be shown as the already-read text it is.

#### Scenario: A scenario the change does not touch
- **GIVEN** a `MODIFIED` requirement whose original can be found
- **WHEN** one of its scenarios is identical to the original's
- **THEN** that scenario SHALL be marked unchanged and drawn less prominently
  than the rest

#### Scenario: A scenario the change edits
- **WHEN** a scenario's text differs from the original's
- **THEN** it SHALL be marked as edited

#### Scenario: A scenario the change introduces
- **WHEN** a scenario has no counterpart in the original requirement
- **THEN** it SHALL be marked as added

#### Scenario: No original to compare with
- **WHEN** the requirement being modified cannot be found
- **THEN** its scenarios SHALL carry no mark, rather than a mark that means
  nothing

### Requirement: The comparison is offered only where it is honest
A comparison SHALL be offered only for a change that has not been archived. For
such a change the specs under `openspec/specs/` are the text the delta modifies,
because nothing has been applied yet. For an archived change they are the result
of applying it, and of every change archived since, so they SHALL NOT be
presented as the original.

#### Scenario: An active change
- **GIVEN** a change that has not been archived
- **WHEN** a `MODIFIED` requirement is shown
- **THEN** the live spec's version of that requirement SHALL be used as the
  original

#### Scenario: An archived change
- **GIVEN** an archived change
- **WHEN** a `MODIFIED` requirement is shown
- **THEN** no comparison SHALL be offered, and the view SHALL say that the text
  the change modified is no longer on disk

#### Scenario: Nothing is inferred from a match
- **GIVEN** an archived change whose delta happens to equal the live spec today
- **WHEN** that requirement is shown
- **THEN** the view SHALL NOT report that the change altered nothing

### Requirement: A node with something to compare offers old, new and the difference
Where a node has an original, the card SHALL offer three views of it: the
difference, the original, and the text the change proposes. The view SHALL open
on the difference. A node with no original SHALL offer no such choice, its
absence being how the view says there is nothing to compare.

#### Scenario: Choosing a view
- **WHEN** a node with an original is shown
- **THEN** the card SHALL offer the difference, the original and the new text,
  and SHALL be showing the difference

#### Scenario: Moving between the views
- **WHEN** the user presses `left` or `right` on such a node
- **THEN** the card SHALL move to the next or previous of the three, stopping at
  either end rather than wrapping

#### Scenario: A node with no original
- **WHEN** an added requirement, a removed requirement or a node of an archived
  change is shown
- **THEN** no choice SHALL be offered, and `left` and `right` SHALL do nothing

#### Scenario: The choice does not follow the cursor
- **WHEN** the user has chosen a view and then moves to another node that also
  has an original
- **THEN** the new node SHALL open on the difference, each node being read from
  its difference first

### Requirement: The deltas are read by OpenSpec's grammar for a change
A spec file in a change is a delta, and the rules that make a main spec valid
invert in it. A delta header SHALL be the structure rather than a fault, a
`## Purpose` section SHALL be optional, and a requirement carrying no scenario
SHALL be accepted where the operation is a removal.

#### Scenario: A delta header
- **WHEN** a change's spec file carries `## ADDED Requirements` or any other
  delta header
- **THEN** it SHALL be read as the operation for the requirements beneath it,
  and SHALL NOT be reported as a fault

#### Scenario: A removal with no scenarios
- **WHEN** a requirement under `## REMOVED Requirements` carries a reason and a
  migration instead of scenarios
- **THEN** it SHALL be accepted, and its card SHALL show what it carries

#### Scenario: A delta with no Purpose
- **WHEN** a change's spec file has no `## Purpose` section, which is every
  delta for a capability that already exists
- **THEN** that SHALL NOT be reported as a fault

#### Scenario: A main-spec heading in a delta
- **WHEN** a change's spec file carries a `## Requirements` section, which
  OpenSpec does not read in a change
- **THEN** it SHALL be reported, because the requirements under it will not be
  applied when the change is archived

#### Scenario: A delta that cannot be read
- **WHEN** a change's spec file does not fit the grammar
- **THEN** the view SHALL report each reason and the line it was found on, the
  way a main spec that does not fit is already reported

### Requirement: The view keeps the behaviour the spec detail view established
The outline and the card SHALL wrap, scroll, page and keep their cursor exactly
as the spec detail view does, that behaviour being described by
`spec-detail-view` and `document-viewer` and not restated here.

#### Scenario: A long label
- **WHEN** a requirement or scenario title is wider than the outline
- **THEN** it SHALL wrap rather than be cut, and the cursor SHALL cover every row
  it produced

#### Scenario: The reading position survives a rescan
- **WHEN** a delta file is rewritten while it is being read
- **THEN** the cursor SHALL stay on the node it was on

#### Scenario: The card is a document
- **WHEN** a card is taller than its pane and holds the keyboard
- **THEN** it SHALL scroll and report its reading position
