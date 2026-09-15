## 1. Dependencies

- [ ] 1.1 Replace the three modules with their v2 paths: `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2`. The `github.com/charmbracelet/*/v2` paths do not resolve
- [ ] 1.2 Pin `x/ansi` to whatever the v2 line requires rather than letting `go get` drag unrelated modules forward; check the diff of `go.mod` before accepting it
- [ ] 1.3 `go mod tidy`, then confirm `termenv` has dropped out, since lipgloss v2 no longer needs it
- [ ] 1.4 Recompute `vendorHash` and write it to **both** `package.nix` and `flake.nix`. A stale hash in either fails `nix flake check` and blocks every future ship

## 2. The message loop

- [ ] 2.1 `case tea.KeyMsg` becomes `case tea.KeyPressMsg` in `update`
- [ ] 2.2 Leave `key := msg.String()` and every `case "j":` style arm alone; `String()` still exists and production code never reads `msg.Type`, `msg.Runes` or `msg.Alt`
- [ ] 2.3 Rewrite the 83 `tea.KeyMsg{}` literals in the six test files to construct `tea.KeyPressMsg`
- [ ] 2.4 Check that the rewritten literals still produce the same `String()` values the handlers switch on, particularly for the ctrl and arrow keys

## 3. View and the program

- [ ] 3.1 Change `View() string` to `View() tea.View`
- [ ] 3.2 Set `AltScreen` on the returned value, replacing `tea.WithAltScreen()` in `Run`
- [ ] 3.3 Reduce `Run` to `tea.NewProgram(m)` with the options that survive
- [ ] 3.4 Confirm the existing `View` tests still pass unchanged. They assert structure rather than a golden file, so they should carry over; if one needs rewriting, that is a signal worth reading rather than a formality

## 4. Viewport

- [ ] 4.1 `docViewport.YOffset` becomes `YOffset()` at all 31 sites, most of them test assertions
- [ ] 4.2 `ViewDown`/`ViewUp` become `PageDown`/`PageUp`
- [ ] 4.3 `LineDown(n)`/`LineUp(n)` become `ScrollDown(n)`/`ScrollUp(n)`
- [ ] 4.4 `HalfPageDown` and `HalfPageUp` no longer return lines; drop the unused results
- [ ] 4.5 Re-check the scroll distance tests. They assert exact offsets after a page and a half page, so a changed definition of "page" shows up there and nowhere else

## 5. Styling

- [ ] 5.1 Delete the `withColor` helper and its 10 call sites. v2's `Render` always emits, so forcing a profile is no longer needed to stop the tests being vacuous
- [ ] 5.2 Remove the `termenv` import from the tests that used it
- [ ] 5.3 Leave the style definitions alone: `lipgloss.Color("236")` still parses ANSI indices and `Background`/`Foreground` accept what it returns

## 6. Verification

- [ ] 6.1 `go build ./...` and `go vet ./...` clean
- [ ] 6.2 The whole suite passes. Grep the output for FAIL rather than counting PASS lines, which cannot see a failure
- [ ] 6.3 `bash scripts/coverage-gate.sh` passes and the floors hold. A drop means behaviour was lost, not that the gate is inconvenient
- [ ] 6.4 `nix flake check` passes, which also proves both vendor hashes agree
- [ ] 6.5 Run `spg` and confirm the alt screen is taken: the UI must not be left in the scrollback after `q`. No test can check this
- [ ] 6.6 Run `spg` and exercise the things that only render: the picker overlay, a confirmation modal, scrolling an open change, the specs tab focus, the nav bar at a narrow width
- [ ] 6.7 Resize the terminal while a document is open and confirm it re-wraps and stays within the frame

## 7. Afterwards

- [ ] 7.1 Do not move `copy-change-name-and-path` to `tea.SetClipboard`. It is OSC 52 and fire-and-forget, so a refused write reports success, which one of that feature's requirements forbids. The reasoning is in this change's design
- [ ] 7.2 Note that light and dark adaptation is now possible, via `tea.RequestBackgroundColor()` and `lipgloss.LightDark`, and belongs in its own change
