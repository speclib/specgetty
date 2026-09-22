# continuous-integration Specification

## Purpose
TBD - created by archiving change add-continuous-integration. Update Purpose after archive.

## Requirements

### Requirement: The gate runs on every push and every pull request
The repository SHALL check itself automatically, and SHALL check the same thing
a developer checks before shipping. That check is `nix flake check`, which builds
the package, vets it, runs the suite and enforces the coverage ratchet.

CI SHALL NOT define a check of its own. Two gates that can disagree are worse
than one, because the disagreement is found when a change is already believed to
be finished.

#### Scenario: A push is checked
- **WHEN** a commit is pushed to any branch
- **THEN** the gate SHALL run against it

#### Scenario: A pull request is checked
- **WHEN** a pull request is opened or updated
- **THEN** the gate SHALL run against the merge result, so that a contributor is
  answered without a maintainer reading the diff first

#### Scenario: The gate is the one that already exists
- **WHEN** the check workflow runs
- **THEN** it SHALL run `nix flake check`, the same command
  `scripts/ship-change.sh` runs, rather than a list of steps that reproduce it

#### Scenario: A failing gate is reported as failing
- **WHEN** the gate fails
- **THEN** the run SHALL fail, and the failure SHALL name which of the build, the
  vet, the tests or the coverage ratchet was not satisfied

### Requirement: What CI measured is what is published
The coverage number shown for the repository SHALL be the one the gate measured
on the commit it describes. It SHALL NOT be measured a second way, because a
second measurement can disagree with the one that decides whether a change may
ship.

#### Scenario: The number comes from the gate
- **WHEN** coverage is published
- **THEN** it SHALL be the total `scripts/coverage-gate.sh` printed during the
  run, that script being the single source of truth the flake already reads

#### Scenario: A commit that did not pass publishes nothing
- **WHEN** the gate fails
- **THEN** no coverage number SHALL be published for that commit, rather than a
  stale number being left to describe it

### Requirement: The repository publishes what it is made of
The specs, requirements, task ratio and open changes of the project SHALL be
published as badges, so that a reader can see the shape of the project without
cloning it. These come from the OpenSpec badge action rather than being counted
here, so that the numbers and the tool that produces them cannot drift.

#### Scenario: The metrics are refreshed
- **WHEN** a commit lands on the main branch and the gate passes on it
- **THEN** the specs, requirements, tasks and open changes badges SHALL be
  regenerated from that commit

#### Scenario: A commit the gate rejected describes nothing
- **WHEN** the gate fails on a commit on the main branch
- **THEN** no badge SHALL be regenerated from it, so that nothing published
  describes a commit that did not pass

#### Scenario: Where they are published
- **WHEN** badges are generated
- **THEN** they SHALL be written to the `gh-pages` branch, which SHALL be created
  if it does not exist

### Requirement: Decoration cannot report a broken repository
Publishing badges SHALL be separate from the gate, and a failure to publish SHALL
NOT fail the check. A badge is decoration and the gate is not: a branch push that
raced, or a token that expired, is not evidence that the software is broken.

#### Scenario: Badge publishing fails
- **WHEN** a badge run fails for any reason
- **THEN** the check on that commit SHALL be unaffected, and the commit SHALL NOT
  be reported as failing

#### Scenario: The gate fails
- **WHEN** the gate fails
- **THEN** that SHALL be reported as a failure, badges being no part of that
  judgement
