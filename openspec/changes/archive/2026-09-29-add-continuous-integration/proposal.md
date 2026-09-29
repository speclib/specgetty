## Why

Nothing runs on push. `.github/workflows/` holds one file, `release.yml`,
triggered on `v*` tags, so between one release and the next the repository
reports nothing about itself: no run says the tests pass, no number says what
the coverage is, and a pull request from anyone else would be judged by reading
it.

This surfaced while restructuring the README. Four badge types were asked for
and two of them, build status and code coverage, have nothing to point at: a
badge with no run behind it reads as a broken project rather than an absent
feature. The OpenSpec metric badges have the same shape of problem, needing a
workflow and a `gh-pages` branch that does not exist yet.

There is a second reason, and it is the one that decides how CI should work
here. The gate is already defined: `scripts/ship-change.sh` runs `nix flake
check`, which builds the package, runs `go vet`, runs the suite and enforces the
coverage ratchet in `scripts/coverage-gate.sh`. A CI that checked something
slightly different would be a second opinion about whether the tree is good, and
two gates that can disagree are worse than one gate, because the disagreement is
discovered at the worst moment.

## What Changes

- **A check workflow runs the same gate CI and a developer already share.** On
  every push and every pull request, `nix flake check` and nothing else, so that
  what CI enforces and what `ship-change.sh` enforces cannot drift apart.

- **The coverage number CI measured is published**, as a shields endpoint on the
  `gh-pages` branch. `scripts/coverage-gate.sh` already prints it and is already
  the single source of truth the flake reads; this takes the number it printed
  rather than measuring coverage a second way.

- **The OpenSpec metrics are published** by
  [openspec-badge-action](https://github.com/wearetechnative/openspec-badge-action),
  on push to main, to the same branch. For this repository they read 27 specs and
  186 requirements.

- **A failed badge run never fails the check.** Badges are decoration and the
  gate is not, so the two are separate jobs and the decorative one cannot report
  a broken repository because a branch push raced.

## Capabilities

### New Capabilities
- `continuous-integration`: what runs on a push, what it publishes, and what a
  failure of each means

## Impact

- `.github/workflows/check.yml`: new, the gate on push and pull request
- `.github/workflows/badges.yml`: new, the metrics and the coverage number
- the `gh-pages` branch: created by the first badge run
- `README.md`: the badges themselves are added by `restructure-the-readme`,
  which is sequenced after this

## Rollback

Deleting the two workflows returns the repository to reporting nothing, and the
`gh-pages` branch can be deleted with them. No code depends on either.
