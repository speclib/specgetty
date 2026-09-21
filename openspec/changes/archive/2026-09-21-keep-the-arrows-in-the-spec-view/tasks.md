## 1. Implementation

- [x] 1.1 Give `left` and `right` an arm for the spec detail level in `ui.go`, so
      neither falls through to the project tab bar. The handlers test
      `m.level == levelChange` and let every other level reach the `else`, which
      is what carried a third level into the second
- [x] 1.2 Prefer a `switch m.level` over the `else if` chain, so a level added
      later has to state what these keys do rather than inheriting the project
      tab bar by omission

## 2. Verification

- [x] 2.1 Press `right` and `left` in the spec detail view and assert
      `m.detailTab`, `m.focus` and `m.level` are all unchanged
- [x] 2.2 Press `right` with the card holding the keyboard and assert the card
      still holds it, which is the symptom a reader actually sees
- [x] 2.3 Press the keys several times, then `esc`, and assert the specs tab is
      active. This passes today because `esc` forces the tab, which is what
      masked the defect; it must still pass afterwards
- [x] 2.4 Assert the keys still switch the project tab bar on the change list and
      still switch artifact sub-tabs in an open change, neither being what this
      change touches
- [x] 2.5 Revert 1.1 and confirm 2.1 and 2.2 fail. Check the build succeeds first
- [x] 2.6 `nix flake check` passes, coverage floors included

## 3. Notes

- [x] 3.1 Bean `specgetty-pimg`. Shipped in `open-a-spec-in-detail`, commit
      `af367e7`, with no test pressing either key at the spec level
- [x] 3.2 `follow-the-spec-grammar` gives the spec level a second state for a
      file that does not parse. In that state there is no tree, so `splitTab()`
      is false and `defaultFocus()` returns `focusDetail`, which neither box
      claims. Fixing this first means that state is written against arrow keys
      that already behave
- [x] 3.3 The `tinychange` schema comes from
      https://github.com/speclib/openspec-tinychange-schema
