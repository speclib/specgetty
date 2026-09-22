---
# specgetty-nms0
title: UI and UX tasks should improve
status: completed
type: task
priority: normal
created_at: 2026-09-22T22:03:18Z
updated_at: 2026-09-22T23:01:05Z
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

## Summary of Changes

Shipped as `e95cc1b`, OpenSpec change
`openspec/changes/archive/2026-09-23-select-the-whole-task`.

Three of the four items were one fault: the tasks pane selected a source line,
and a task is not a source line. `src/ui/taskitems.go` groups the rendered
lines into task items and hands the arithmetic to `itemLines` in
`src/ui/lines.go`, which the change list, the spec outline and the properties
list already use. The tasks pane is now the fourth caller.

- **cursor doesn't follow page-up/pdown.** `pgdown`, `pgup`, `ctrl+f`,
  `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` now split three ways the way `j`
  and `k` always did: a document with a cursor moves its cursor, one without
  moves its rows, anything else is a list. `docPage()` sits beside `listPage()`.
- **the selected task item should be completely highlighted.** The band covers
  every row of the item, the checkbox line and its indented continuations, and
  stops there.
- **switch to show only open tasks.** Dropped, recorded as a non-goal in the
  proposal with the reasons.
- **flicker when change checkboxes.** `m.scanning` meant two things. It is now
  two fields. A first scan keeps the modal and the key guard; a refresh of a
  project already on screen gets neither, so a burst of toggles loses no
  keystroke. The toggle's own read goes through the refresh flag so the
  watcher's notice of the same write coalesces rather than queueing a third read.

Also: the cursor visits tasks only, passing over headings and blank lines; the
nav bar offers `space toggle` and calls `jk/↑↓` navigate on that pane; the
`changes.gif` demo was re-recorded because it reaches the tasks artifact.

Gate passed: total coverage 91.0% against a 90.5% floor, ui 90.1% against 89.6%.
