# Design

## Context

Two functions read capabilities off the filesystem, `parseChangeDir` for a
change's deltas and `ParseProjectInfo` for a project's main specs, and both do
it with one `os.ReadDir` and a `spec.md` lookup per entry. Everything
downstream receives `SpecNames []string` and `SpecContents map[string]string`.

Those consumers were checked before choosing an approach. They index the slice,
look up the map, build a document key, and build a file path with
`filepath.Join(root, "openspec", "specs", name, "spec.md")`. None of them splits
a name or assumes it has one segment. `pathOf` in `deltaparse.go` joins with
`/` for outline node paths, but those paths are only ever compared for equality
by `indexOfPath`, never parsed.

The one place the two sides meet is `buildChangeSpecTree`, which compares a
delta against `live[cap]`. That makes correctness depend on both walks producing
the same string for the same capability.

## Goals / Non-Goals

**Goals:**

- One discovery rule, in one helper, used by both walks, so a delta and its live
  spec cannot disagree about a capability's name.
- The spec list reads as the tree it is, at the width the pane actually has.
- No new arithmetic for a list that draws lines the cursor cannot land on.

**Non-Goals:**

- Grouping the change delta outline. Its capability nodes already read
  `tenancy/projects`, and two deltas no longer collide once discovery is fixed.
  Nesting them would need a new node kind, an indent rule, a card, and gutter
  handling in the marked view. That is its own change.
- Distinguishing "spec.md missing" from "spec.md unreadable" in the specs tab.
  Both show the same message.

## Decisions

### Discovery is recursive, and shared

A new helper walks a `specs/` directory with `filepath.WalkDir` and returns the
paths of every directory that directly holds a `spec.md`, relative to that
`specs/` directory, `/`-separated. Both call sites use it, so the naming cannot
drift between a change and the project it is read against.

`WalkDir` rather than nested `ReadDir` calls: it is one pass, it does not follow
symlinked directories (neither does the current `ReadDir` plus `IsDir`), and its
error callback is the natural place to skip a subtree that cannot be read.

### The capability rule is the `spec.md`, not the container

A directory is a capability when it directly holds a `spec.md`. Two alternatives
were rejected:

- *A directory is a capability unless it holds subdirectories.* This keeps
  listing a leaf directory that is still empty, and would have kept the existing
  "Spec without spec.md" scenario untouched. It invents a distinction between a
  container and a capability that `openspec list --specs` does not make, and a
  capability that both has a `spec.md` and groups others would fall on the wrong
  side of it.
- *Only the first level is a capability, deeper levels are variants.* This is
  what the code does now, and it is the bug.

The consequence is deliberate: a directory without `spec.md` disappears, which
matches the list `openspec list --specs` prints. In `chatcms-app` the project
goes from 14 empty capabilities to 52 real ones.

### Capability identity stays a flat string

`SpecNames []string` and `SpecContents map[string]string` keep their shape. A
`Spec{Group, Name}` struct was considered and rejected: it touches every one of
the twelve consumers to express something none of them asks for, and the string
is already the identity OpenSpec reports. The path a document opens is
`filepath.Join`, which cleans `a/b` into the platform's separators, so a name
carrying a slash needs no special case.

### Grouping is derived where it is drawn

The scanner returns capabilities. The specs tab derives the groups from the
first path segment at render time. Putting groups in the scanner would push
presentation into it and leave the change delta tree with no equivalent.

### The list reuses `itemLines`

`src/ui/lines.go` already exists because three lists share this shape. The specs
list is the fourth: the cursor indexes capabilities, the pane counts drawn
lines, and headers and blank lines are lines with owner `-1`. `offsetFor` gives
the scroll position and `fit` gives the page size, so `listPage` and
`renderSpecsTab` stop doing their own arithmetic. Hand-rolling the offset is how
a grouped list ends up scrolling one line short of the selected row.

### A row is truncated, not elided

Leaves are short (`reference-instance`), headers shorter. The one long case is a
deep path, and cutting the tail keeps the distinguishing part of a header and
the first part of a deep leaf. `ansi.Truncate` already does this in the detail
pane. A middle elision was considered and rejected as machinery for a case the
grouping mostly removes.

### An unreadable subtree is skipped, not fatal

The helper continues past a directory it cannot read, matching the rule
`openspec-scanning` already states for scan directories: log the error, skip
what cannot be read, report the rest. `ParseProjectInfo` returns no error today,
and giving it one would turn a permission problem in one corner of `specs/` into
an empty project row.

## Risks / Trade-offs

- **A test pins the old count.** `TestParseProjectInfo` asserts that a directory
  without `spec.md` is counted, using `cap-c` as the example. It must be
  rewritten around the new rule, not deleted. The specs-tab scenario "Spec
  without spec.md" is kept, reworded to the unreadable case, because `openspec
  validate` refuses a MODIFIED block that drops a scenario and archive would
  refuse to lose it.
- **Counts jump for real projects.** A nested project's spec count rises from
  the number of containers to the number of capabilities. Intended, but it is
  the visible change in the project picker.
- **A capability named like its own container.** `specs/tenancy/spec.md` beside
  `specs/tenancy/projects/spec.md` sorts the bare `tenancy` row immediately
  before the `tenancy` header, so it reads as part of that group when it is not.
  Accepted: OpenSpec trees are uniformly two segments, and any rule that fixed
  it would make one capability's row appear inside another's group.
- **"No spec.md found" changes meaning.** It now reports a listed capability
  whose file could not be read, rather than a directory that never had one.
- **A very deep tree makes a long list.** Headers add rows, but paging exists and
  moves by what the pane shows.

## Migration Plan

No data changes. Nothing on disk is renamed or moved, and projects with a flat
`specs/` keep rendering as they do today. Reverting the change restores the old
listing without touching any project.

## Open Questions

None.
