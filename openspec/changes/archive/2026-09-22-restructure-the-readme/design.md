## Context

See proposal.md for why. What matters here is what the landing page says, in
what order, and what each moved page has to keep working.

The repository already holds more than the README admits: an MIT `LICENSE`, a
`goreleaser` configuration that builds `spg` archives per tag, five tags up to
`v0.7.2`, a nix flake, a `Makefile` with `build`, `test` and `cover`, and seven
screen recordings made against a fixture harbour in `demo/`. Most of this change
is pointing at things that already exist.

## Goals / Non-Goals

Goals: a landing page a stranger can read in a minute; reference material that
is easier to find than it is today, not harder; no claim in the README that the
repository cannot back.

Non-Goals: a documentation site, a rewrite of the reference material, CI or
anything that publishes a badge, or any change to the software.

## Decisions

### The opening says why, and it is not the old why

The README has had two openings and neither does the job. The first asked
whether you had lost track of which OpenSpec projects exist on your machine,
which is a discovery problem and the smaller half of what this is now. The
second states what the tool is and leaves the reader to supply a reason.

The reason is that the work changed. Writing software has moved from typing code
to reviewing what a machine proposes, and the bottleneck moved with it: there
are more specs and more changes to read than there used to be, and they are read
in editors built for writing code rather than for reading a behaviour contract.
specgetty is built for that reading, and it is fast because reading is something
you do many times a day.

That paragraph is also the only natural place to say what OpenSpec is and link
it, which is the single largest gap in the current page.

### What stays on the landing page

In this order, because it is the order a stranger asks the questions in:

```
  title, one line, hero recording      what is this
  badges                               is it alive
  why                                  should I care
  features, nine bullets               what does it do
  install: binary, go, nix             how do I get it
  quick start                          how do I run it
  one more recording                   what does it feel like
  documentation table                  where is everything else
  contributing, help, licence, related what if I want in
```

About 120 lines. Everything below the documentation table is three or four lines
each and exists because the checklist in the linked article names each of them
as a reason people walk away.

### What moves, and how it is cut

Five pages, cut along the seams the README already has rather than by size:

| page                    | why these belong together                      |
| ----------------------- | ---------------------------------------------- |
| `docs/keys.md`          | every key table, so a reader can search one page for a keystroke |
| `docs/reading-specs.md` | the three levels, the vocabulary and the report: one subject, reading a spec |
| `docs/projects.md`      | the picker, stores, the cache and properties: working across projects |
| `docs/change-list.md`   | search, columns, grouping and export: working within one project's changes |
| `docs/editor.md`        | `E` and the environment it reads, which is self-contained and long |

The recordings travel with their sections. Only the hero stays on the landing
page, which takes about 2 MB off it.

### Badges, and why this change is sequenced second

| badge         | source                                  | published by                |
| ------------- | --------------------------------------- | --------------------------- |
| licence       | shields.io, static                      | nothing, it is static       |
| Go version    | shields.io, from `go.mod` (1.25.0)      | nothing, it is static       |
| build status  | the check workflow                      | `add-continuous-integration` |
| code coverage | the number the gate printed             | `add-continuous-integration` |
| OpenSpec      | `wearetechnative/openspec-badge-action` | `add-continuous-integration` |

Three of the five need something to point at, and nothing points anywhere today:
`.github/workflows/` holds only `release.yml`, triggered on tags. A badge with no
run behind it reads as a broken project rather than an absent feature, so this
change waits for the one that creates them and then adds all five at once. The
alternative, shipping the two static badges now and three later, edits the README
twice for the same section.

That ordering is also what keeps this change honestly free of spec impact. A
workflow has behaviour, and this repository specs its automation already.

### Where the OpenSpec numbers come from

The action generates `number_of_specs`, `number_of_requirements`, `tasks_status`
and `open_changes`. For this project the first two read 27 and 186, which is an
unusual thing to be able to put on a landing page and fits a tool whose subject
is specs.

### Link integrity is the one way this change can break something

Every moved section can carry a relative link, an image path, or an anchor that
another section points at. The images are the sharpest case: `demo/recordings/*`
is relative to the README today and becomes `../demo/recordings/*` from inside
`docs/`. A link check is a task rather than an afterthought.

### The stale passages are corrected, not carried

Two sections describe behaviour that no longer exists, and both would otherwise
be moved into a new file and preserved there. The change list documents an `f`
key that cycles active, archived and both; that key, the `change_mode` setting
and the `--change-mode` flag all went with `group-the-change-list`. The specs tab
says a spec that does not follow the structure reports on the status line, which
`follow-the-spec-grammar` replaced with the report view that the same README
documents forty lines further down.

## Risks / Trade-offs

- A reader who knows the current README loses a single searchable page. The
  documentation table and per-page tables of contents are the mitigation, and
  the material is more findable once it is in five named files than it is at
  line 300 of one.
- Five new files is five more things to keep true. They were already five
  sections that had to be kept true.
- This change executes nothing. Everything that runs moved to
  `add-continuous-integration`, which is also what this change now depends on:
  if that one is delayed, the README either ships with three badges pointing at
  nothing or waits.

## Migration Plan

None. No user of the software is affected; only readers of the repository.
