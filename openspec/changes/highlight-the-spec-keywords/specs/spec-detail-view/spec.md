## MODIFIED Requirements

### Requirement: A scenario's card lays its clauses out
A scenario's content is raw text, and a clause is the shape that text usually
takes: one source bullet, `- **WHEN** the user...`. Where the card recognises
that shape it SHALL put the keyword on a row of its own with the clause indented
under it, rather than rendering the bullet as it stands. The title SHALL be bold,
a backticked span SHALL be highlighted, and the prose SHALL be left legible.

A clause keyword SHALL be drawn by its role rather than as one kind of word: a
condition, an assertion and a continuation SHALL be told apart. The keywords that
bind SHALL be drawn wherever they appear in a clause's text, by the same rule the
renderer follows everywhere else.

Recognising the shape is a presentation choice and not a condition of being
shown. Content the card does not recognise SHALL be drawn as prose, wrapped and
legible, in its place in the scenario. A clause written in a shape the card does
not lay out SHALL still have its keywords drawn, so that the structure the layout
could not give it is readable anyway.

#### Scenario: A clause is laid out
- **WHEN** a scenario clause is drawn
- **THEN** its keyword SHALL be on a row of its own and the clause text SHALL be
  indented beneath it

#### Scenario: A clause longer than the card
- **WHEN** a clause is longer than the card can show on one row
- **THEN** it SHALL be wrapped and every row after the first SHALL keep the
  clause indentation

#### Scenario: The keywords
- **WHEN** a clause opens with `GIVEN`, `WHEN`, `THEN` or `AND`
- **THEN** that keyword SHALL be drawn in the style its role carries, and a
  condition, an assertion and a continuation SHALL be distinguishable from one
  another

#### Scenario: Code spans
- **WHEN** a clause contains a span between backticks
- **THEN** that span SHALL be highlighted

#### Scenario: Content in an unrecognised shape
- **WHEN** a scenario's content is not in the shape the card lays out
- **THEN** it SHALL be drawn as wrapped prose rather than omitted

#### Scenario: A binding keyword inside a clause
- **WHEN** a clause's text contains `SHALL`, `SHALL NOT`, `MUST` or `MUST NOT`
- **THEN** that keyword SHALL be drawn distinctly from the text around it

#### Scenario: Keywords in prose the card did not lay out
- **GIVEN** a scenario whose clauses are written in a shape the card does not lay
  out, so its content is drawn as prose
- **WHEN** that prose is drawn
- **THEN** the keywords opening its lines SHALL still be drawn as keywords
