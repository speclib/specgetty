## 1. The rule

- [x] 1.1 Point the paging and jump keys back at the focus-gated predicate, undoing `page-from-either-half`, and delete the predicate it added now that nothing asks it
- [x] 1.2 Give each of those keys a branch for the list that holds the keyboard, so a key that reaches neither a document nor the log moves a list
- [x] 1.3 Size a page to the rows the surface is currently showing, and half a page to half of it, so the keys mean the same thing at any terminal height. Proven by asserting a taller terminal moves further

## 2. The change list

- [x] 2.1 Move the selection by a page with `pgdown` and `ctrl+f`, and back with `pgup` and `ctrl+b`
- [x] 2.2 Move by half a page with `ctrl+d` and `ctrl+u`
- [x] 2.3 Move to the first and last change with `gg` and `G`
- [x] 2.4 Stop at either end rather than wrapping or running out of range. Proven by paging past both ends from every starting row
- [x] 2.5 Land on a change rather than a group header when a page crosses the boundary. Proven against a fixture whose page size puts a header mid-jump
- [x] 2.6 Remember the selection the way a single-row move does, so a rescan or a filter keeps it
- [x] 2.7 Move within what a filter left, not within the whole list

## 3. The side lists

- [x] 3.1 Page and jump the spec list while it holds the keyboard, leaving the spec content where it is
- [x] 3.2 Page and jump the properties rows the same way
- [x] 3.3 Keep the content paging when the content holds the keyboard, on both tabs. This is the half that `page-from-either-half` got right
- [x] 3.4 Handle a list shorter than a page by moving to the end and stopping

## 4. The picker

- [x] 4.1 Add the page and half-page keys to the picker's own key block
- [x] 4.2 Leave its `g` and `G` alone, which already work
- [x] 4.3 Move within what the picker's filter left

## 5. Verification

- [x] 5.1 Assert the whole rule as one table: for each surface that can hold the keyboard, which of the list cursor, the document offset and the log offset moves, and that the other two do not
- [x] 5.2 Confirm an open change still pages its artifact, since it was never a split
- [x] 5.3 Confirm the log panel still takes these keys while it holds the keyboard
- [x] 5.4 Confirm `j` and `k` are untouched on every surface
- [x] 5.5 Revert task 1.2 and confirm the change list paging test fails. Check the build succeeds first
- [x] 5.6 Open this project with the built binary and page down the 38-row change list, the 24-spec list and the properties rows
- [x] 5.7 `nix flake check` passes, coverage floors included

## 6. Notes

- [x] 6.1 Update the README where it lists the keys, which documents them per view
