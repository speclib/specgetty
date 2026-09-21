---
# specgetty-m3u5
title: one line margin between active and archived changes
status: todo
type: task
priority: normal
created_at: 2026-09-21T16:59:57Z
updated_at: 2026-09-21T17:36:46Z
---

One blank line between the ACTIVE and ARCHIVED groups in the change list.

Specified as OpenSpec tinychange `space-the-change-groups`.

Settled during exploration:

- The gap is unconditional: drawn whether or not either group has rows.
- It is a third kind of `groupLine`, not a header with an empty label. The test
  `TestEveryDrawnRowIsAChangeOrAHeader` exists to catch exactly that shortcut.
- It is a real line in the slice, not a newline added while rendering, because
  `renderGroupedTable` budgets in lines and the frame is asserted to be exactly
  the terminal height.
- No gap above the first group.
- When the scroll offset would put the gap at the top of the pane, the pane
  starts at the group header below it instead.

Found while tracing the line budget, split off as bean `specgetty-rj7k`: the
page key already spends the body's line count on changes, so it skips the two
changes the group headers displace, and three once this lands.
