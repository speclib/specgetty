## ADDED Requirements

### Requirement: The spec list can be paged and jumped through
The spec list SHALL move its cursor by a page, by half a page, and to either
end, while it holds the keyboard.

#### Scenario: Paging the spec list
- **WHEN** the spec list holds the keyboard and a page or half-page key is
  pressed
- **THEN** the selected spec SHALL move, and the spec content SHALL NOT scroll

#### Scenario: Jumping the spec list
- **WHEN** the spec list holds the keyboard and `gg` or `G` is pressed
- **THEN** the selection SHALL move to the first or last spec

#### Scenario: The content still pages when it holds the keyboard
- **WHEN** the spec content holds the keyboard and a page key is pressed
- **THEN** the content SHALL scroll, as it always has

## REMOVED Requirements

### Requirement: The specs tab pages from either half
**Reason**: Added one commit ago and wrong. With the spec list pageable, a key
pressed on the list half has two things it could move, and this requirement said
to move the one the keyboard is not on. Replaced by the rule that the keys
follow the keyboard, which is what `j` and `k` in this capability already do.
