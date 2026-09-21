## 1. Implementation

- [x] 1.1 Give `groupLine` a third kind for the blank line between groups, so it
      is neither a change nor a header rather than a header with no label
- [x] 1.2 Emit that line in `groupLines` before every group but the first,
      unconditionally, so an empty group above or below does not suppress it
- [x] 1.3 Draw it as an empty line in `renderGroupedTable`, unstyled, so no
      background or foreground is painted across the width
- [x] 1.4 When the computed offset would put that line at the top of the pane,
      start one line further down. The offset puts the cursor on the last visible
      row, so moving the window down by one keeps the cursor in range

## 2. Verification

- [x] 2.1 Amend `TestEveryDrawnRowIsAChangeOrAHeader` to allow the third kind and
      to assert each line is exactly one of the three, then confirm it still
      fails for a line that is both a change and a header
- [x] 2.2 Assert one blank line sits between the last active row and the
      `ARCHIVED` header in a rendered frame, and none between the column header
      and `ACTIVE`
- [x] 2.3 Assert the blank line is still drawn when the active group is empty and
      when the archived group is empty, both by a filter and by the project
      having none
- [x] 2.4 Assert `lineOfRow` and the scroll offset count the blank line, by
      stepping the cursor across the group boundary in a pane too short to hold
      both groups and checking the selected change is on screen at every step
- [x] 2.5 Assert that no cursor position produces a frame whose first body line is
      blank, by walking the cursor through the whole list at a height where the
      boundary can reach the top
- [x] 2.6 Assert `renderFrame` is still exactly `m.height` lines and no wider than
      `m.width` at 60x20, 92x30 and 120x50, which is where a line nobody budgeted
      for would show up
- [x] 2.7 Remove the offset adjustment from 1.4 and confirm 2.5 fails. Check the
      build succeeds first
- [x] 2.8 `nix flake check` passes, coverage floors included

## 3. Notes

- [x] 3.1 The blank line costs one of the thirteen body lines an 80x24 terminal
      has, so eleven changes become ten. It is on screen only while the group
      boundary is, which is when the list is short enough to afford it
- [x] 3.2 `listPage` spends the body's line count on changes, so the page key
      already skips the two changes the group headers displace, and this makes it
      three. That is bean `specgetty-rj7k`, not this change
