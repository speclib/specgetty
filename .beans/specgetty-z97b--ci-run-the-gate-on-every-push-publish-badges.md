---
# specgetty-z97b
title: 'CI: run the gate on every push, publish badges'
status: completed
type: task
priority: normal
created_at: 2026-09-22T15:25:04Z
updated_at: 2026-09-22T15:51:42Z
parent: specgetty-ab9u
---

Nothing runs on push: .github/workflows holds only release.yml, on v* tags. No run says the tests pass, no number says what coverage is, and a pull request would be judged by reading it.

OpenSpec change `add-continuous-integration`. New capability `continuous-integration`: the gate on every push and pull request, the coverage number CI measured, the OpenSpec metrics, and the rule that decoration cannot fail the gate.

CI runs `nix flake check` and nothing else, the same command scripts/ship-change.sh runs, so that what CI enforces and what shipping enforces cannot drift apart.


## Summary of Changes

Shipped as OpenSpec change `add-continuous-integration`, commit `54267b9`, with
a follow-up fix in `0725ebf`. New capability `continuous-integration`, 4
requirements.

Every push and every pull request now runs `nix flake check` and nothing else:
the same command `scripts/ship-change.sh` runs, so what CI enforces and what
shipping enforces cannot drift apart. The coverage number the gate measured and
the four OpenSpec metrics are published to `gh-pages`.

## What was verified, and how

Before landing, on a throwaway branch, because `workflow_run` fires only for a
workflow already on the default branch:

| what                          | run         | result                        |
| ----------------------------- | ----------- | ----------------------------- |
| a failing test fails the gate | 35747526608 | failed, log names the test    |
| a failing ratchet fails it    | 35748030961 | every test passed, run failed |
| the happy path                | 35748627084 | passed                        |
| a pull request                | 35749229755 | checked out `pull/7/merge`    |

The ratchet run is the one worth keeping: every test passed and the run still
failed, naming the package and its floor. That is the arm a workflow watching
only `go test` would have reported as green.

After landing:

- `gh-pages` created on the first badge run, with the four SVGs.
- The numbers are right and explain themselves: 28 specs and 190 requirements,
  which is the 27 and 186 counted before this change plus the four requirements
  of the capability it added. 1 open change and 0/21 tasks is
  `restructure-the-readme`.
- `badges/coverage.json` reads 90.7%, the number the gate measured in CI, and
  shields renders it: HTTP 200.
- All five files share the branch, written by one job in sequence.

## Two things found by running it rather than reading it

**The cache action was costing time.** `magic-nix-cache-action` returned 400 on
every restore and the run took 3m24s. Removed, the same run takes 1m00s. It was
deprecated as well, but the number is what decided it.

**The badge action failed after succeeding.** It published the SVGs correctly and
then exited 1 on `git checkout main`, because the job had checked out a bare
commit and a detached HEAD has no such branch. Fixed by naming the branch before
the action runs, which keeps the badges describing the commit that passed rather
than whatever main's tip has become. Verified end to end afterwards.

## Owed and paid

The publishing half could not be verified before landing, for the reason above.
Those tasks were reworded to say what was checkable beforehand rather than
checked off as though they had been exercised, the first run on main was watched
immediately, and the one failure was fixed forward the same minute.

## Worth knowing

The ratchet floors are tight against what CI measures: scanner sits at exactly
its floor with no margin, and ui and TOTAL have 0.1. They pass today. The floors
were set from local readings and CI reads slightly lower, so an unrelated change
could turn main red on a rounding difference. `scripts/coverage-gate.sh` forbids
lowering a floor, so this is reported rather than quietly adjusted.
