## 1. Implementation

- [x] 1.1 Join table cells with two spaces instead of one, in both the header row and the body rows
- [x] 1.2 Make `layoutFields` charge two columns per gap instead of one, so the widths it hands back still add up to the table width. The gap count and the join string are the same fact written twice and must move together
- [x] 1.3 Leave `minFlexWidth` where it is. The wider gaps eat into the flexible column's budget, which means a column starts being dropped at a slightly wider terminal than before. That is the existing rule doing its job, not a new one

## 2. Verification

- [x] 2.1 A table whose values exactly fill their columns renders with two blank columns between them, asserted on the text
- [x] 2.2 The rendered row is still exactly the table width, at several widths. If the gap count and the join string disagree, the row is one or two columns wrong per gap and nothing else notices
- [x] 2.3 A test pins which columns survive at a given narrow width, so the shift in the drop threshold is recorded rather than discovered on a laptop
- [x] 2.4 The picker gets the same separation, since it shares `renderTable`
- [x] 2.5 Change the gap count without changing the join string and confirm 2.2 fails. Check the build succeeds first
- [x] 2.6 `nix flake check` passes, coverage floors included

## 3. Notes

- [x] 3.1 Independent of `pad-every-view`. Either can be reverted without the other, which is why it is not folded into it
