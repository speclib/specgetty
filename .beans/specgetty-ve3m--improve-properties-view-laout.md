---
# specgetty-ve3m
title: improve properties view laout
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T21:38:26Z
updated_at: 2026-09-21T22:00:03Z
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
