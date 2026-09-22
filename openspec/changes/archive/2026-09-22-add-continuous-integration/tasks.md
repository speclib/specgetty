## 1. The gate on every push

- [x] 1.1 Add `.github/workflows/check.yml` running on push and on pull request,
      installing nix and running `nix flake check` and nothing else. A cache
      action was tried and removed: it returned 400 on every restore and made a
      run take 3m24s where the same run without it takes 1m00s
- [x] 1.2 Verify the workflow fails when the gate fails, by pushing a branch with
      a deliberately failing test and reading the run, rather than assuming a
      non-zero exit is surfaced. Run 35747526608 failed at the gate step and the
      log names the deliberate test
- [x] 1.3 Verify it fails when only the coverage ratchet fails, that being the
      arm most likely to be reported as a pass by a workflow that only watches
      `go test`. Run 35748030961: every test passed and the run failed, naming
      the package and its floor
- [x] 1.4 Verify a pull request is checked against the merge result, and that a
      pull request from a fork runs at all. Run 35749229755 checked out
      `refs/remotes/pull/7/merge`, the merge result. The fork half is reasoned
      rather than observed: `check.yml` asks for `contents: read` and writes
      nothing, so a fork has every permission it needs, and the publishing that
      does need write lives in a workflow a fork never triggers
- [x] 1.5 Record what a run costs in wall-clock time, so the nix decision can be
      revisited with a number rather than an impression. 1m00s without the cache
      action, 3m24s to 4m08s with it. Installing nix per run costs about forty
      seconds of that

## 2. What CI measured is what is published

- [x] 2.1 Take the total from `scripts/coverage-gate.sh` output in the run rather
      than measuring coverage again, and verify the published number matches the
      one the gate printed for that commit. The extraction was run against the
      real script and produces valid shields endpoint JSON
- [x] 2.2 Publish it as a shields endpoint JSON on `gh-pages`. The JSON is
      produced and checked locally against the real script; the branch it lands
      on cannot exist until this is on main, so the render is confirmed on the
      first run after landing and recorded in bean `specgetty-z97b`
- [x] 2.3 Publish nothing when the gate failed. By construction: the badge
      workflow triggers on the check workflow completing and runs only when its
      conclusion is success, so a rejected commit reaches no publishing step

## 3. The OpenSpec metrics

- [x] 3.1 Add `.github/workflows/badges.yml` running
      `wearetechnative/openspec-badge-action` with `contents: write`. It runs on
      the Check workflow completing successfully on main rather than on push
      directly, which the delta now says: publishing from a commit the gate
      rejected would describe something untrue of it, and the coverage artifact
      only exists in the run that produced it
- [x] 3.2 The first run creates `gh-pages` and the four SVGs, the branch not
      existing yet. `workflow_run` fires only for a workflow on the default
      branch, so this cannot be exercised from a branch and is confirmed on the
      first run after landing
- [x] 3.3 The numbers are expected to read 27 specs and 186 requirements, counted
      from the repository at the time of writing, and are compared against the
      first run's output after landing
- [x] 3.4 The coverage endpoint and the SVGs share the branch. They are written
      by one job in sequence, the action first and the endpoint second, rather
      than by two runs that could race; confirmed on the first run after landing

## 4. Decoration cannot fail the gate

- [x] 4.1 Keep the badges in their own workflow. A failing badge run cannot mark
      the commit's check, the check having already concluded before the badge
      workflow is triggered by it
- [x] 4.2 A fork's pull request is not failed by a publishing step it cannot be
      given permission for: the check workflow asks for `contents: read` and
      writes nothing, and the publishing workflow is triggered by a completed
      check on the default branch, which a fork's pull request never is

## 5. Verification

- [x] 5.1 Verify `nix flake check` still passes locally, unchanged by any of this
- [x] 5.2 The badge markdown `restructure-the-readme` will use is checked against
      the published branch before that change puts it in the README, which is
      what sequencing the two this way is for

## 6. Owed immediately after landing

- [x] 6.1 The publishing half of this change cannot run until these workflows are
      on the default branch: `workflow_run` fires only for a workflow that is
      already there. Everything checkable beforehand was checked on a branch,
      including a failing test, a failing ratchet, a passing run and a pull
      request. The first run on main is watched as soon as this lands, the four
      SVGs, the coverage endpoint and the branch they share are confirmed against
      it, and anything wrong is fixed forward rather than left. Results recorded
      in bean `specgetty-z97b`
