---
# specgetty-nms0
title: UI and UX tasks should improve
status: in-progress
type: task
priority: normal
created_at: 2026-09-22T22:03:18Z
updated_at: 2026-09-22T22:48:34Z
---

- cursor doesn't follow when using - page-up/pdown
- switch to show only open tasks
- flicker when change checkboxes
- the selected task item should be completely highlighted, not the line only

## Scope after exploring

- [ ] cursor follows page-up/page-down, ctrl+f/b, ctrl+d/u, gg and G
- [ ] the selected task item is highlighted whole, not one source line
- [ ] the cursor visits tasks only, not headings and blank lines
- [ ] no flicker and no swallowed keys when a checkbox is toggled
- [x] decide what "show only open tasks" should be (dropped, see below)
- [ ] the nav bar advertises `space` and stops calling `jk` scroll on the tasks pane

## Decisions

**Dropped: show only open tasks.** Not wanted.

**The unit.** A task item is a `- [ ] ` or `- [x] ` line at column zero plus the
indented lines that follow it. It ends at the first line that is unindented,
blank, a heading, or the next checkbox. The rule can be widened later.

**The cursor visits tasks only.** Headings and blank lines are chrome, which is
how the three other lists in the application already model them
(`lineOwner = -1` in `src/ui/lines.go`).

## Why three of these are one change

`src/ui/lines.go` already carries the shape: a slice with one entry per drawn
line holding the item that line belongs to. The change list, the spec outline
and the properties list all use it. The tasks pane is the one cursored surface
that never got wired to it, which is why its highlight covers a source line
rather than an item and why its page keys move rows rather than a cursor.

## The flicker, as diagnosed

`m.scanning` means two different things and the second one is wrong:

- at startup the screen is empty, the filesystem is being searched, and a modal
  saying so is honest
- on a refresh the screen is full, one known project is being re-read, and the
  modal at `src/ui/ui.go:1660` blinks over everything

Worse, `src/ui/ui.go:356` drops every key except `q` and `ctrl+c` while the flag
is set, so ticking several boxes quickly loses presses. Toggling writes the
file, the watcher sees it 200ms later, and the flag goes up.

The fix is to split the flag: a startup scan keeps the modal and the key guard,
a refresh gets neither.
