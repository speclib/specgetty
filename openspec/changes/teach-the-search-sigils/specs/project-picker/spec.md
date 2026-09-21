## ADDED Requirements

### Requirement: The picker names its matchers and suggests the contents search
The picker's prompt SHALL name the three matchers while it is focused and empty,
and a name search that matched no project SHALL suggest the same term as a
contents search.

The picker's contents search reaches the file paths and the file contents of
every project found on the machine, which is the least guessable capability in
the application.

#### Scenario: The picker's prompt is opened
- **WHEN** the picker's prompt is focused and its query is empty
- **THEN** it SHALL name the default matcher, the literal-name sigil and the
  contents sigil

#### Scenario: A failed name search in the picker
- **WHEN** a query with no sigil, or one beginning with `'`, matches no project
- **THEN** the message SHALL say so and SHALL suggest the same term prefixed
  with the contents sigil

#### Scenario: A failed contents search in the picker
- **WHEN** a query beginning with `:` matches no project
- **THEN** the message SHALL say so and SHALL suggest nothing

#### Scenario: One wording for two meanings
- **WHEN** the naming is shown in either the picker or the change list
- **THEN** it SHALL use the same words, even though the contents sigil reaches
  artifact and spec text in a change and file paths and contents in a project
