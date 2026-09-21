## ADDED Requirements

### Requirement: The paging keys act on whatever holds the keyboard
`pgdown`, `pgup`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` SHALL act
on whichever list, document or panel holds the keyboard, by the same rule `j`
and `k` already follow.

One rule rather than one per surface. A list that cannot be paged is walked a
row at a time however long it is, and a document that pages while a list holds
the keyboard moves something the user is not looking at.

#### Scenario: A document that holds the keyboard
- **WHEN** a document holds the keyboard and one of these keys is pressed
- **THEN** the document SHALL scroll or jump

#### Scenario: A list that holds the keyboard
- **WHEN** a list holds the keyboard and one of these keys is pressed
- **THEN** the list's cursor SHALL move, and no document SHALL scroll

#### Scenario: The log panel
- **WHEN** the log panel holds the keyboard
- **THEN** it SHALL scroll, as it already does

### Requirement: A page is what the surface can show
Moving by a page SHALL move by the number of rows the surface is currently
showing, and by half that for the half-page keys, so that the keys mean the same
thing at any terminal size.

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

## REMOVED Requirements

### Requirement: The paging keys do not wait for focus
**Reason**: Wrong, and corrected one commit later. It made the paging keys
bypass focus and always reach the document, which was a misreading of a report
that turned out to be about lists. With lists pageable there is a list and a
document in play at once on a split tab, and a key that ignores focus has no way
to choose between them. Replaced by "The paging keys act on whatever holds the
keyboard".
