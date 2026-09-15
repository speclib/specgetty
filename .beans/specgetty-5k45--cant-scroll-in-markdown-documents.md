---
# specgetty-5k45
title: can't scroll in markdown documents
status: completed
type: bug
priority: normal
created_at: 2026-09-15T15:54:19Z
updated_at: 2026-09-15T18:04:21Z
---

cursors and vim keys should work and page down / up

Proposal written: `openspec/changes/scroll-markdown-documents`.

Two defects behind this, not one:

1. No scroll offset anywhere. Every content pane ends in truncateContent (src/ui/ui.go:1563), which cuts the line slice at the panel height. At levelChange the vertical keys are excluded by `m.level != levelChange` (src/ui/ui.go:535, 562) and the page and jump keys move m.fileCursor, a leftover from a file listing with no renderer.
2. renderMarkdown takes a width and never uses it (src/ui/ui.go:1447). The panel box wraps the already-truncated string afterwards (lipgloss style.go:368) and MaxHeight cuts the overflow, so a long paragraph pushes its neighbours off screen. Content is lost inside the visible height too.

Decided: full pages in documents only, lists keep their half page; percentage in the panel title, no scrollbar; no mouse wheel (bubbletea capture would take native text selection away). Specs tab deferred to specgetty-ze7f.

Overlaps with the active project-picker change, which deletes fileCursor and renumbers the level constants. Sequence them rather than running both at once.

## Summary of Changes

Shipped as openspec change `scroll-markdown-documents`, archived to
`openspec/changes/archive/2026-09-15-scroll-markdown-documents/`. Commit 86aa405.

### Two defects, not one

Scrolling was the reported problem. The deeper one was that `renderMarkdown`
took a `width` parameter and never used it, so nothing wrapped. The lipgloss box
wrapped afterwards and `MaxHeight` cut the overflow, which lost content inside
the visible height, not only below it.

Wrapping now happens in `renderMarkdown` and `renderYAML`, after styling, so the
rows the renderer produces are the rows the terminal shows. Only then can a
viewport report a position that matches the screen.

### Decisions worth keeping

`ansi.Wrap`, not `ansi.Wordwrap`. Wordwrap lets a token longer than the limit
overflow, and an overflowing row gets re-wrapped by the panel box, which puts
the row count back out of step. Breakpoints include `/` so store paths and URLs
break sensibly.

Styles are re-opened on every wrapped row (`reopenStyles`). `ansi.Wrap` leaves
the opening sequence on the first row only and relies on terminal state carrying
across the newline. That holds when rows are printed together and breaks the
moment a viewport slices them: scroll until a continuation row is at the top and
it renders unstyled. Re-opening costs no display width.

One viewport keyed by document identity (project, tab, change, artifact). One
rule covers every reset: different document starts at the top, same document
keeps its place. That is what makes the file watcher bearable while reading.

`Update` is wrapped so `syncDocument` runs on every path, including the early
returns the overlays take. Otherwise the viewport holds stale rows.

### Deviation from the task list

Task 8.5 asked to confirm the project list and change list still page by half a
page. That premise expired with `project-picker`: those handlers belonged to the
old project list panel and its dead file listing, both removed. No list binds
the paging keys now. The test asserts the true behaviour instead, and says why.

Task 7.2 asked to remove `truncateContent` if it had no callers left. It has
five, so it stays.

### Verification

Each group of tests was checked by reintroducing the bug it guards and
confirming failure. Removing the wrapping fails five tests; removing the
reset rule fails one. A first attempt at that check was itself wrong: the
revert broke the build, so `go test` never ran and the grep found no failures.

Coverage in `src/ui` went from 59.0% to 64.1%, total 56.9% to 61.2%.
