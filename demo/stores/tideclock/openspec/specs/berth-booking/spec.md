# berth-booking Specification

## Purpose
Holds a berth for a boat that has said it is coming, and gives it up again when
the boat does not arrive, so that an empty berth is never held by a promise.

## Requirements

### Requirement: A berth is held only as long as it is promised
A booking SHALL hold a berth from the hour it names until two hours after it,
and SHALL release it afterwards without anyone asking.

#### Scenario: The boat arrives
- **WHEN** a boat takes the berth it booked
- **THEN** the booking SHALL be closed and the berth SHALL be shown as occupied

#### Scenario: The boat does not arrive
- **GIVEN** a booking whose hour passed two hours ago
- **WHEN** the berth is still empty
- **THEN** the booking SHALL be released and the berth offered again

### Requirement: A berth is never promised twice
Two boats arriving for one berth is the failure this exists to prevent.

#### Scenario: Two bookings for one berth
- **WHEN** a berth already held is booked again for an overlapping hour
- **THEN** the second booking SHALL be refused, and SHALL say which booking
  holds it
