---
# specgetty-t3d0
title: layout improvement change view
status: completed
type: task
priority: normal
created_at: 2026-09-18T12:39:53Z
updated_at: 2026-09-18T14:38:18Z
---

Have a look at this gimped image. This is how it should look like design/specgetty-layout-improvement.png


## OpenSpec changes

Split into three so each can be reverted without the others:

1. `unify-focus-state` (tinychange, no visual change). One `focus` value
   replacing `activeView` plus `specsFocus`. Must land first.
2. `box-the-tab-content` (spec-driven). The mockup: the tab content gets its
   own border. Chrome above it, content inside it.
3. `light-the-focused-pane` (tinychange). The specs tab splits into one border
   per pane, lit by containment: a border is lit when the keyboard is inside it.

Only 2-before-3 is a hard dependency.

The `.xcf` working file is gitignored; the two PNGs are committed as the record.


## Summary of Changes

Shipped as three commits on top of v0.5.0, each revertible without the others:

| commit  | change                 | effect                             |
|---------|------------------------|------------------------------------|
| 35f160b | unify-focus-state      | none on screen                     |
| 588df8d | box-the-tab-content    | the mockup: content gets a border  |
| a64e380 | light-the-focused-pane | the specs tab lights its focus     |

Three things came out of the work that were not in the mockup:

- The config tab filename moved above the border and the search prompt stayed
  inside it, settled by a rule rather than case by case: a line that names the
  content is chrome, a line that reports on the content is content.
- The specs tab now says which half has the keyboard. `tab` has moved it since
  the log panel landed, and the only sign was the selected spec dimming.
- Collapsing `activeView` and `specsFocus` turned out to be load-bearing rather
  than tidy. The pair could describe "the log has the keyboard and so does the
  spec content", which was harmless while nothing read it, and would have been
  a visibly wrong lit border once each border asks where the keyboard is.

Cost: two rows and four columns everywhere. A 20-row terminal shows nine
changes where it showed eleven.

Verified by capturing 27 views at three terminal widths before and after each
step, and by reverting each edit in turn to confirm a test fails, with the
build checked first.

## Not done

The scroll percentage still sits in the panel title rather than on the box that
scrolls. Raised during exploration, never decided, so deliberately left alone.
`openspec/specs/change-list-view` still requires the panel title to report it,
so nothing is inconsistent.

## Note

`design/` including the 2.7 MB `.xcf` was committed in `80b30f9` before this
work started. `.gitignore` now carries `*.xcf`, which only affects future files.
