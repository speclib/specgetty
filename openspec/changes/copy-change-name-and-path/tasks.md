## 1. Clipboard

- [ ] 1.1 Promote `github.com/atotto/clipboard` from indirect to direct in `go.mod`; it is already in the module graph via `bubbles`
- [ ] 1.2 Copy through `clipboard.WriteAll`, and treat both a returned error and the package-level `clipboard.Unsupported` as failure, since the two indicate different problems and neither should look like success

## 2. The keys

- [ ] 2.1 Bind `y` to copying the selected change's name; both `y` and `Y` are currently unbound
- [ ] 2.2 Bind `Y` to copying the selected change's absolute directory path
- [ ] 2.3 Gate both on the changes tab having a selected row, and on no overlay holding the keyboard, the way `a`, `d` and `e` already are
- [ ] 2.4 Build the path from `ChangeInfo.DirName`, never from `Name`

      DirName is the directory on disk; for an archived change it keeps the
      `YYYY-MM-DD-` prefix that Name has stripped for display. Using Name here
      produces a path that does not exist. That is exactly the defect
      `doExportChange` shipped with from f2f6105 until 2026-09-15, and it was
      invisible because the wrong path only fails at the moment someone uses it.
      A path copied to the clipboard fails even later, in a shell, far from the
      tool that produced it.

- [ ] 2.5 Add the hints to the nav bar for the changes tab, alongside the existing action hints

## 3. The status line

Nothing in specgetty reports a result without a modal today; archive, discard
and export all require a keypress to dismiss. That weight is wrong for a copy,
and silence is worse, because a failed copy would be indistinguishable from a
successful one until the user tries to paste.

- [ ] 3.1 Add a status message field to the model and render it in place of the nav bar while it is set
- [ ] 3.2 Clear it on the next key press rather than on a timer, so there is no `tea.Tick` and no re-render loop; the message lasts exactly as long as the user's attention does
- [ ] 3.3 Use the same line for failure, naming what went wrong
- [ ] 3.4 Keep it general rather than copy-specific. If a second caller appears, lift it into its own capability; it is specified here only because this is the change that needs it

## 4. Tests

- [ ] 4.1 `y` copies the display name; `Y` copies the absolute path
- [ ] 4.2 The path of an archived change contains the date prefix and resolves on disk, built against a real temporary project rather than a string comparison
- [ ] 4.3 The path of an active change resolves on disk
- [ ] 4.4 Neither key acts with no row selected, on another tab, or while the picker or the search prompt holds the keyboard
- [ ] 4.5 A failed copy produces a status message that says so, by injecting a failing writer rather than by breaking the machine's clipboard
- [ ] 4.6 The status message appears in `View()` output and is gone after the next key
- [ ] 4.7 The nav bar shows the hints with a row selected and omits them when the list is empty
- [ ] 4.8 Inject the clipboard writer so the suite never touches the real system clipboard; a test that clobbers the developer's clipboard is a bad neighbour

## 5. Verification

- [ ] 5.1 `go build ./...` and `go vet ./...` clean
- [ ] 5.2 `bash scripts/coverage-gate.sh` passes, floors raised
- [ ] 5.3 `nix flake check` passes
- [ ] 5.4 Copy a real archived change's path from this repository and confirm the result is a directory that exists, since a path that does not resolve is the whole failure mode this change guards against
- [ ] 5.5 Run `spg` and press both keys, then paste. The clipboard reaches the terminal and the window manager, neither of which the test suite can see

## 6. Notes

- [ ] 6.1 Decided against OSC 52: it wins only when running over SSH, which is not a case that matters here, and `atotto` is verified working on this machine (`wl-copy` and `wl-paste` both present, and its `init()` checks `WAYLAND_DISPLAY` first)
- [ ] 6.2 Decided not to wait for bubbletea v2 (`specgetty-p4e2`), which brings its own OSC 52 clipboard. With SSH out of scope, v2 offers nothing this needs
- [ ] 6.3 The project picker is out of scope; copying a project path is a sibling feature nobody has asked for
- [ ] 6.4 The `tinychange` schema comes from https://github.com/speclib/openspec-tinychange-schema and is vendored at `openspec/schemas/tinychange/`
