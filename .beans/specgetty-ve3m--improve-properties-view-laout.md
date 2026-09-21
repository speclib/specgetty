---
# specgetty-ve3m
title: improve properties view laout
status: completed
type: task
priority: normal
created_at: 2026-09-21T21:38:26Z
updated_at: 2026-09-21T22:10:03Z
---

the sidepanel is very small and we need sections with titles in the sidepanel:


Project Info (not selectable)
- Config
- Stores

Schemas (not selectable)
- schema 1


## OpenSpec change

`group-the-properties-list`, validated strict. 25 tasks.

Decisions taken from this bean's exploration:

- Sections `PROJECT` (config, store) and `SCHEMAS` (one row per schema), headers
  uppercase and not selectable, rows indented, a blank line between groups.
- No counts on the headers. `ACTIVE (3)` earns its count on 47 changes; across
  the 33 projects on this machine 25 use one schema and 8 use two.
- The `project` row becomes `config`, since the group header now carries the
  context the old label was doing.
- One store row, singular.
- The `problems` row from `specgetty-p44t` gets a seat under `PROJECT` and is not
  built here.

The width rule had to be replaced rather than worked around. `config-tab-display`
required the list to take only the width its labels need, "so that the content
keeps the room its paths require". Measured, that reason does not hold: the
longest line the content shows is a store root of 60 columns, which wraps at an
80-column terminal and fits at a 100-column one under both the old policy and a
floor of 20. The policy never decides whether a path wraps. The delta removes
that requirement with the measurement as the reason and adds one that keeps its
other four scenarios.

The change also collapses the line arithmetic shared by three lists. The cursor
indexes items, the pane counts drawn lines, and an item may occupy several: the
change list has it, the spec outline has it, and this list would have been the
third copy. The two shipped sites move onto the helper first, where their own
tests prove the move changes nothing, and the new list is built on it.

## Summary of Changes

Shipped as `group-the-properties-list` (commit `91c05e4`).

- **The list is grouped.** `PROJECT` holds `config` and `store`, `SCHEMAS` holds
  one row per schema. Headers are uppercase, not selectable, and a group keeps
  its header when it holds nothing. Rows are indented two, with a blank line
  between groups.
- **`project` is now `config`.** The header it sits under carries the context
  the old label was doing.
- **The list takes a floor of twenty columns**, keeping the one-third cap. The
  divider stays put across a resize and every extra column goes to the content.
- **One line arithmetic for three lists.** `src/ui/lines.go` holds `span`,
  `fit` and `offsetFor` over a slice of item indices, one entry per drawn line.
  The change list and the spec outline moved onto it first with no test edited,
  and the new list was built on it.

## What the work turned up

- **A real defect, caught by task 3.4.** At height 14 the pane opened on the
  blank line between groups. The change list already solved this with a spacer
  skip; the properties list now does the same.
- **Task 5.1 failed to fail, twice.** Reverting the cursor to index drawn lines
  broke nothing, because no test asserted *which* drawn row carried the
  highlight. The first replacement test still passed, matching styled text in
  the content half beside the list. It now asserts the exact string the
  renderer emits for a selected row. Without that, a renderer highlighting the
  wrong row passed every other test in the file.
- **Task 5.2 proved the extraction.** Breaking `offsetFor` alone fails the
  scroll tests of all three lists, which is what says the helper carries the
  arithmetic rather than sitting beside it.
- **Six existing tests asserted the old order or the old label** and were
  updated. One of them, the tab-bar test, was scanning the whole frame for
  `" config "` and started matching the new row label; it now asks
  `renderTabHeader`, which is its actual subject.
- Task 5.3 rendered the tab for all 39 OpenSpec projects on this machine, 378
  rows at three terminal sizes, with no line wider than its pane.

Coverage: src/ui 89.7% (floor 88.3), total 90.7%.
