## Context

See proposal.md for why. What matters here is that the gate already exists and
this change only arranges for it to run somewhere else.

```
  today                              after
  -----                              -----
  a developer:  ship-change.sh       a developer:  ship-change.sh
                  -> nix flake check                 -> nix flake check
  CI:           nothing on push      CI:           check.yml
                release.yml on tags                  -> nix flake check
                                                   badges.yml on main
                                                     -> metrics, coverage
                                                   release.yml on tags
```

`nix flake check` builds the package, runs `go vet`, runs the suite and runs
`scripts/coverage-gate.sh`, which prints a line ending `TOTAL   90.8%` and fails
when any package falls below its floor.

## Goals / Non-Goals

Goals: one gate, run in both places; a coverage number that came from the run it
describes; badges that are true.

Non-Goals: a matrix across Go versions or operating systems, caching strategy
beyond what the nix installer action offers, a release change, or any coverage
service.

## Decisions

### CI runs nix, not a reproduction of it

The alternative is `actions/setup-go` with `make test` and a coverage step, which
is faster and simpler to read. It is rejected because it creates a second
definition of "good": CI would pass on something `ship-change.sh` would refuse,
or the reverse, and which one is right would be argued at the moment someone is
trying to ship.

The cost is a nix installation per run. The project is nix-first already, with
`flake.nix`, `package.nix` and a gate that assumes it, so this is the cheaper
side of the trade.

### The coverage number is taken, not measured

`scripts/coverage-gate.sh` prints the total it enforced. The badge publishes that
line rather than running `go test -coverprofile` again. Measuring twice invites
two numbers, and the one that matters is the one that decides whether a change
may ship.

It is published as a shields endpoint JSON on `gh-pages`, next to the OpenSpec
SVGs, rather than through a coverage service. That keeps the repository free of
a third-party account for a number it already computes, and puts both kinds of
badge in one place with one mechanism.

### Two workflows, not one job with two purposes

The gate and the badges are separated so that a badge failure cannot mark a
commit as failing. They also have different triggers: the gate answers for every
push and pull request, including from a fork, while the badges describe the main
branch and need write permission to publish.

A fork's pull request cannot be given write permission, which is another reason
the two cannot be one job: a contributor's pull request would fail on the
publishing step through no fault of theirs.

### What a first run has to create

`gh-pages` does not exist. The OpenSpec badge action creates it on its first run.
The coverage endpoint is written to the same branch, so its first run has to
tolerate the branch having just been created, or be ordered after the action that
creates it.

## Risks / Trade-offs

- Every push installs nix. Slower than a Go-only workflow, and the first run
  cannot be cached. Accepted for one gate rather than two.
- A workflow with `contents: write` runs on every push to main. It writes only to
  `gh-pages` and only files under `badges/`, and it is the only thing in this
  repository with that permission besides the release.
- A badge is published from the main branch only, so a reader of a pull request
  sees the main branch's numbers. That is what a README badge means everywhere,
  and the check itself does report per pull request.

## Migration Plan

None for a user. The first push after this lands creates `gh-pages` and starts
reporting. `restructure-the-readme` adds the badges to the README and is
sequenced after this, so that the README is edited once rather than twice.
