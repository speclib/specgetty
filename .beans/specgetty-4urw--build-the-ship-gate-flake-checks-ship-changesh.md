---
# specgetty-4urw
title: 'build the ship gate: flake checks + ship-change.sh'
status: completed
type: task
priority: high
created_at: 2026-09-15T14:28:01Z
updated_at: 2026-09-15T14:41:38Z
---

/mip:ship step 4 cannot run. Three pre-existing gaps found 2026-09-15.

## Findings

1. scripts/ship-change.sh does not exist. Only scripts/release.sh is tracked.
2. nix flake check fails to evaluate: flake.nix:14 does
   'nixosModules.default = import ./module.nix self;' but module.nix has never
   been tracked in this repo (git log --all -- module.nix is empty).
3. flake.nix has no checks output at all, so nix flake check would be a no-op
   even once it evaluates. The coverage thresholds named in the skill
   (>=70% overall, >=80% core) exist nowhere in the repo. release.sh runs no
   tests either.

## Baseline

go test ./... passes clean. Coverage as of today:

| package      | coverage |
| ------------ | -------- |
| src/watcher  | 84.6%    |
| src/scanner  | 49.3%    |
| src (main)   | 19.0%    |
| src/ui       | 18.1%    |
| total        | 26.0%    |

## Tasks

- [x] Fix or remove the module.nix reference in flake.nix
- [x] Add a checks output: build, go vet, go test, coverage threshold
- [x] Write scripts/ship-change.sh (stage, gate, archive, commit as Pim Snel, push)
- [x] Decide the coverage threshold policy (see below)

## Open decision

At the skill's stated 70/80 thresholds the gate blocks every change until
coverage roughly triples. A ratchet (today's number is the floor, it may only
rise) makes the gate useful immediately and turns 70/80 into a goal rather than
a wall.

## Summary of Changes

Removed the nixosModules.default line from flake.nix. It imported module.nix,
which has never been tracked in this repo, so nix flake check failed to
evaluate at all.

Added a checks output with two derivations: build (the package compiles) and
tests (go vet, the full suite, and the coverage ratchet).

Added scripts/coverage-gate.sh as the single source of truth for coverage
floors, so nix flake check and a local run enforce the same thing. Floors are
todays measured values per package, and may only be raised.

Added scripts/ship-change.sh: refuses a non-main branch, refuses a change with
unchecked tasks, validates strict, stages before gating (nix only sees tracked
files), runs nix flake check, archives, commits as Pim Snel, pushes main. No
bypass flag.

### Worth knowing

Coverage is not identical across environments. src/watcher measures 84.6%
locally and 87.2% inside the nix sandbox, and TOTAL 26.0 vs 26.1. Floors track
the lowest observed value or the two disagree. This is documented in the gate
script itself.

The 70%/80% target from the ship skill is NOT what the gate enforces. At 26.0%
overall it would have blocked every change. Raising the floors toward that
target is unfinished work.
