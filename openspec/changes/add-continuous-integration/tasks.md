## 1. The gate on every push

- [ ] 1.1 Add `.github/workflows/check.yml` running on push and on pull request,
      installing nix and running `nix flake check` and nothing else
- [ ] 1.2 Verify the workflow fails when the gate fails, by pushing a branch with
      a deliberately failing test and reading the run, rather than assuming a
      non-zero exit is surfaced
- [ ] 1.3 Verify it fails when only the coverage ratchet fails, that being the
      arm most likely to be reported as a pass by a workflow that only watches
      `go test`
- [ ] 1.4 Verify a pull request is checked against the merge result, and that a
      pull request from a fork runs at all
- [ ] 1.5 Record what a run costs in wall-clock time, so the nix decision can be
      revisited with a number rather than an impression

## 2. What CI measured is what is published

- [ ] 2.1 Take the total from `scripts/coverage-gate.sh` output in the run rather
      than measuring coverage again, and verify the published number matches the
      one the gate printed for that commit
- [ ] 2.2 Publish it as a shields endpoint JSON on `gh-pages`, and verify the
      badge renders from the raw URL
- [ ] 2.3 Publish nothing when the gate failed, and verify the previous number is
      not left describing a commit that did not pass

## 3. The OpenSpec metrics

- [x] 3.1 Add `.github/workflows/badges.yml` running
      `wearetechnative/openspec-badge-action` with `contents: write`. It runs on
      the Check workflow completing successfully on main rather than on push
      directly, which the delta now says: publishing from a commit the gate
      rejected would describe something untrue of it, and the coverage artifact
      only exists in the run that produced it
- [ ] 3.2 Verify the first run creates `gh-pages` and the four SVGs, the branch
      not existing yet
- [ ] 3.3 Verify the numbers match the repository: 27 specs and 186 requirements
      at the time of writing
- [ ] 3.4 Verify the coverage endpoint and the SVGs coexist on that branch, the
      two being written by different runs

## 4. Decoration cannot fail the gate

- [ ] 4.1 Keep the badges in their own workflow, and verify a failing badge run
      leaves the commit's check green
- [ ] 4.2 Verify a fork's pull request is not failed by a publishing step it
      cannot be given permission for

## 5. Verification

- [ ] 5.1 Verify `nix flake check` still passes locally, unchanged by any of this
- [ ] 5.2 Verify the badge markdown the README will use resolves, before
      `restructure-the-readme` puts it there
