---
# specgetty-2kij
title: the archived column has no margin between the matched column when filtering.
status: completed
type: task
priority: normal
created_at: 2026-09-21T16:56:10Z
updated_at: 2026-09-21T18:01:34Z
---

## Summary of Changes

Fixed in `teach-the-search-sigils` (commit 1e35b56), where it was folded in
because that change teaches people to use `:`, and `:` is the only thing that
renders the hint column.

The cause was not the archived column: the hint was appended with no separator
at all, relying on the last column's own padding to stand in for one. `date`
fills its ten-wide column exactly, so there was none. The fix prepends the
column gap, so any last column behaves.
