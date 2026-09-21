## 1. The legend

- [ ] 1.1 Give the shared search prompt the width it is drawn in, which it does not take today
- [ ] 1.2 Name the three matchers while the prompt is focused and its query is empty. One wording serves both surfaces, since the contents sigil reaches artifact text in a change and file paths in a project
- [ ] 1.3 Drop the naming when the prompt has less room than it needs, rather than wrapping. Proven at the 60-column minimum on both surfaces, asserting the query and the count survive and the line does not grow
- [ ] 1.4 Show nothing once a character is typed, and nothing while the prompt is unfocused
- [ ] 1.5 Confirm the prompt still reports the result count in every one of those states

## 2. The failed search

- [ ] 2.1 Suggest the contents search, carrying the same term, when a query with no sigil matches nothing
- [ ] 2.2 Do the same for a query beginning with `'`, which is also a name search
- [ ] 2.3 Suggest nothing when a query beginning with `:` matches nothing
- [ ] 2.4 Name the term rather than the sigil, so the suggestion can be typed as read
- [ ] 2.5 Do all of the above in the project picker's empty message too

## 3. The match hint

- [ ] 3.1 Separate the hint from the table by the column gap, in the change list's grouped renderer and in the shared table the picker uses
- [ ] 3.2 Separate the header the same way, rather than leaving it to the last header's length
- [ ] 3.3 Prove it with a row whose last column is exactly full, which is what the archive date always is. This is the case that has been wrong since `date` became a default column

## 4. Verification

- [ ] 4.1 Open this project and run `:lipgloss`, confirming the hint is legible beside the date and the groups count what was found
- [ ] 4.2 Open the picker and run `:lipgloss`, confirming the same
- [ ] 4.3 Run a name search that fails on both surfaces and confirm the suggestion names the term
- [ ] 4.4 Confirm the legend appears on `/` and is gone after one keystroke, on both surfaces
- [ ] 4.5 Revert task 3.1 and confirm the gap test fails. Check the build succeeds first
- [ ] 4.6 `nix flake check` passes, coverage floors included

## 5. Notes

- [ ] 5.1 Update the README where it documents the sigils, which currently describes them only in prose
