## Why

The tasks pane is the one surface in the application where the cursor selects
something rather than scrolls something, and it is the one surface that was
never taught what it is selecting. It selects a source line. A task is not a
source line. In this repository's own `openspec/changes/add-continuous-integration/tasks.md`,
task 1.1 is two lines and task 1.2 is three:

```
- [ ] 1.1 Add `.github/workflows/check.yml` running on push and on pull request,
      installing nix and running `nix flake check` and nothing else
- [ ] 1.2 Verify the workflow fails when the gate fails, by pushing a branch with
      a deliberately failing test and reading the run, rather than assuming a
      non-zero exit is surfaced
```

Three faults follow from that one gap, and they are the same fault:

- the green band covers the first line of a task and stops, so a task reads as
  half selected
- `j` and `k` step through the remainder of a task, and through the headings and
  the blank lines between tasks, so walking a list of eight tasks takes twenty
  presses
- the page keys do not move the cursor at all, so after `ctrl+f` the cursor can
  be off screen and `space` then toggles a task the reader cannot see

`src/ui/lines.go` already holds the answer and says so in its own words: three
lists share one shape, a cursor over items whose rows are not all one row tall,
and the arithmetic for it is written once. The change list, the spec outline and
the properties list all use it. The tasks pane is the fourth instance and the
only one not wired to it.

There is a second, unrelated fault in the same pane, and it is felt on the same
keystroke. `m.scanning` means two different things. At startup the screen is
empty and the filesystem is being searched, so a modal saying so is honest. On a
refresh the screen is full and one known project is being re-read, and the same
modal blinks over everything (`src/ui/ui.go:1660`). Worse, every key except `q`
and `ctrl+c` is dropped while the flag is set (`src/ui/ui.go:356`). Toggling a
task writes the file, the watcher notices 200ms later, the flag goes up, and
presses made in that window are lost. Ticking boxes quickly is exactly when this
happens and exactly when it is least forgivable.

## What Changes

- **A task is the unit the tasks pane addresses.** A task item is a `- [ ] ` or
  `- [x] ` line at column zero plus the indented lines that follow it, ending at
  the first line that is unindented, blank, a heading, or the next checkbox. The
  cursor selects an item, the highlight covers every row of it, and the page
  keys move by items the way the other three lists already do.

- **The cursor visits tasks only.** Headings and blank lines become chrome,
  which is what `lineOwner = -1` in `src/ui/lines.go` already models for the
  other three lists. `j space j space` becomes the rhythm of working a list.

- **A refresh of a project already on screen is silent.** The scanning flag is
  split in two. A startup scan keeps the modal and keeps refusing keys, because
  there is nothing on screen to interact with. A refresh gets neither, because
  there is. Nothing is drawn for it: the checkbox changing shape is the
  feedback, which is the argument `src/ui/checkbox.go` already makes about the
  status line.

- **The nav bar says what the pane does.** It currently advertises `jk/↑↓
  scroll` on a pane where `jk` moves a cursor, and never mentions `space` at
  all. Both are corrected, gated on the pane having a cursor, by the
  rule the nav bar already follows for `E` and for `←→ diff/old/new`.

## Non-goals

- **Filtering the list to open tasks only.** Considered and dropped. It collides
  with the toggle (the task you tick vanishes under the cursor), it raises a
  question about headings left with nothing under them, and it is not wanted.

- **Widening the boundary rule.** Continuation by indentation covers what
  OpenSpec's own generators write. A task whose continuation is a nested list or
  a fenced block can be handled when one turns up.

- **Optimistic redraw of the toggled checkbox.** Flipping the glyph in the model
  before the rescan lands would cut the last few milliseconds of latency. Once
  the refresh is silent there is nothing left to hide, and it buys a second copy
  of the truth about what the file says.

## Capabilities

### Modified Capabilities
- `task-checkboxes`: the cursor selects a task item rather than a source line,
  and the nav bar advertises the toggle
- `document-viewer`: the paging keys move a cursored document's cursor rather
  than its rows
- `fs-watch`: a refresh of content already on screen interrupts nothing

## Impact

- `src/ui/docview.go`: the cursor and the highlight address items
- `src/ui/checkbox.go`: the item boundary rule
- `src/ui/lines.go`: the fourth caller, no change expected to the file itself
- `src/ui/listpaging.go`: a `docPage()` beside `listPage()`
- `src/ui/ui.go`: the paging and `gg`/`G` branches, the split scanning flag, the
  nav bar hints

The scanning flag has seven uses and they are all in `src/ui/ui.go`. Nothing
outside the tasks pane changes behaviour: a document without a cursor scrolls by
rows exactly as it does today.

## Rollback

The two halves are independent and can be reverted separately. The flag split is
self-contained; reverting it restores the modal and the key guard on refreshes.
Reverting the item unit returns the cursor to source lines, which is a smaller
edit than making it, because the arithmetic it delegates to stays where it is.
