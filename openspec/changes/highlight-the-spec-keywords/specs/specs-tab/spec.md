## ADDED Requirements

### Requirement: The spec vocabulary is drawn as a vocabulary
A spec is written in a small set of words that carry its meaning, and the
renderer SHALL draw them as such rather than as prose. This SHALL apply wherever
markdown is rendered, so that a spec reads the same on the specs tab, on a
change's spec deltas, and in any other document the same renderer draws.

The words and their roles:

- `GIVEN` and `WHEN` open a condition, `THEN` opens an assertion, and `AND`
  continues the clause above it. These SHALL be recognised at the start of a
  line only, with or without a list marker and with or without bold marks.
- `SHALL`, `SHALL NOT`, `MUST` and `MUST NOT` bind. These SHALL be recognised
  anywhere in prose, because that is where they are written.

A keyword SHALL be recognised only in upper case, and `SHALL NOT` and `MUST NOT`
SHALL each be drawn as one keyword rather than as a keyword followed by a word.

#### Scenario: A clause keyword at the start of a line
- **WHEN** a line begins with a clause keyword, as a bare word, as a list item,
  or wrapped in bold marks
- **THEN** that keyword SHALL be drawn distinctly from the prose after it

#### Scenario: A clause keyword in the middle of a sentence
- **WHEN** a clause keyword appears inside a sentence rather than opening a line
- **THEN** it SHALL be left as prose, being an ordinary word in that position

#### Scenario: A normative keyword inside a sentence
- **WHEN** a requirement reads "the bar SHALL display the battery"
- **THEN** `SHALL` SHALL be drawn distinctly from the words around it

#### Scenario: A two-word keyword
- **WHEN** a line contains `SHALL NOT` or `MUST NOT`
- **THEN** the two words SHALL be drawn as one keyword, and never as a drawn
  keyword followed by an undrawn word

#### Scenario: Lower case is prose
- **WHEN** a line contains the word `shall` in lower case
- **THEN** it SHALL be left as prose, being ordinary English in that case

#### Scenario: A keyword inside a code span
- **GIVEN** a spec naming a keyword rather than using one, as in a sentence about
  `SHALL`
- **WHEN** the line is drawn
- **THEN** the keyword SHALL be drawn as the code span it is part of, and SHALL
  NOT be drawn as a keyword as well

#### Scenario: Roles are told apart
- **WHEN** a document contains a condition, an assertion, a continuation and a
  normative keyword
- **THEN** each SHALL be distinguishable from the others, so that the shape of a
  scenario can be read without reading the words

## MODIFIED Requirements

### Requirement: Backticked spans are highlighted
A span between backticks in a rendered markdown document SHALL be drawn
distinctly from the prose around it, on the specs tab and in every other document
the same renderer draws.

Every segment of a line SHALL be drawn with exactly one style. Styles SHALL NOT
be applied around one another: the end of an inner style ends the outer one, so
a line drawn that way loses its styling from the first nested mark onward.

#### Scenario: A code span in a spec
- **WHEN** a spec contains a span between backticks
- **THEN** that span SHALL be highlighted, and the backticks themselves SHALL NOT
  be shown

#### Scenario: An unclosed backtick
- **WHEN** a line contains a single backtick with no closing one
- **THEN** the line SHALL be rendered unchanged rather than swallowing the rest
  of the document

#### Scenario: Marks inside a code span
- **GIVEN** a line containing a code span that itself contains bold marks, as a
  spec quoting `- **WHEN** ...` does
- **WHEN** the line is drawn
- **THEN** the whole span SHALL be drawn as a code span, and the marks inside it
  SHALL NOT end that styling part way through

#### Scenario: A code span binds tighter than bold
- **WHEN** a line contains both a code span and bold marks
- **THEN** the code span SHALL be read first, as markdown defines, and its
  contents SHALL NOT be read as other marks
