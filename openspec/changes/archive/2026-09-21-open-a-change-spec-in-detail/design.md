## Context

See proposal.md for why. What matters here is that `follow-the-spec-grammar`
landed first and changed what this builds on.

`parseSpec` now returns `(specTree, []specProblem)` and reads OpenSpec's grammar
for a **main** spec: a delta header is an error, `## Requirements` is required,
`## Purpose` is required, and a requirement with no scenario is a fault. Every
one of those inverts in a change.

It also changed a scenario's content model:

```go
type specPart struct {
    kind    int    // partClause or partProse
    keyword string // for partClause
    text    string
}
```

Ordered parts rather than clauses beside a body, and nothing is dropped. That
removes work this change would otherwise have needed: a `REMOVED` requirement's
`**Reason**:` and `**Migration**:` are not clause keywords, so they arrive as
`partProse` and already render in full. An earlier sketch of this change added a
fourth card kind for them. It is not needed.

The view it extends already has: an outline with wrapped labels and a node
cursor, `specDetailSplit`, `nodesInRows` for the page walk, the cursor kept by
node path across a re-parse, and a report state for a file that does not parse.

## Goals / Non-Goals

Goals: read a change's deltas by the grammar OpenSpec applies to a change; say
what each requirement does to its capability; say which parts of a `MODIFIED`
requirement are actually different.

Non-Goals: reconstructing an archived change's original from git. Widening which
clause shapes get laid out. Editing. Validating a change in specgetty's own
right.

## Decisions

### The operation is a mark, not a level

Three shapes were considered: the operation as an outline level, a two-step
navigation with a list of the change's spec files first, and the operation as a
mark on the requirement row.

The mark wins on three counts. A requirement belongs to exactly one operation,
so the operation is a property of it and not a container for it. A level would
put four indents in an outline that is four tenths of the panel, where a
scenario title already has a median width of 28 columns and a maximum of 63. And
a mark reads in place, where a heading has to be scrolled to and remembered.

The two-step shape was rejected because the median change in this repository
carries one spec file. It would add a navigation step whose list has one row.

### The scenario marks are the feature, and they need no diff

```
  scenarios identical    53%
  scenarios edited       25%
  scenarios added        22%
  scenarios dropped       0%
```

Half of what a `MODIFIED` requirement restates is unchanged. Marking those and
dimming them answers the reader's question without rendering anything. The diff
is the smaller half of the value, which is why the marks are specified as their
own requirement and the three-way card row as another: the first can ship and be
useful if the second turns out badly.

The zero is structural rather than lucky. OpenSpec cannot drop a scenario from a
modified requirement, which is why two `REMOVED` blocks in this repository give
that as their reason for deleting a requirement instead of amending it. So the
scenario vocabulary is three marks, not four. Defensive handling of a fourth
costs nothing and is not specified.

### The comparison is refused for an archived change

For an active change the live spec is the pre-change text by construction, and
both sides are already in memory: `info.SpecContents` and `ci.SpecContents`. No
file read, no subprocess.

For an archived change the live spec is the result of applying the delta, and of
every change archived since. Measured over 24 pairs, 23 would report nothing
changed and one would attribute a later change's edits to this one. Both are
lies, and the second is the worse kind because it looks like information.

Git can reconstruct it, and that is how the numbers in the proposal were
measured: `git show <archive-commit>~1:openspec/specs/<cap>/spec.md`. It needs
the archive commit found by heuristic and a subprocess per capability, against a
case nobody reads. Refusing is cheap and honest, and the refusal is specified so
that a match between delta and live spec is never reported as "this change
altered nothing".

### A delta entry point, not a flag on parseSpec

`parseSpec` encodes the main-spec grammar in its control flow, not in a table:
the Purpose search, the delta-header scan, the `## Requirements` bounds and the
outside-the-section check are each a pass over the file. A `isDelta bool`
threaded through all of them would put a branch in every pass and leave one
function stating two grammars.

A separate entry point shares what is genuinely shared, which is the mechanical
half: `fenceMask`, `sectionBody`, `headingLevel`, `scenarioName`, `partsOf`,
`clauseOf`. It states its own validity rules, which are almost the complement of
the other's. The tree it produces is the same `specTree`, with an operation on
each requirement node, so every renderer downstream is unchanged except for the
mark.

### The three-way row takes left and right

The change view one level up already draws a sub-tab row and moves it with
`left` and `right`. The card's row is the same idiom in the same place on
screen, so the keys are the ones already learned.

This is why `keep-the-arrows-in-the-spec-view` had to land first. Those keys
used to fall through to the project tab bar from any level that did not claim
them; they now do nothing at a spec level, which is the state a new claim can be
added to. Claiming them here on top of a fall-through would have been two bugs
interacting.

Each node opens on the difference rather than keeping the last choice. The
choice is about one node, and carrying it to the next one would mean the reader
sees an original without having asked for it.

## Risks / Trade-offs

- A delta whose `MODIFIED` requirement title does not match anything live, from
  a typo or a rename. → Specified: no original means no marks and no row, rather
  than marks that mean nothing.
- Matching a requirement by its heading text is exact, so a reworded heading
  reads as a different requirement. → That is what OpenSpec itself does when it
  applies a delta, so specgetty agreeing with it is correct even when the result
  is unhelpful.
- Two parsers for one file format could drift. → They share every mechanical
  helper and produce the same tree. What differs is the validity rules, which
  are the part that genuinely differs.
- A fourth level in `m.level`. → It is an enum of views with explicit enter and
  esc mappings rather than a depth counter, which `levelSpec` already
  established by being reached from `levelProject`.

## Migration Plan

None. Without pressing `enter` on a change's specs sub-tab nothing behaves
differently.
