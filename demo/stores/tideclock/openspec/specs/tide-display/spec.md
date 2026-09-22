# tide-display Specification

## Purpose
Shows the state of the tide at a glance, so that a glance is enough. The dial is
read from the quayside in poor light by people carrying things, which is what
decides nearly everything below.

## Requirements

### Requirement: The dial shows how far the tide has run
The clock SHALL show the tide's position between low and high water as an angle
on a dial, rather than as a time. A time has to be subtracted from the current
time before it means anything, and the answer wanted is almost always "how much
water is there now".

#### Scenario: Halfway to high water
- **WHEN** the tide is halfway from low water to high water
- **THEN** the dial SHALL point halfway between the two marks

#### Scenario: At the turn
- **WHEN** the tide is within ten minutes of high or low water
- **THEN** the dial SHALL rest on the mark rather than creeping past it, the
  turn being slower than the dial can usefully show

#### Scenario: The direction is visible
- **WHEN** the dial is read
- **THEN** whether the tide is making or ebbing SHALL be visible without waiting
  to see which way the hand moves

### Requirement: The dial reads in poor light
The quayside is lit by one lamp and read at dusk. The dial SHALL be legible
without a light of its own.

#### Scenario: At dusk
- **WHEN** the dial is read in falling light
- **THEN** the marks and the hand SHALL be distinguishable by shape as well as
  by colour

#### Scenario: In direct sun
- **WHEN** the dial is in direct sunlight
- **THEN** the face SHALL NOT reflect into the reader's eyes at the angle it is
  mounted
