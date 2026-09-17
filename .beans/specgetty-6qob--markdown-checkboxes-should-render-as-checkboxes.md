---
# specgetty-6qob
title: markdown checkboxes should render as checkboxes
status: todo
type: feature
priority: normal
created_at: 2026-09-15T15:55:40Z
updated_at: 2026-09-17T17:24:41Z
---

the checkboxes should be rendered as utf8 checkboxen and I want the lines to be selectable, highlighted. Space should change the state of the checkbox. Saving should be atomic at every change

## Planned (2026-09-17)

OpenSpec change `toggle-task-checkboxes` (spec-driven) holds proposal, design,
spec deltas and 41 tasks. New capability `task-checkboxes`; `document-viewer`
modified.

## Decisions

| Topic         | Decision                                                     |
| ------------- | ------------------------------------------------------------ |
| scope         | `tasks.md` only                                               |
| glyphs        | `▢` U+25A2 and `▣` U+25A3                                     |
| colour        | deferred to a wider markdown highlighting change              |
| undo          | none; `space` is its own inverse                              |
| `j`/`k`       | move the cursor a source line at a time, view follows         |
| save          | content-addressed, atomic, mode preserved                     |

## Why not the ballot boxes

`☐` and `☑` were the first pick and were rejected after looking at them:
Ghostty draws `☑` as a coloured emoji while `☐` stays a text glyph, so the
sizes do not match and the checked state is wider than `ansi.StringWidth`
reports. A wrapper counting one cell while the terminal draws two drifts a
column on every wrapped line.

`▢` and `▣` are both one cell, both from Geometric Shapes, and neither has an
emoji presentation. Since colour is deferred, the filled centre carries the
state on its own, which also survives the green selected row.

## The measurement that shaped the design

At an 86 column pane, across 511 real task lines: 33% occupy one screen row,
60% two, 6% three, a few four. A task is not a screen row, so the cursor moves
by source line and the highlight covers every row that line produced. The
renderer currently discards that mapping on purpose, so the change adds it back
as one structure that serves both the highlight and the save.

Dropping `- [ ]` (three cells) to one cell returns two columns to the text and
moves one-row tasks from 33% to 36%.

## The risk the design exists to close

This is the first time specgetty writes to a file the user also edits; its
writes so far are its own cache, a directory rename, and a zip written
elsewhere. The filesystem watcher exists precisely because the user edits these
files in an editor, so a whole-buffer write from memory could discard an editor
save.

The save therefore re-reads the file, finds the line by its exact text rather
than by counting checkboxes, and refuses when the line is missing or ambiguous.
Counting would survive the re-read but not an insertion above the cursor.

## Open question left in the design

Whether the cursor stops on every source line or only on checkbox lines. The
design assumes every line, which matches "the lines should be selectable" at the
cost of pressing `j` through blank lines between task groups. Cheapest thing to
change after using it once.
