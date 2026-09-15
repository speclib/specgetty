---
# specgetty-g4pa
title: change list should take full view
status: completed
type: feature
priority: normal
created_at: 2026-09-15T10:46:08Z
updated_at: 2026-09-15T15:46:03Z
---

selecting a change and clicking enter should open the change in a deeper level

This view should have filters for showing:
- open/archived/both
- search filter (substring/fzf)

This view should show this fields (configurable)
- title
- total tasks
- open tasks
- schema
- complete (all documents available) (according to schema)
- total specs
- total new specs
- total spec modications (delete/update)
- creation date
- last update date

This view should have the following actions:
- export to zip
- delete (only when all tasts are open)
- archive

## Scope after split (2026-09-15)

Explored in /opsx:explore. Split into two beans. This bean is part A, the
navigation and list work. The column data that needs new scanner capability
moved to [[specgetty-vru8]].

### In scope here

- L1 becomes a full-width change table; L2 is a single change with its own
  artifact sub-nav ([design][proposal][tasks][specs]), which stops left/right
  from spilling out of the change into the parent tab bar (src/ui/ui.go:428)
- changes tab and archive tab merge into one list with an open/archived/both
  filter
- `/` search filter: fuzzy on change names, `:` sigil for literal substring in
  artifact and spec text, score-sorted, smart-case
- Columns limited to what a scan already computes: name, tasks, spec count
- Field selection via config file plus `--change-fields=`; no runtime picker

### Decisions

| Topic         | Decision                                                      |
| ------------- | ------------------------------------------------------------- |
| search key    | `/` opens input mode, live filter                             |
| nav while typing | up/down and ctrl+n/ctrl+p move the cursor                  |
| enter         | opens highlighted change at L2, filter stays applied          |
| esc           | clears filter and exits input mode                            |
| order         | score-sorted, fuzzy.Find                                      |
| lifetime      | survives L1 to L2 and rescans, clears on project switch       |
| cursor        | remembered by change name, restored if it survives, else clamped |
| empty result  | "No changes match `query`", query echoed                      |

### Notes

- sahilm/fuzzy is already in the module graph via bubbles; bubbles/textinput is
  available for the prompt. No new vendoring decision.
- Every consumer of the change list already routes through `currentChanges()` /
  `currentArchivedChanges()`, so the merge and the filter have one insertion point.
- Delete action (from the original body) still needs a key; `d` is taken by
  discard. Overlaps [[specgetty-tyri]].
- [[specgetty-opwv]] (generic cross-project search) left undecided on purpose.
- Related: [[specgetty-g10p]] tab reordering is partly answered by the merge,
  [[specgetty-edee]] and [[specgetty-oarh]] share the nav-level stack.

### OpenSpec change

`openspec/changes/change-list-full-view/` holds proposal.md, design.md and spec
deltas for `change-list-view`, `change-search`, `changes-tab`, `archive-tab`
and `detail-tabs`. Validates strict. tasks.md has 11 groups, 67 items.

## Summary of Changes

Shipped as openspec change `change-list-full-view`, archived to
`openspec/changes/archive/2026-09-15-change-list-full-view/`. Commit 113b89b.

### What landed

Navigation is now an explicit depth (`level`) rather than a `zoomed` boolean:
projects, project, change. `enter` descends, `esc` ascends. A change fills the
panel at its own level, so its artifact sub-tabs own left/right and no longer
spill into the project tab bar.

The changes tab and the archive tab are one list, filtered by `f` between open,
archived and both. It renders as a full-width table whose columns come from
`change_fields` in the config or `--change-fields`.

`/` filters the list: fuzzy on names, `'` for a literal name match, `:` for the
text inside artifact and spec files. Smart case throughout. Results are ranked
strongest first. A row matched on file contents names the files that matched.

### Two things worth remembering

Every consumer of the change list already routed through `currentChanges()`, so
the merge, the mode filter and the search all landed in one function
(`currentRows()`), with nothing downstream needing to know a filter exists.

`fuzzy.Find` compares with `equalFold` and has no case-sensitive mode, so smart
case is enforced by post-filtering its `MatchedIndexes` against the pattern.

### Deferred, deliberately

- Columns needing new scanner data: [[specgetty-vru8]]
- A key for the delete action. `d` is discard; overlaps [[specgetty-tyri]]
- Cross-project search: [[specgetty-opwv]], left undecided
- What sits above the change list: [[specgetty-jdif]]. The specs describe only
  the change list and the change below it, so jdif can redefine the top of the
  stack without contradicting an archived spec. [[specgetty-edee]] and
  [[specgetty-oarh]] describe the model jdif replaces.
