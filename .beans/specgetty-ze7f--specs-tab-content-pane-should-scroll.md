---
# specgetty-ze7f
title: specs tab content pane should scroll
status: completed
type: task
priority: normal
created_at: 2026-09-15T16:32:06Z
updated_at: 2026-09-15T18:49:40Z
blocked_by:
    - specgetty-5k45
---

The specs tab shows a spec list on the left and the spec content on the right (src/ui/ui.go:renderSpecsTab). The content pane truncates and does not scroll, because j/k already moves the spec list cursor.

Deferred out of the openspec change `scroll-markdown-documents`, which makes the change artifact pane and the config tab scrollable and gives every markdown pane real wrapping. The specs tab gets the wrapping fix from that change but not the scrolling.

The focus-axis question is already decided: `tab` toggles focus between the spec list and the content pane, and the focused pane gets the active border treatment renderPanel already applies. `tab` is free at that level, the handler at src/ui/ui.go:440 explicitly does nothing once inside a project.

Implementation is then the document-viewer capability applied to a third pane.

- [ ] add a focus field for the specs tab (list or content)
- [ ] bind `tab` to toggle it at levelProject on the specs tab
- [ ] render the focused pane with the active border colour
- [ ] route the vertical keys to the content viewport when it has focus
- [ ] reset scroll when the selected spec changes
- [ ] update the nav bar hints

## Summary of Changes

Shipped as openspec change `scroll-specs-tab`, archived to
`openspec/changes/archive/2026-09-15-scroll-specs-tab/`. Commit 7fe3610.

The specs tab now has a focus. `tab` moves the keyboard between the spec list
and the spec content, and on through the log panel when it is open. With the
list focused `j`/`k` change spec as before; with the content focused they scroll
it, and it takes the same page and jump keys as an open change.

### The archived answer did not survive contact

`scroll-markdown-documents` recorded the plan for this: "tab toggles focus
between the spec list and the content, and the focused pane gets the active
border treatment renderPanel already applies."

The first half held. The second did not. `renderPanel` colours the border of a
whole panel, and `project-picker` removed the split layout, so there is only one
panel and the two halves live inside it. There is no border to colour.

Focus is shown instead by which cursor is lit: the selected spec is highlighted
while the list holds the keys and dimmed while the content does. The reading
percentage in the title appears only when the content is focused, which gives a
second, independent cue.

### One seam worth knowing

`specsSplit` now owns the list/content width division, and both `renderSpecsTab`
and `docRegion` call it. If those two ever disagree about the content width, the
box re-wraps the rows and the reported position stops matching the screen. That
is the same failure the wrapping work existed to remove, so the split has one
owner.

### On the config tab

Reported as broken in the same breath as this. It is not: measured on this
project, the config document is 21 rows, so it scrolls in a terminal 24 or 30
rows tall and fits entirely at 40 or more, where nothing moves because nothing
needs to. If it ever fails to scroll in a *short* terminal, that is a real bug
and worth a separate report with the terminal height.

Coverage in `src/ui` went from 64.1% to 68.1%, total 61.2% to 64.4%.
