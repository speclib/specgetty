## ADDED Requirements

### Requirement: The prompt names the matchers it accepts
The search prompt SHALL name the three matchers while it is focused and no query
has been typed. That is the moment between asking to search and knowing what to
type, and it is the only moment the naming costs nothing: the first keystroke
replaces it with the query.

#### Scenario: The prompt is opened
- **WHEN** the search prompt is focused and the query is empty
- **THEN** the prompt SHALL name the default matcher, the literal-name sigil and
  the contents sigil

#### Scenario: The first character is typed
- **WHEN** any character is typed into the prompt
- **THEN** the naming SHALL disappear, so that a query and its explanation never
  share the line

#### Scenario: The prompt is not focused
- **WHEN** a query is applied but the prompt does not hold the keyboard
- **THEN** the naming SHALL NOT be shown, because nothing is about to be typed

#### Scenario: Too narrow to say it
- **WHEN** the prompt has less room than the naming needs
- **THEN** the naming SHALL be dropped rather than wrapped, leaving the query
  and the result count, which is the same choice the nav bar makes about its
  hints

### Requirement: A failed name search suggests searching the contents
When a query that matched on names alone finds nothing, the application SHALL
suggest the same term as a contents search.

A name search that found nothing is the moment a person is already looking for
another way, which is a better time to be told about one than any other.

#### Scenario: A fuzzy query that matched nothing
- **WHEN** a query with no sigil matches no change
- **THEN** the message SHALL say so and SHALL suggest the same term prefixed
  with the contents sigil

#### Scenario: A literal name query that matched nothing
- **WHEN** a query beginning with `'` matches no change
- **THEN** the same suggestion SHALL be made, because it is also a name search

#### Scenario: A contents query that matched nothing
- **WHEN** a query beginning with `:` matches no change
- **THEN** the message SHALL say so and SHALL suggest nothing, because there is
  no further matcher to try

#### Scenario: The suggestion carries the term
- **WHEN** the suggestion is shown for the query `inotify`
- **THEN** it SHALL name `:inotify` rather than describing the sigil in the
  abstract, so that it can be typed as read
