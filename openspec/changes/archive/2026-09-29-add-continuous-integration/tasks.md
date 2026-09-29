## 1. The gate on every push

- [x] 1.1 Add `.github/workflows/check.yml` running on push and on pull request,
      installing nix and running `nix flake check` and nothing else
- [x] 1.2 Verify the workflow fails when the gate fails, by pushing a branch with
      a deliberately failing test and reading the run, rather than assuming a
      non-zero exit is surfaced

      Run 35747526608, branch `ci-verify`: failed, the log naming
      `--- FAIL: TestCIProbeDeliberateFailure` in `src/ui`.
- [x] 1.3 Verify it fails when only the coverage ratchet fails, that being the
      arm most likely to be reported as a pass by a workflow that only watches
      `go test`

      Run 35748030961, same branch: every package reported `ok`, and the run
      still failed at `==> coverage ratchet`, naming
      `src/ui 89.7% < 99.9% FAILED`.
- [x] 1.4 Verify a pull request is checked against the merge result, and that a
      pull request from a fork runs at all

      Run 35749229755 checked out `refs/remotes/pull/7/merge`, its HEAD being the
      merge of the pull request into the base. The fork half is read from the
      workflow rather than exercised, there being no second account to fork from:
      check.yml names no secret and asks for `contents: read`, which is what a
      fork's token already grants.
- [x] 1.5 Record what a run costs in wall-clock time, so the nix decision can be
      revisited with a number rather than an impression

      Ten runs on main: 62s to 146s, most of them near 80s. The cache action that
      made it 3m24s is gone.

## 2. What CI measured is what is published

- [x] 2.1 Take the total from `scripts/coverage-gate.sh` output in the run rather
      than measuring coverage again, and verify the published number matches the
      one the gate printed for that commit

      Run 35844604752 printed `TOTAL 90.8%` from the ratchet and
      `Coverage measured by the gate: 90.8%`; `badges/coverage.json` on gh-pages
      reads 90.8%.
- [x] 2.2 Publish it as a shields endpoint JSON on `gh-pages`, and verify the
      badge renders from the raw URL

      The raw JSON answers 200, and shields renders it from
      `img.shields.io/endpoint?url=...`, also 200.
- [x] 2.3 Publish nothing when the gate failed, and verify the previous number is
      not left describing a commit that did not pass

      Runs 35843944711, 35843838695 and 35842607400 were skipped, their commits
      529c2d3 and bafc146 having failed the gate, and gh-pages carries no commit
      naming either.

## 3. The OpenSpec metrics

- [x] 3.1 Add `.github/workflows/badges.yml` running
      `wearetechnative/openspec-badge-action` with `contents: write`. It runs on
      the Check workflow completing successfully on main rather than on push
      directly, which the delta now says: publishing from a commit the gate
      rejected would describe something untrue of it, and the coverage artifact
      only exists in the run that produced it
- [x] 3.2 Verify the first run creates `gh-pages` and the four SVGs, the branch
      not existing yet

      `65b451c Initial commit for badges` created the branch empty, and `bd414da`
      wrote the four SVGs into it.
- [x] 3.3 Verify the numbers match the repository: 27 specs and 186 requirements
      at the time of writing

      The badges read 28 specs and 192 requirements, which is what the repository
      holds today: the 27 and 186 counted before this change, plus the capability
      it adds and what has landed since.
- [x] 3.4 Verify the coverage endpoint and the SVGs coexist on that branch, the
      two being written by different runs

      All five files are on the branch, the SVGs written by `bd414da` and
      `badges/coverage.json` by `77941f3`.

## 4. Decoration cannot fail the gate

- [x] 4.1 Keep the badges in their own workflow, and verify a failing badge run
      leaves the commit's check green

      Commit 54267b9 carries `gate: success` beside `publish: failure`. The gate
      is untouched by the badge job, which is what the requirement asks; GitHub's
      aggregate state for the commit does count the failed publish job, so the
      badge failure is visible, just not as a verdict on the software.
- [x] 4.2 Verify a fork's pull request is not failed by a publishing step it
      cannot be given permission for

      Read from the workflow rather than exercised: badges.yml runs only on a
      `workflow_run` of Check with `branches: [main]`, which a fork's pull request
      never produces, so the `contents: write` job cannot run against one at all.

## 5. Verification

- [x] 5.1 Verify `nix flake check` still passes locally, unchanged by any of this
- [x] 5.2 Verify the badge markdown the README will use resolves, before
      `restructure-the-readme` puts it there

      Every badge URL in README.md answers 200, the shields endpoint included.
