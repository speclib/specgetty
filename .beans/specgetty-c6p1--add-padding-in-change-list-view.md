---
# specgetty-c6p1
title: add padding to every view
status: completed
type: task
priority: normal
created_at: 2026-09-17T17:35:49Z
updated_at: 2026-09-18T10:59:29Z
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


## Summary of Changes

Shipped as three commits on top of v0.5.0, each revertible without the others:

| commit  | change                    | effect                          |
|---------|---------------------------|---------------------------------|
| 5d3c4c4 | unify-content-width       | none on screen                  |
| 64c0a22 | pad-every-view            | the gutter, 13 views            |
| 131e92c | widen-table-column-gaps   | two columns between table cells |

`renderPanel` now carries the inset, so every view inside it gets one column of
air on each side without rendering any of its own. Thirteen views, the log panel
included, and any view added later.

Three things came out of doing it that were not in the bean:

- The project header carried `Padding(1, 1)` of its own. Under the panel inset
  it sat two columns in while everything below it sat at one. It is `(1, 0)`
  now: the vertical half is what gives it room above and below, which you
  already judged sufficient.
- The nav bar filled the space between the key hints and the version with bare
  unstyled spaces, so its background had a hole in the middle. Now painted.
- Two tests asserted the panel content width by restating `m.width - 2` rather
  than asking for it. They now derive it, which is what `unify-content-width`
  exists to make possible.

Verified by capturing 27 views at 100x30, 72x24 and 60x20 before and after each
step: line counts unchanged and every line exactly the terminal width. Each edit
was also reverted in turn to confirm a test fails, with the build checked first.

## Not done

Top and bottom padding, which you removed from the ticket on 18 September
having reached the same conclusion: there is already enough room there.
