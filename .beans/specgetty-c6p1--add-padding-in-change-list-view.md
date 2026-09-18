---
# specgetty-c6p1
title: add padding to every view
status: in-progress
type: task
priority: normal
created_at: 2026-09-17T17:35:49Z
updated_at: 2026-09-18T10:48:50Z
---

All views will look better with at least on char padding on the left and right sides.

The l1 views:
- change list view
- specs list
- config.yaml view

The l2 view:
- all markdown views
- the specs view

Am I missing views?


## OpenSpec changes

Split into three so each can be reverted without the others:

1. `unify-content-width` (tinychange, no visual change). Collapses the three
   copies of `m.width - 2` into one method. Must land first.
2. `pad-every-view` (spec-driven). The one-column inset on the panel, the
   header giving up its own horizontal padding, and the nav bar.
3. `widen-table-column-gaps` (tinychange). Two blank columns between table
   columns, in the change list and the picker, which share a renderer.

Top and bottom padding is out of scope, as you narrowed the bean to left and
right on 18 September.
