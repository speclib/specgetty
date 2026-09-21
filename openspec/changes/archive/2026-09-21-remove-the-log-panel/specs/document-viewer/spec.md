## ADDED Requirements

### Requirement: The paging keys act on whichever list or document holds the keyboard
`pgdown`, `pgup`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` SHALL act
on whichever list or document holds the keyboard, by the same rule `j` and `k`
already follow.

One rule rather than one per surface. A list that cannot be paged is walked a
row at a time however long it is, and a document that pages while a list holds
the keyboard moves something the user is not looking at.

#### Scenario: A document that holds the keyboard
- **WHEN** a document holds the keyboard and one of these keys is pressed
- **THEN** the document SHALL scroll or jump

#### Scenario: A list that holds the keyboard
- **WHEN** a list holds the keyboard and one of these keys is pressed
- **THEN** the list's cursor SHALL move, and no document SHALL scroll

#### Scenario: An overlay that holds the keyboard
- **WHEN** the project picker is open and holds the keyboard
- **THEN** its cursor SHALL move, and nothing behind it SHALL

## REMOVED Requirements

### Requirement: The paging keys act on whatever holds the keyboard
**Reason**: One of its scenarios names the log panel as a place the keyboard can
be, and OpenSpec cannot drop a scenario from a modified requirement. Replaced by
the same rule stated over the surfaces that remain, a list or a document.
