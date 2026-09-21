---
# specgetty-m1qx
title: group by view for changes window
status: completed
type: feature
priority: normal
created_at: 2026-09-15T15:40:13Z
updated_at: 2026-09-21T16:48:01Z
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

## Summary of Changes

Shipped as `group-the-change-list` (commit 93a9a30), archived to
`openspec/changes/archive/2026-09-21-group-the-change-list/`.

- one table, grouped, active first, each group counting its rows and keeping its
  header at zero
- the filter modes are gone: `f`, `change_mode`, `--change-mode`,
  `ResolveListMode`, `Config.ChangeMode` and the three mode names. A config
  still carrying `change_mode` is reported at startup, since YAML discards
  unknown keys silently
- two fixed orders: active by name, archived newest first
- `date` joined the default columns, blank on active rows; `archived` left
  them as redundant against the group header
- `renderGroupedTable` is the change list's own, reusing `layoutFields` and
  `fitCell`, so `renderTable` and the picker are untouched

## Defects fixed

- the archive was oldest first, so the change archived most recently sat at the
  bottom of 36 rows
- in the combined view nothing indicated which rows were archived with default
  columns, which `change-list-view` required
- unrelated crash found while covering the configuration paths:
  `spg --config <broken.yml> <dir>` dereferenced a nil config

## Corrected along the way

Task 4.2 asked for `g` and `G` to reach the first and last change. They are
document keys and have never been bound on the change list, so there was no jump
for a header to catch; binding them would have been unrequested behaviour. The
task now pins what is true: stepping to either end stops on a change.

Coverage floors raised: src 85.9 to 87.7, scanner 94.8 to 94.9, ui 83.3 to 83.9,
total 86.0 to 86.5.
