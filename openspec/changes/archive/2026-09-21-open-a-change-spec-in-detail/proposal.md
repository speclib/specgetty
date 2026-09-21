## Why

Bean `specgetty-zqhs`. A change's specs sub-tab renders every spec delta of the
change one after another, as markdown, exactly the way a live spec was read
before `open-a-spec-in-detail`. The deferred half of that change is this one.

A delta is harder to read than a main spec, not easier. It restates a whole
requirement to change one sentence of it, and nothing on screen says which part
is the change. Measured over the last fourteen archived changes, inside a
`MODIFIED` requirement:

```
  scenarios identical to the ones being modified   56   53%
  scenarios edited                                 27   25%
  scenarios added                                  23   22%
  scenarios dropped                                 0    0%
```

More than half of what a `MODIFIED` requirement restates is unchanged text the
reader has to recognise as unchanged by reading it. That is the whole problem,
and a badge answers it without rendering a diff at all.

The zero is not a coincidence. OpenSpec cannot drop a scenario from a modified
requirement, which is why two `REMOVED` blocks in this repository give exactly
that as their reason for deleting a requirement rather than amending it.

## What Changes

- `enter` on a change's specs sub-tab opens the change's deltas as an outline
  beside a card, the way `enter` on the specs tab opens a live spec. `esc`
  returns to the sub-tab.

- The outline carries the operation as a mark on each requirement, not as a
  level of its own. A requirement belongs to exactly one operation, so the
  operation is a property of it:

```
╭──────────────────────────────╮ ╭───────────────────────────────────╮
│ specs-tab                    │ │  ┌ diff ┐ old   new               │
│ ~ The specs tab has a focus  │ │  │      └────────────────────────  │
│     Default focus            │ │  │                                 │
│     Moving the focus         │ │  │   Scenario: Focus is visible    │
│   ~ Focus is visible         │ │  │                                 │
│   + Paging from either half  │ │  │   WHEN                          │
│ + Backticked spans are …     │ │  │      the spec list holds the    │
│ change-list-view             │ │  │      keyboard                   │
│ - The old grammar            │ │  │                                 │
╰──────────────────────────────╯ ╰───────────────────────────────────╯
```

  A requirement is `+` added, `~` modified or `-` removed. A scenario inside a
  `~` requirement is unchanged, `~` edited or `+` added. An unchanged scenario
  is dimmed, because "you have read this before" is what the reader needs from
  it.

- Each delta file is a node at the root of the outline, since a change carries
  one spec file at the median and seven at the most in this repository. Its card
  reports the capability and what the change does to it.

- A node with something to compare gains a row above the card: `diff`, `old`,
  `new`, moved with `left` and `right`, opening on `diff`. A node with nothing to
  compare has no row, which is how the view says so.

- The comparison is offered for an **active** change only. For an active change
  the live spec under `openspec/specs/` is the pre-change text by construction,
  and both sides are already in memory. Checked against history over the last
  fourteen archived changes, every one of the 25 `MODIFIED` requirements found
  its original there:

```
  MODIFIED requirements checked         25
    original found to diff against      25
    no original                          0
```

  For an archived change the live spec has already absorbed the delta, so there
  is nothing honest to compare with. Of those same 24 pairs, 23 would report
  nothing changed, which is false, and one would attribute a later change's
  edits to this one. An archived change shows the new text and says why there is
  no comparison.

- The parser gains a delta mode. `follow-the-spec-grammar` made `parseSpec` read
  OpenSpec's grammar for a main spec, where a delta header is an error, a
  `## Requirements` section is required and a requirement without scenarios is a
  fault. In a change every one of those inverts.

Not in scope, and the reason for each:

- **A comparison for an archived change.** Git can reconstruct the original, and
  that is how the numbers above were measured, but it needs the archive commit
  found by heuristic and a subprocess per capability. The active case needs
  neither and is the case people read.

- **Recognising `**Reason**:` and `**Migration**:` as laid-out keywords.**
  `follow-the-spec-grammar` keeps every unrecognised line as prose and argues
  against widening `clauseOf`, so a `REMOVED` card already shows its reason in
  full. Widening it would be presentation, and that change already declined it
  with numbers.

## Capabilities

### New Capabilities
- `change-spec-detail-view`: the view, the operation marks, the comparison and
  the delta grammar it reads

### Modified Capabilities
- `change-list-view`: `enter` on the specs sub-tab descends
- `edit-in-editor`: the surface that had no single file to open now has one

## Impact

- `src/ui/specparse.go`: a delta entry point sharing the fence mask, the section
  reader and `partsOf`, with the validity rules inverted
- `src/ui/changespec.go`: new, the outline marks, the comparison and the card row
- `src/ui/specdetail.go`: the outline and card renderers take a mark per node
- `src/ui/ui.go`: a fourth level, entered from a change and escaping back to it
- `README.md`: the level and its keys

## Rollback

Its own commit. Reverting returns the specs sub-tab to rendering every delta as
one markdown document, which is what ships today.
