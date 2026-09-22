## Why

The dial is read at dusk, and `tide-display` already says it must be legible in
falling light without a light of its own. Winter is the case that rule does not
cover: from November the quay is dark before the last ferry, and shape alone is
not enough when there is nothing to see the shape against.

The harbourmaster's log records eleven occasions last winter when someone came
to the office to ask what the tide was doing. The dial was four metres away.

## What Changes

- The dial is lit from behind, dimly, between dusk and dawn.
- The light is not a feature anyone operates. It is on when the dial cannot be
  read without it, which is a property of the hour and not of a switch.
- A sentence is added to the poor light requirement, and one of its scenarios is
  reworded: `by shape as well as by colour` was written for dusk and reads as an
  instruction not to rely on the light this change adds.

Not in scope: lighting the rest of the quay, which is the harbour's business and
not this dial's.

## Capabilities

### Modified Capabilities
- `tide-display`: the dial reads in the dark as well as in poor light

## Impact

- The dial's housing gains a lamp and a light sensor.

## Rollback

The lamp can be unplugged. The dial reads as it did before, which is to say not
at all after four in the afternoon in January.
