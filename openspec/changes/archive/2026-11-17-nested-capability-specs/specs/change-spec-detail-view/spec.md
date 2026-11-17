# Spec Delta

## MODIFIED Requirements

### Requirement: The outline covers every delta the change carries
A change may carry several spec files, one per capability it touches. The
outline SHALL list each of them, in a stable order, with that file's
requirements under it and each requirement's scenarios under that.

#### Scenario: A change touching several capabilities
- **WHEN** a change carries deltas for more than one capability
- **THEN** the outline SHALL show each capability as a node, with its own
  requirements beneath it

#### Scenario: Deltas that share a parent segment
- **WHEN** a change carries `specs/tenancy/projects/spec.md` and
  `specs/publish/act/spec.md`
- **THEN** the outline SHALL show `tenancy/projects` and `publish/act` as two
  capability nodes, neither having replaced the other

#### Scenario: The order is stable
- **WHEN** the same change is opened twice
- **THEN** the capabilities SHALL be in the same order both times, and each
  capability's requirements SHALL be in the order its file gives them

#### Scenario: A capability node's card
- **WHEN** the cursor is on a capability node
- **THEN** the card SHALL name the capability and report what the change does to
  it

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

#### Scenario: A delta for a nested capability
- **GIVEN** an active change whose delta is at
  `specs/tenancy/projects/spec.md`
- **WHEN** the project's live spec at
  `openspec/specs/tenancy/projects/spec.md` holds the requirement being
  modified
- **THEN** that requirement SHALL be found as the original, so the delta's
  scenario marks are drawn rather than left off

#### Scenario: An archived change
- **GIVEN** an archived change
- **WHEN** a `MODIFIED` requirement is shown
- **THEN** no comparison SHALL be offered, and the view SHALL say that the text
  the change modified is no longer on disk

#### Scenario: Nothing is inferred from a match
- **GIVEN** an archived change whose delta happens to equal the live spec today
- **WHEN** that requirement is shown
- **THEN** the view SHALL NOT report that the change altered nothing
