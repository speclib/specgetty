# board-legibility Specification

## Purpose
Keeps the signs around the harbour readable by the people who need them rather
than by the people who wrote them.

## Requirements

### Requirement: A sign says one thing
A sign SHALL carry one instruction. A sign carrying three is read as carrying
none.

#### Scenario: A sign with more than one instruction
- **WHEN** a sign would carry more than one instruction
- **THEN** it SHALL be split into one sign per instruction
