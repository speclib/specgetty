---
# specgetty-y662
title: 'README: a landing page, not a manual'
status: completed
type: task
priority: normal
created_at: 2026-09-22T15:25:04Z
updated_at: 2026-09-22T16:00:18Z
parent: specgetty-ab9u
---

The README is 535 lines and 89% of it is reference manual. The content is good; it is in the wrong place.

OpenSpec change `restructure-the-readme`, skip_specs. Cuts the README to about 120 lines, moves 478 lines into five pages under docs/, and adds the five badges.

Three failures it fixes: OpenSpec is named eleven times and never explained or linked; the install section documents only `go install @master` while releases, archives and a nix flake exist; 2.5 MB of GIFs load on the landing page.

Ships after `add-continuous-integration`, which publishes three of the five badges.


## Summary of Changes

Shipped as OpenSpec change `restructure-the-readme`, commit `a6d7a88`.
`skip_specs: true`: documentation moved and rewritten, no behaviour changed.

The README is 130 lines against 534, and loads the hero recording alone, 334 KB
against 2.5 MB. The 478 lines of reference material are five pages under
`docs/`: keys, reading-specs, projects, change-list and editor. Each of the
seven recordings now sits with its subject, exactly once.

## The install section did not work, and had not for some time

This is the find worth keeping. The README said:

```bash
go install github.com/mipmip/specgetty@master
```

That fails twice over. There is no `master` revision, the branch being `main`,
and the module has no package at its root, `main` living in `./src`. Every path
in the new section was run before it was written down:

| command                                      | result                        |
| -------------------------------------------- | ----------------------------- |
| `go install .../specgetty@master`            | unknown revision              |
| `go install .../specgetty@latest`            | no package at the module root |
| `go install .../specgetty/src@latest`        | works, installs it as `src`   |
| `releases/latest/download/spg_linux_amd64..` | 404, archives are per version |
| `gh release download --pattern 'spg_*...'`   | works                         |
| `nix run github:speclib/specgetty`           | works, prints 0.7.2           |

So the README now documents nix first, the release archive by pattern, and
`go install .../src@latest` with the rename it needs. That last one is a wart
worth fixing in the repository rather than the README: moving the program to
`cmd/spg` would make the obvious command work and install the right name.

## What else was checked rather than assumed

- Every local link and anchor, by a checker: 0 broken of either kind, including
  the cross-page anchors the split created.
- Every badge URL fetched: six of six return 200, and the Go badge reads v1.25.0,
  which is what `go.mod` says.
- The heading sets before and after: three differ and all three are deliberate
  renames, so no section was lost.
- The six questions a stranger asks, all answered within the first 92 lines.

## Corrections carried out while moving

Two passages described behaviour that no longer exists and would otherwise have
been moved into a new file and preserved there: the `f` key removed by
`group-the-change-list`, and the claim that a spec which does not parse reports
on the status line, which `follow-the-spec-grammar` replaced with the report
view. Both now say what the software does.

## Also added

`CONTRIBUTING.md`, a licence line, a related-projects list, and an explanation
of what OpenSpec is with a link to it. The word appeared eleven times in the old
README and was never once defined, which was the largest comprehension gap on
the page.
