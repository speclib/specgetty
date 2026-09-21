<!-- Level-four headings that do not read `Scenario:`. openspec counts these
     as scenarios and says so in a comment: "a level-4 child whose header is
     not literally `Scenario:` (e.g. `#### Edge case`) is still a scenario".
     24 such headings in 4 files of the local corpus, all in mip.rs. -->
# rendering Specification

## Purpose
Describes how a document is turned into what the reader sees.

## Requirements

### Requirement: Markdown is rendered
The renderer SHALL render CommonMark.

#### Parsing
- **WHEN** a document is read
- **THEN** it SHALL be parsed as CommonMark

#### Scenario: Rendering
- **WHEN** a parsed document is drawn
- **THEN** the output SHALL be styled text

#### Edge case ####
- **WHEN** the document is empty
- **THEN** nothing SHALL be drawn
