---
# specgetty-zqhs
title: detail level of specs inside changes
status: completed
type: feature
priority: normal
created_at: 2026-09-21T18:31:22Z
updated_at: 2026-09-21T21:49:03Z
---

just like the spec detail view there should be a detail view for specs inside a change.


The left sidepanel should show information about `new` `modified` and `deleted` requirements. or other changes. Full research before implementation

## OpenSpec change

`open-a-change-spec-in-detail`, validated strict.

Settled during exploration:

- The operation is a mark on the requirement row (`+` added, `~` modified,
  `-` removed), not an outline level. A requirement belongs to exactly one
  operation, so it is a property of it, and a level would put four indents in a
  pane four tenths of the panel wide.
- Inside a `MODIFIED` requirement each scenario is marked unchanged, edited or
  added, and an unchanged one is dimmed. Measured over the last fourteen
  archived changes: 53% identical, 25% edited, 22% added, 0% dropped. That
  marking is the feature; the diff is the smaller half.
- The zero is structural. OpenSpec cannot drop a scenario from a modified
  requirement, which is why two `REMOVED` blocks here give that as their reason.
- The comparison is offered for active changes only. 25 of 25 `MODIFIED`
  requirements found their original in the live spec. For an archived change 23
  of 24 would report nothing changed and one would attribute a later change to
  this one.
- A node with an original gets a `diff` / `old` / `new` row above the card,
  moved with `left` and `right`, opening on `diff`. Each node opens on the diff
  rather than keeping the last choice.
- The parser gets a delta entry point rather than a flag: the main-spec rules
  invert almost completely in a change.

## What `follow-the-spec-grammar` removed from this

Its `specPart{kind, keyword, text}` keeps unrecognised lines as prose, so a
`REMOVED` requirement's `**Reason**:` and `**Migration**:` already render in
full. An earlier sketch added a fourth card kind for them. Not needed.

## Summary of Changes

Shipped as `open-a-change-spec-in-detail` (commit `e795dae`).

`enter` on a change's specs sub-tab opens its deltas as an outline beside a
card, at a fourth level that escapes back to the sub-tab.

- **The operation is a mark, not a level.** `+` added, `~` modified, `-`
  removed, drawn on the requirement's own row in a gutter every row of a
  change outline carries. A main spec's outline is unchanged.
- **Scenarios inside a modified requirement are marked too**: added, edited, or
  unchanged and drawn back. This is the useful half. Compared on parsed parts
  rather than source lines, so a rewrap is not reported as an edit.
- **Three views where there is something to compare**: the difference, the
  original, the new text, chosen with `left`/`right`, opening on the difference.
  The chooser appears only on a node with an original, so its absence is the
  answer for the rest. Each node opens on its own difference.
- **Active changes only.** An archived change's deltas have been applied to
  the project's specs, so those specs are the result rather than the original.
- **A delta entry point** rather than a flag on `parseSpec`: the main-spec rules
  invert almost completely in a change. Both share the fence mask, the section
  reader and `partsOf`.
- **`E` opens the delta file** the node belongs to, the specs sub-tab one level
  up still offering no key.

## Notes

- The word diff runs per clause, so the keyword stays on its own row in all
  three views and only the marks move between them.
- `TestEveryDeltaInThisRepositoryIsRead` walks every delta here: 101 files, 292
  requirements, 834 scenarios, all parsing.
- Coverage: src/ui 88.7% (floor 88.3), total 89.9%.
