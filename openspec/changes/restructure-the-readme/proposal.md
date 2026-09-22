## Why

The README is being read by people deciding whether to try specgetty, and it is
not built for that. Measured:

```
  README.md   535 lines, 3729 words, 7 GIFs totalling 2.5 MB

  what a newcomer needs      51 lines    10%   #
  reference manual          478 lines    89%   #########################
```

The content is good. That is worth saying plainly, because it decides the shape
of this change: almost nothing needs writing, and almost everything needs
moving. What is left is a landing page that never had room to exist.

Three failures cost the most, in order.

**OpenSpec is never explained or linked.** The word appears eleven times in the
README and is not once defined or pointed at. A reader who does not already know
lands on "a text-mode UI tool for reviewing OpenSpec changes and specifications"
and has nowhere to go. Every feature, every screenshot and the whole value
proposition sit behind that one unexplained noun.

**The install section documents the path fewest people want.** Releases are
tagged, `goreleaser` builds `spg` archives for each one, and a working nix flake
is in the repo. The README offers only `go install github.com/mipmip/specgetty@master`:
no binary, no nix, no version, under a heading that says "Source-mode
installation" while offering no other mode. That install path also names
`mipmip` where the repository is `speclib`, which resolves through a redirect and
works, but tells the reader two different things about where the project lives.

**Length and weight.** 535 lines with no table of contents, and 2.5 MB of GIFs
that all load on the landing page. Splitting fixes both: six of the seven GIFs
travel with the sections they illustrate.

Against the consensus of the READMEs on
[awesome-readme](https://github.com/matiassingers/awesome-readme), and the
checklist in "Top ten reasons why I won't use your open source project" that it
links, what is missing is: badges, a table of contents, any mention of the
licence that is already in the repository, a contributing section, and a line
saying where to get help. What is present and strong, and unusual for a project
this size, is seven real screen recordings.

## What Changes

- **The README says why before what.** Writing software has moved from typing
  code to reviewing what a machine proposes, and the bottleneck moved with it:
  there are more specs and more changes to read than before, and they are read in
  editors built for code. That is the opening, and it is where OpenSpec gets
  explained and linked.

- **The README becomes a landing page of about 120 lines**: the why, the feature
  list, installation, a quick start, one screen recording, a table of links to
  the documentation, and then contributing, licence and related projects.

- **Five pages under `docs/` take the 478 lines of reference material**, each
  with the recordings that belong to it:

| page                   | holds                                              | lines |
| ---------------------- | -------------------------------------------------- | ----- |
| `docs/keys.md`         | every key table, at every level                     |   151 |
| `docs/reading-specs.md` | the three levels, the vocabulary, the report view  |   106 |
| `docs/projects.md`     | the picker, stores, the cache, the properties tab   |    99 |
| `docs/change-list.md`  | search, columns, grouping, export                   |    80 |
| `docs/editor.md`       | `E`, `$VISUAL` and `$EDITOR`                        |    42 |

- **Installation covers the three paths that work**: a released binary, `go
  install` at a tagged version, and the nix flake.

- **All five badges, in one pass**: licence, Go version, build status, coverage,
  and the four OpenSpec metrics. `add-continuous-integration` publishes the last
  three and is sequenced before this change, so the README gains its badges once
  rather than twice.

- **Stale passages are corrected while their sections move.** The change list
  still documents an `f` key that `group-the-change-list` removed along with the
  `change_mode` setting, and the specs tab still says a spec that does not parse
  reports on the status line, which `follow-the-spec-grammar` replaced with the
  report view documented forty lines further down the same file.

Not in scope, with the reason for each:

- **Everything that publishes a badge.** The workflows, the `gh-pages` branch and
  the coverage number belong to `add-continuous-integration`, which this change
  is sequenced after. That change has behaviour and a spec delta; this one moves
  documentation and has neither, and drawing the line there is what lets each be
  honest about its spec impact.
- **A documentation site.** `docs/` in the repository is the answer for now.
- **`demo/`.** It stays what it is, the fixture harbour the recordings are made
  against.

## Capabilities

None. This change moves and rewrites documentation, so `.openspec.yaml` sets
`skip_specs: true` rather than inventing a requirement to justify it.

That was not true of the first draft, which also added a badge workflow. This
repository specs its automation: `release-process` carries a requirement about
how the release workflow's two jobs divide the work. A workflow that runs on
push and publishes to a branch has behaviour by that same standard, so it moved
to `add-continuous-integration` where it can be specified, rather than riding
along under a `skip_specs` marker that would have been wrong for it.

## Impact

- `README.md`: cut from 535 lines to about 120
- `docs/`: new, five pages
- `CONTRIBUTING.md`: new, short
- every internal link and image path in the moved sections

## Rollback

Its own commit, and a revert restores a single long README. Nothing outside the
documentation moves, so nothing else can break with it.
