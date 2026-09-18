## 1. The border

- [x] 1.1 Wrap the tab content in a bordered box in `renderDetailPanel`, directly under the tab bar with no blank row between them, with no title on it and the panel's existing dim border colour
- [x] 1.2 Give the box the same one-column inset the panel has, so the content is inset from both borders and `panel-layout`'s existing requirement still holds inside the new box
- [x] 1.3 Take two rows off the content height budget and four columns off the width. Both come from `panelContentWidth` and the `height - 5` in `renderDetailPanel`, which become the box's outside dimensions rather than the content's
- [x] 1.4 Do the same for an open change: the artifact document goes in the box, the change header line and the sub-tab row stay above it

## 2. Chrome and content

- [x] 2.1 Leave the config tab's source filename above the border, beside the tab bar, and take it out of the region handed to the box
- [x] 2.2 Keep the changes tab's search prompt inside the border, under the table
- [x] 2.3 Check the remaining content areas land on the right side of that line: the empty states, the unimplemented-tab placeholder, and the "No project selected" prompt

## 3. Verification

- [x] 3.1 Every view is exactly `m.height` lines and `m.width` columns with the box drawn, at 100x30, 72x24 and the 60x20 minimum. A two-column overflow shows as an extra line, not a wide one
- [x] 3.2 The box's top border sits on the row immediately after the tab bar, asserted on the rendered rows rather than on the arithmetic
- [x] 3.3 The config filename is outside the box and the search prompt is inside it, both asserted by which side of the border row they fall on
- [x] 3.4 A document still wraps at the new content width, asserted on its text, and `docRegion` still agrees with what the box is drawn at. That agreement is what `unify-content-width` exists to hold and it now has a second border to survive
- [x] 3.5 The change table keeps all three columns at the 60-column minimum, so the narrower content has not silently dropped one
- [x] 3.6 Capture every view before and after and diff. Every line should differ by the border and the two rows, and nothing else
- [x] 3.7 Reintroduce each edit in turn and confirm a test fails. Check the build succeeds first: a revert that does not compile makes `go test` print nothing, which reads like a pass
- [x] 3.8 `nix flake check` passes, coverage floors included

## 4. Notes

- [x] 4.1 Ships after `unify-focus-state` and before `light-the-focused-pane`. Only the second ordering matters: the focus change needs these borders to exist
