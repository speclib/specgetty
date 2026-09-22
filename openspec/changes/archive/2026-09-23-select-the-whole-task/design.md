## Context

Two problems in one pane, with nothing in common except the keystroke that
exposes both. They are proposed together because a reader ticking boxes meets
them in the same second, and because the tests for both live in
`src/ui/checkbox_test.go` and `src/ui/docview_test.go`. They can be shipped
separately if that turns out to be easier.

## The item unit

### What a task item is

```
- [ ] 1.1 Add `.github/workflows/check.yml` running on push and on   ┐
      installing nix and running `nix flake check` and nothing else  ┘ one item
- [ ] 1.2 Verify the workflow fails when the gate fails, by pushing  ┐
      a deliberately failing test and reading the run                ┘ one item
                                                                       blank ends it
## 2. What CI measured is what is published                            chrome
```

A run starts at a line matching `taskUncheckedPrefix` or `taskCheckedPrefix` at
column zero, which is the definition `src/ui/checkbox.go` already uses and which
the scanner's totals already agree with. It continues while the next line begins
with whitespace and is not itself a checkbox. It ends on the first line that is
unindented, blank, a heading, or the next checkbox.

Tying the boundary to the same two prefixes matters. `checkbox.go` says already
that an indented or `*` prefixed checkbox is deliberately not one of ours,
because the scanner does not count it and a box drawn beside a number that
ignores it would be a lie. An indented `- [ ]` is therefore a continuation line
of the item above, not an item.

### Why it goes through itemLines

`src/ui/lines.go` exists for exactly this and argues its own case:

> Three lists in this application share one shape [...] Two instances was a
> coincidence. Three is where the arithmetic should be one thing.

`span` gives the highlight its range and `fit` gives the page key its item
count. The tasks pane computes the first by hand today and does not compute the
second at all. Wiring it to `itemLines` makes it the fourth caller and leaves
`lines.go` untouched.

`offsetFor` is deliberately not used, though it looks like the third fit. It
answers where the pane would sit if it were being laid out from scratch, which
is the right question for the outline, laid out afresh every frame, and the
wrong one for a document. A document holds a reading position across rescans
and resizes, which the document-viewer spec requires of it, and `offsetFor`
would throw that position away on every keystroke by scrolling further than the
cursor needed. The existing `scrollCursorIntoView` already scrolls the smallest
amount that shows the whole item and already shows the top of an item taller
than the pane, so it keeps its logic and is fed the item's span instead of a
source line's.

`specdetail.go` shows the pattern to copy: `ownersOfRows(outlineRows(...))`
builds the slice from whatever the renderer laid out, and `nodesInRows` wraps
`fit` for the page key. The tasks pane wants the same two functions over
`[]sourceLine`.

### Decision: the cursor visits tasks only

`m.docLines` today holds every source line. It becomes a list of items, with
chrome marked `lineOwner`. Consequences:

- `j` and `k` step task to task. Headings and blanks are passed over, and the
  scroll that follows the cursor brings them into view on the way.
- `G` lands on the last task rather than on the trailing empty line a file's
  final newline produces, which is an improvement nobody asked for and everybody
  wanted.
- A `tasks.md` with no checkboxes has no items, so `docHasCursor()` is false,
  `docActive()` is true, and `j`/`k` scroll the viewport the way they do on
  `proposal.md`. This falls out of the existing predicates with no special case.

The alternative, keeping a line cursor and only widening the highlight, was
rejected: `j` from the checkbox line would land on a continuation line without
moving the highlight, which reads as a dead key.

### Decision: a page moves items, not rows

`document-viewer` already states the rule for surfaces whose items are not one
row tall:

> Where a surface draws items whose rows are not all one row tall, a page SHALL
> move the cursor by as many items as fill that number of rows.

