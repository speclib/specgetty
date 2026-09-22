# departure-board Specification

## Purpose
Says when the next ferry leaves, and says so from the far end of the quay.

## Requirements

### Requirement: The next departure is the largest thing on the board
The board SHALL give the next departure more room than everything else together.

#### Scenario: Read from the quay
- **WHEN** the board is read from the far end of the quay
- **THEN** the next departure SHALL be legible when the rest of the board is not
