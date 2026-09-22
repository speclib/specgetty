## MODIFIED Requirements

### Requirement: The dial reads in poor light
The quayside is lit by one lamp and read at dusk. The dial SHALL be legible
without a light of its own. From dusk until dawn the dial SHALL be lit from
behind, dimly, so that the hand reads as a silhouette rather than disappearing
with the light.

#### Scenario: At dusk
- **WHEN** the dial is read in falling light
- **THEN** the marks and the hand SHALL be distinguishable by shape as well as
  by colour

#### Scenario: In direct sun
- **WHEN** the dial is in direct sunlight
- **THEN** the face SHALL NOT reflect into the reader's eyes at the angle it is
  mounted, and the lamp SHALL be off

#### Scenario: In the dark
- **WHEN** the dial is read between dusk and dawn
- **THEN** it SHALL be lit from behind, and the hand SHALL read as a silhouette
  against the face

## ADDED Requirements

### Requirement: The light is not operated
Nobody SHALL be required to switch the dial's light on or off. It is on when the
dial cannot be read without it.

#### Scenario: Dusk falls
- **WHEN** the ambient light falls below the threshold
- **THEN** the lamp SHALL come on without anyone acting

#### Scenario: A passing cloud
- **WHEN** the ambient light crosses the threshold briefly
- **THEN** the lamp SHALL NOT flicker on and off

## REMOVED Requirements

### Requirement: The dial carries a torch bracket
**Reason**: The bracket held a torch nobody ever left there, and the lamp
removes the reason it existed.
**Migration**: The bracket can stay on the housing. Nothing reads it.