That requirement was written for the three lists and is extended to the
cursored document. Its sibling requirement, that the paging keys act on
whichever surface holds the keyboard, needs one clause: a document that holds
the keyboard scrolls, unless it has a cursor, in which case the cursor moves and
the view follows it. That mirrors the three-way split `j` and `k` already have
in `src/ui/ui.go`:

```
docHasCursor()  →  moveDocCursor(±1)          →  moveDocCursor(±docPage())
docActive()     →  docViewport.ScrollUp(1)    →  docViewport.PageUp()
else            →  list cursor                →  moveListCursor(listPage())
     j / k                                            pgup / ctrl+b / ...
     (today)                                          (missing today)
```

`gg` and `G` take the same split: `GotoTop`/`GotoBottom` for a cursorless
document, first and last item for a cursored one. The view follows, so the end
result on screen is the same and the cursor is where the eye is.

## The refresh

### Decision: split the flag rather than suppress the event

`m.scanning` conflates two situations:

```
      startup scan                    background refresh
      ─────────────                   ──────────────────
      screen is empty                 screen is full
      user is waiting                 user is reading, typing
      searching the filesystem        re-reading one known project
      ─────────────                   ──────────────────
      modal: honest                   modal: a lie
      refuse keys: nothing to do      refuse keys: loses input
```

The flag has seven uses, all in `src/ui/ui.go`: the declaration, the key guard,
two in the `fsChangeMsg` branch, the clear and the re-arm in `scanMsg`, the
modal, and the startup assignment. A second field for the refresh case touches
those and nothing else.

Two alternatives were weighed and rejected:

**Suppress the watcher's echo of our own write.** Remember the path and mtime
just written and ignore the next event for it. Racy, and it fixes less: the same
modal and the same swallowed keys appear when `tasks.md` is edited in another
window while spg is open, which is a workflow this tool invites by offering `E`.

**Update the model optimistically and let the rescan reconcile.** Zero perceived
latency, but it does not touch the modal or the key guard, because the watcher
fires either way. It is an improvement on top of the split, not instead of it,
and it costs a second copy of the truth about the file.

### Decision: a refresh draws nothing

Not a spinner in the nav bar gutter, not a dimmed frame. `src/ui/checkbox.go`
already argues the position for the status line:

> Success needs no message: the box changes shape, which is the feedback.

The same holds for the rescan behind it. If an external edit to a large project
turns out to be slow enough to be puzzling, a spinner in the nav bar is a small
addition later, and the nav bar is where `docScrollPercent` and the version
already live.

### Known and accepted: two reads per toggle

`toggleSelectedTask` scans immediately and the watcher scans again 200ms later.
The `scanPending` coalescing in the `scanMsg` branch does not catch this, because
the toggle's scan sets no flag and so does not participate. Routing the toggle's
scan through the refresh flag lets the watcher event coalesce into it when the
timing lines up. Worst case stays two reads of one project, which is not worth
more machinery than that until it measurably stutters.

## Where the refresh requirement lives

It is put in `fs-watch`, beside "Auto-rescan on filesystem changes in the open
project", which is the requirement that says a rescan refreshes the display.
`modal-presentation` names "the scan indicator" among the modals it covers, and
that stays true: the indicator is still a whole-frame modal, it is raised
in fewer situations. If a reviewer would rather see the rule in
`modal-presentation`, moving it costs nothing.

## Risks

- **The boundary rule meets a file it was not written for.** A task whose
  continuation is a fenced code block whose content is unindented would end the
  item early, leaving the fence as chrome. Harmless on screen, and the toggle is
  unaffected because it matches the checkbox line's text and nothing else.

- **A highlight that swallows a heading.** The failure mode of getting the
  boundary wrong is visible immediately and is what the boundary tests are for.

- **The key guard protects something.** Removing it for refreshes means a key
  can be handled against a model that is one scan behind. Every handler either
  reads the model as it is or acts on a file by name, and `toggleTaskInFile`
  already re-reads from disk and refuses a line it cannot identify, so a stale
  model cannot cause a wrong write. That refusal is the safety net, not the
  guard.
