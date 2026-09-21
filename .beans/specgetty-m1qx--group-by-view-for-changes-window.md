---
# specgetty-m1qx
title: group by view for changes window
status: in-progress
type: feature
priority: normal
created_at: 2026-09-15T15:40:13Z
updated_at: 2026-09-21T16:38:09Z
---

A group by view should be available for suitable filters like [archived/active]

## What this actually turned into

Active and archived changes live in one list behind a three-state filter cycled
by `f`. That filter is the least visible state in the application: its only
indication is a label in the nav bar. Grouping says the same thing by position,
permanently, for both groups at once.

## OpenSpec change

`group-the-change-list` (`openspec/changes/group-the-change-list/`),
validated strict. Not yet implemented.

Decisions taken during exploration:

- grouping replaces the filter modes. `f`, `change_mode`, `--change-mode`
  and the three mode names all go; a config still carrying the key is reported
  rather than silently ignored
- two fixed orders, no sorting: active by name, archived newest first. Sorting
  stays with [[specgetty-vru8]], which already plans the columns that would want
  it, and needs a sort key per field because the rendered value sorts wrong
- one column set for both groups, since the date field already renders blank on
  active rows. `date` joins the defaults, the state column leaves them
- a group keeps its header and a count of zero when empty

## Defects this fixes, found while exploring

- the archive is ordered oldest first, so the change archived most recently sits
  at the bottom of 35 rows
- in `active+archived` with the default columns nothing indicates which rows
  are archived, which `change-list-view` requires
