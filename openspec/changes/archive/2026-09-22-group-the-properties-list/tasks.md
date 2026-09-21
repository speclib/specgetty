## 1. The shared arithmetic, extracted before it is used

- [x] 1.1 Add `src/ui/lines.go` with a slice of item indices, one entry per drawn
      line and `-1` for chrome, carrying `span`, `fit` and `offsetFor`, and
      verify unit tests cover a one-line-per-item list, an item spanning three
      lines, and an item taller than the pane
- [x] 1.2 Verify `offsetFor` reduces to `max(0, last-height+1)` for one-line
      items and clamps down to the item's first line when the item is taller than
      the pane, by a test asserting both arms
- [x] 1.3 Verify `fit` reduces to the row count when every item is one line, by a
      test, since that is the property the change list relies on
- [x] 1.4 Move the change list onto the helper, keeping the spacer-skip where it
      is, and verify the existing grouped-list tests pass unchanged
- [x] 1.5 Move the spec outline onto the helper, and verify the existing outline
      paging and wrapped-node scrolling tests pass unchanged
- [x] 1.6 Verify no behaviour changed, by running the suite before and after the
      move with no test edited. A test that had to be changed to pass is a
      behaviour change and belongs in the proposal rather than in this task

## 2. The sections

- [x] 2.1 Give `propSection` a group, and order the rows as configuration, store,
      then one per schema, and verify a test asserts the order for a project with
      two schemas and for one with none
- [x] 2.2 Rename the `project` row to `config`, and verify a test asserts the
      label and that the row still reports the configuration file and its source
- [x] 2.3 Lay the list out as drawn lines, a header per group and a blank line
      between groups, and verify a test asserts the drawn shape for a project
      with two schemas
- [x] 2.4 Draw a header for a group that holds nothing, and verify a test asserts
      the schemas header is present for a project with no schema rows
- [x] 2.5 Draw headers in `sectionHeaderStyle` and uppercase, rows indented by
      two, and verify a test asserts a header is not indented and a row is

## 3. The cursor

- [x] 3.1 Keep the cursor indexing sections rather than drawn lines, and verify a
      test asserts moving down from the last row of a group lands on the first
      row of the next
- [x] 3.2 Verify no key can leave the cursor on a header, by a test that presses
      every vertical key at every position and asserts the selected row is always
      a row
- [x] 3.3 Make `gg` and `G` land on the first and last row, and verify a test
      asserts neither lands on a header
- [x] 3.4 Scroll the list by drawn lines through the helper, and verify a test
      asserts a row below the fold is brought into view and that the pane does
      not open on a blank line

## 4. The width

- [x] 4.1 Give the list a floor of twenty columns, keeping the one-third cap, and
      verify a test asserts the width at 60, 80, 100, 140 and 200 columns
- [x] 4.2 Verify the list is the same width at 100 and at 200 columns and that
      every extra column goes to the content, by a test asserting both halves
- [x] 4.3 Verify neither half is drawn wider than the panel at the minimum
      terminal size, by a test at 60x20
- [x] 4.4 Verify the frame still fits at 60x20, 100x30 and 200x50 with the
      properties tab active

## 5. Verification

- [x] 5.1 Revert 3.1 and confirm the cursor lands on a header. Check the build
      succeeds first
- [x] 5.2 Revert 1.4 and confirm the change list's own tests fail, which is what
      says the helper is carrying the arithmetic rather than sitting beside it
- [x] 5.3 Verify the properties tab renders for every project on this machine
      without a panic or a row wider than its pane, by a test over a directory
      given by an environment variable and skipped when it is unset
- [x] 5.4 `nix flake check` passes, coverage floors included

## 6. Documentation

- [x] 6.1 Update the properties tab picture in README.md to the grouped list
- [x] 6.2 Add the change to CHANGELOG.md under Changed
