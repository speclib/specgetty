## MODIFIED Requirements

### Requirement: A page is what the surface can show
Moving by a page SHALL move by the number of rows the surface is currently
showing, and by half that for the half-page keys, so that the keys mean the same
thing at any terminal size. Where a surface draws items whose rows are not all
one row tall, a page SHALL move the cursor by as many items as fill that number
of rows, so the rule stays a rule about rows and the cursor stays on an item.

#### Scenario: A page in a taller pane
- **GIVEN** two terminal heights
- **WHEN** a page key is pressed in each
- **THEN** the taller one SHALL move further

#### Scenario: Jumping to the ends
- **WHEN** `gg` or `G` is pressed
- **THEN** the surface SHALL go to its first or last row

#### Scenario: Paging past an end
- **WHEN** a page key would move beyond the first or last row
- **THEN** it SHALL stop there rather than wrapping or going out of range

#### Scenario: A page over items of unequal height
- **GIVEN** a list in which some items occupy more rows than others
- **WHEN** a page key is pressed
- **THEN** the cursor SHALL move by as many items as occupy one pane of rows,
  which is fewer items where the items are taller

#### Scenario: Items that are all one row tall
- **GIVEN** a list in which every item occupies one row
- **WHEN** a page key is pressed
- **THEN** the cursor SHALL move by the number of rows the list is showing, which
  is what it does today
