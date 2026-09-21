## ADDED Requirements

### Requirement: The structure read is the structure OpenSpec defines
The view SHALL read a spec by the same rules OpenSpec's own parser uses, so that
a file specgetty can structure is a file the rest of the toolchain can read.
Those rules are: the requirements of a main spec live under a `## Requirements`
heading; a requirement is a `### Requirement: <name>` heading inside that
section; a scenario is any level-four heading with content under it; a heading
inside a fenced code block is not a heading; and a delta header in a main spec
is an error, because it truncates the section OpenSpec parses.

specgetty SHALL NOT add a rule of its own to this set. Where OpenSpec is silent,
as it is about how a scenario's content is written, the view SHALL accept what
it is given rather than require a convention no tool enforces.

#### Scenario: Requirements outside the requirements section
- **GIVEN** a main spec with `### Requirement:` headings that are not inside a
  `## Requirements` section
- **WHEN** the spec is opened
- **THEN** it SHALL NOT be structured, because those requirements are invisible
  to `openspec validate`, `list` and `archive` as well

#### Scenario: A delta header in a main spec
- **GIVEN** a main spec whose requirements sit under `## ADDED Requirements`
- **WHEN** the spec is opened
- **THEN** it SHALL NOT be structured, and the reason SHALL say that a delta
  header belongs in a change

#### Scenario: A scenario heading without the Scenario prefix
- **GIVEN** a requirement containing a level-four heading whose text does not
  begin with `Scenario:`
- **WHEN** the spec is opened
- **THEN** that heading SHALL be a scenario, named by its heading text, because
  that is what OpenSpec counts

#### Scenario: A scenario heading with nothing under it
- **WHEN** a level-four heading has no content before the next heading
- **THEN** it SHALL NOT be a scenario, because OpenSpec does not count it as one

#### Scenario: A heading inside a fenced code block
- **GIVEN** a spec with a fenced code block containing a line that begins with
  one or more `#` characters
- **WHEN** the spec is opened
- **THEN** that line SHALL be content, not a heading, and SHALL NOT end the node
  it sits in

### Requirement: A scenario's content is never dropped
A scenario's content SHALL be carried in full and in the order the file gives
it. No line of it SHALL be discarded for being written in a shape the clause
renderer does not recognise, because the shape of a clause is a convention and
not part of what OpenSpec defines.

#### Scenario: Clauses written without bullets
- **GIVEN** a scenario whose lines read `GIVEN ...`, `WHEN ...` and `THEN ...`
  with no list markers
- **WHEN** its card is drawn
- **THEN** every one of those lines SHALL be on the card

#### Scenario: A scenario written as prose
- **WHEN** a scenario's content is a paragraph carrying no keyword at all
- **THEN** that paragraph SHALL be on the card

#### Scenario: Clauses and prose together
- **GIVEN** a scenario holding both keyword clauses and a paragraph between them
- **WHEN** its card is drawn
- **THEN** both SHALL be present in the order the file gives them

#### Scenario: A separator is not content
- **WHEN** a scenario's content includes a horizontal rule
- **THEN** it SHALL be left out, being a separator in the file rather than part
  of the behaviour described

### Requirement: A file that does not fit the grammar opens as a report
`enter` SHALL descend into the spec view whether or not the file can be
structured. When it cannot, the view SHALL report every reason it found, each
with the line it is on, rather than a single summary. The report SHALL offer the
key that opens the file in an editor.

The reader SHALL NOT be left with one transient line for a file with several
faults, and the file SHALL remain readable as markdown on the specs tab above,
which is where it already is.

#### Scenario: Descending into a file that does not fit
- **WHEN** the user presses `enter` on a spec that cannot be structured
- **THEN** the spec view SHALL open showing why, and `esc` SHALL return to the
  specs tab with the same spec selected

#### Scenario: Every reason is given
- **GIVEN** a file with more than one fault
- **WHEN** its report is drawn
- **THEN** each fault SHALL be named, and none SHALL be omitted in favour of the
  first

#### Scenario: A reason names its line
- **WHEN** a reason concerns a particular line of the file
- **THEN** the report SHALL give that line number

#### Scenario: The report offers the editor
- **WHEN** a report is on screen
- **THEN** the key that opens the file in the user's editor SHALL be offered, so
  the fault can be repaired where it is

#### Scenario: The markdown is still reachable
- **WHEN** a spec cannot be structured
- **THEN** its whole file SHALL still be shown as markdown on the specs tab, so
  no content becomes unreachable

## MODIFIED Requirements

### Requirement: A scenario's card lays its clauses out
A scenario's content is raw text, and a clause is the shape that text usually
takes: one source bullet, `- **WHEN** the user...`. Where the card recognises
that shape it SHALL put the keyword on a row of its own with the clause indented
under it, rather than rendering the bullet as it stands. The title SHALL be bold,
the keywords SHALL be coloured, a backticked span SHALL be highlighted, and the
prose SHALL be left legible.

Recognising the shape is a presentation choice and not a condition of being
shown. Content the card does not recognise SHALL be drawn as prose, wrapped and
legible, in its place in the scenario.

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
- **THEN** that keyword SHALL be drawn in its own colour

#### Scenario: Code spans
- **WHEN** a clause contains a span between backticks
- **THEN** that span SHALL be highlighted

#### Scenario: Content in an unrecognised shape
- **WHEN** a scenario's content is not in the shape the card lays out
- **THEN** it SHALL be drawn as wrapped prose rather than omitted

## REMOVED Requirements

### Requirement: A spec that does not fit the structure does not open
**Reason**: The behaviour it describes is replaced by "A file that does not fit
the grammar opens as a report". The requirement turned `enter` back at the door
and put the reason on the nav bar, which cannot hold the several reasons a file
usually has, nor the line each sits on, nor the key that opens an editor on it.
Its criterion was also specgetty's own reading of the structure rather than the
grammar OpenSpec defines, which is what let 136 of 544 live specs open and show
part of what they contain.

**Migration**: None for a user. Its one surviving guarantee, that the whole file
stays readable as markdown on the specs tab, is carried by the scenario "The
markdown is still reachable" in the requirement that replaces it.
