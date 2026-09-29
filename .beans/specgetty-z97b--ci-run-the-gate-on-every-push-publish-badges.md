---
# specgetty-z97b
title: 'CI: run the gate on every push, publish badges'
status: completed
type: task
priority: normal
created_at: 2026-09-22T15:25:04Z
updated_at: 2026-09-29T11:59:00Z
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


## The verification pass, 2026-09-29

The workflows landed in 54267b9 and the follow-up 0725ebf, but the change kept
fifteen unchecked verification tasks and was never archived. Each one is now
ticked against a run that can be read back, and the change is archived as
`2026-09-29-add-continuous-integration` in commit 30edb10.

What the evidence says today, rechecked rather than copied from above:

| task | evidence                                                            |
|------|---------------------------------------------------------------------|
| 1.2  | run 35747526608 failed, naming TestCIProbeDeliberateFailure         |
| 1.3  | run 35748030961: every package ok, ratchet failed on src/ui         |
| 1.4  | run 35749229755 checked out pull/7/merge                            |
| 1.5  | ten runs on main, 62s to 146s, most near 80s                        |
| 2.1  | gate printed 90.8%, gh-pages coverage.json reads 90.8%              |
| 2.2  | raw JSON 200, shields endpoint 200                                  |
| 2.3  | badge runs skipped for 529c2d3 and bafc146, no gh-pages commit      |
| 3.2  | 65b451c created the branch, bd414da wrote the four SVGs             |
| 3.3  | badges read 28 specs and 192 requirements, which is what is here    |
| 3.4  | SVGs from bd414da, coverage.json from 77941f3, all five coexist     |
| 4.1  | commit 54267b9 carries gate success beside publish failure          |
| 5.1  | nix flake check passes locally, TOTAL 90.8% >= 90.5%                |
| 5.2  | every badge URL in README.md answers 200                            |

Two were read from the workflows rather than exercised, there being no second
account to fork from, and the tasks say so: a fork pull request runs the gate
because check.yml names no secret and asks only for `contents: read`, and it
cannot be failed by publishing because badges.yml fires only on a workflow_run
of Check with `branches: [main]`.

One thing worth knowing about 4.1: the gate stayed green under a failing badge
run, which is what the requirement asks, but GitHub counts the failed publish
job in the aggregate state for the commit. The badge failure is visible there;
it is just not a verdict on the software.

No CHANGELOG entry: the user-facing half of this shipped in 0.7.3, which already
describes continuous integration and the badges.
