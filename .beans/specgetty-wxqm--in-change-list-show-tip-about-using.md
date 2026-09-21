---
# specgetty-wxqm
title: 'in change list show tip about using :'
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T16:55:17Z
updated_at: 2026-09-21T17:27:41Z
---

use colons to search inside change artifacts

## OpenSpec change

`teach-the-search-sigils` (`openspec/changes/teach-the-search-sigils/`),
validated strict. Not yet implemented.

Decisions taken during exploration:

- a legend in the search prompt while it is focused and empty, naming all three
  matchers, not just `:`. Dropped rather than wrapped when too narrow
- plus the failed-search moment: a name search that found nothing suggests the
  same term with the contents sigil. The legend teaches at the point of intent,
  this at the point of frustration
- both on the change list and in the picker, which share one prompt renderer
- rejected: a matcher toggle. It would replace the sigils rather than teach
  them, leaving two ways to do one thing

## Defect folded in

The `matched` hint is appended flush against the table. It looked right only
because the last default column used to be `specs`, whose values are one or two
characters in a five-wide column. `group-the-change-list` made `date` a
default, and it fills its ten-wide column exactly, so every `:` result now
reads `2026-09-18design`. Folded in because a tip that leads people to a
visibly broken column is worse than no tip.
