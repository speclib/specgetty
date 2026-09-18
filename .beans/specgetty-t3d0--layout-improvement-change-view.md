---
# specgetty-t3d0
title: layout improvement change view
status: in-progress
type: task
priority: normal
created_at: 2026-09-18T12:39:53Z
updated_at: 2026-09-18T14:25:55Z
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
