## 1. The inset

- [x] 1.1 Give `renderPanel`'s box style a one-column horizontal inset, and pass `renderDetailPanel` a width two columns smaller to match. The top border is drawn by hand from `width` and must still come out at exactly `width` columns
- [x] 1.2 Change the project header's `Padding(1, 1)` to `Padding(1, 0)`, keeping the blank line above and below
- [x] 1.3 Delete the hardcoded two-space indent in "No project selected. Press p to pick one."
- [x] 1.4 Confirm the log panel inherits the inset through the same `renderPanel` call, and that its viewport is set to the new content width

## 2. The nav bar

- [x] 2.1 Indent the nav bar's text one column on each side while its background still spans the full terminal width
- [x] 2.2 Reduce the width budget the key hints are fitted into by two, so a hint that no longer fits is dropped rather than pushing the version off the right edge
- [x] 2.3 Apply the same to the status message that replaces the hints

## 3. Verification

- [x] 3.1 Every view is exactly `m.height` lines with each overlay up. A two-column overflow shows as an extra line, not as a wide one, which is the only reason the picker misalignment was ever caught
- [x] 3.2 A document of known content wraps at the expected column, asserted on the text rather than on the rendered width
- [x] 3.3 `docRegion().width` equals what `renderFrame` passes to `renderDetailPanel`, at several widths including the 60-column minimum
- [x] 3.4 The first character of the project header, the tab bar, the table header and the first table row all sit in the same column
- [x] 3.5 A selected change row's highlight starts and ends one column inside the border
- [x] 3.6 The nav bar's background still covers column zero and the last column while its text does not
- [x] 3.7 Render all twenty-seven captures again and diff against `views-before-v0.5.0.txt` in the scratchpad. Read the diff; every line should differ by the inset and nothing else
- [x] 3.8 Reintroduce each of the four edits in turn and confirm a test fails. Check the build succeeds first: a revert that does not compile makes `go test` print nothing, which reads like a pass
- [x] 3.9 `nix flake check` passes, coverage floors included

## 4. Notes

- [x] 4.1 Ships after `unify-content-width` and before `widen-table-column-gaps`. The order matters only for the first: without it the content width has to be changed in three files by hand
