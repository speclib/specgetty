## 1. Dependencies

- [x] 1.1 Replace the three modules with their v2 paths: `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2`. The `github.com/charmbracelet/*/v2` paths do not resolve
- [x] 1.2 Pin `x/ansi` to whatever the v2 line requires rather than letting `go get` drag unrelated modules forward; check the diff of `go.mod` before accepting it
- [x] 1.3 `go mod tidy`, then confirm `termenv` has dropped out, since lipgloss v2 no longer needs it
- [x] 1.4 Recompute `vendorHash` and write it to **both** `package.nix` and `flake.nix`. A stale hash in either fails `nix flake check` and blocks every future ship
- [x] 1.5 Not foreseen: all three v2 modules declare `go 1.25.0`, which raised this module's own directive. The nix build was on Go 1.24.10 and refused with `go.mod requires go >= 1.25.0`. Fixed with `buildGo125Module`, which the already-pinned nixos-25.05 provides, so the toolchain moved without moving nixpkgs

## 2. The message loop

- [x] 2.1 `case tea.KeyMsg` becomes `case tea.KeyPressMsg` in `update`
- [x] 2.2 Leave `key := msg.String()` and every `case "j":` style arm alone; `String()` still exists and production code never reads `msg.Type`, `msg.Runes` or `msg.Alt`
- [x] 2.3 Rewrite the 83 `tea.KeyMsg{}` literals in the six test files to construct `tea.KeyPressMsg`
- [x] 2.4 Check that the rewritten literals still produce the same `String()` values the handlers switch on, particularly for the ctrl and arrow keys

## 3. View and the program

- [x] 3.1 Change `View() string` to `View() tea.View`
- [x] 3.2 Set `AltScreen` on the returned value, replacing `tea.WithAltScreen()` in `Run`
- [x] 3.3 Reduce `Run` to `tea.NewProgram(m)` with the options that survive
- [x] 3.4 Confirm the existing `View` tests still pass unchanged. They assert structure rather than a golden file, so they should carry over; if one needs rewriting, that is a signal worth reading rather than a formality

## 4. Viewport

- [x] 4.1 `docViewport.YOffset` becomes `YOffset()` at all 31 sites, most of them test assertions
- [x] 4.2 `ViewDown`/`ViewUp` become `PageDown`/`PageUp`
- [x] 4.3 `LineDown(n)`/`LineUp(n)` become `ScrollDown(n)`/`ScrollUp(n)`
- [x] 4.4 `HalfPageDown` and `HalfPageUp` no longer return lines; drop the unused results
- [x] 4.5 Re-check the scroll distance tests. They assert exact offsets after a page and a half page, so a changed definition of "page" shows up there and nowhere else

## 5. Styling

- [x] 5.1 Delete the `withColor` helper and its 10 call sites. v2's `Render` always emits, so forcing a profile is no longer needed to stop the tests being vacuous
- [x] 5.2 Remove the `termenv` import from the tests that used it
- [x] 5.3 Leave the style definitions alone: `lipgloss.Color("236")` still parses ANSI indices and `Background`/`Foreground` accept what it returns
- [x] 5.4 Not foreseen, and the largest behavioural change in the upgrade: `Style.Width(n)` now means the TOTAL rendered width, with border and padding counted inside it. In v1 it was the content width and the chrome was added outside. Every box in the application was affected

      Three tests caught it: the panel top border no longer matched its box, the
      picker box came out two columns narrow, and the scanning modal wrapped
      mid-phrase because its content area had silently shrunk by six columns.
      `renderPanel` now takes the total width and derives the content width from
      it, the picker's split was inverted, and each fixed modal asks for its
      text width plus `modalChrome`.

## 6. Verification

- [x] 6.1 `go build ./...` and `go vet ./...` clean
- [x] 6.2 The whole suite passes. Grep the output for FAIL rather than counting PASS lines, which cannot see a failure
- [x] 6.3 `bash scripts/coverage-gate.sh` passes and the floors hold. A drop means behaviour was lost, not that the gate is inconvenient
- [x] 6.4 `nix flake check` passes, which also proves both vendor hashes agree
- [x] 6.5 Run `spg` and confirm the alt screen is taken: the UI must not be left in the scrollback after `q`. No test can check this
- [x] 6.6 Run `spg` and exercise the things that only render: the picker overlay, a confirmation modal, scrolling an open change, the specs tab focus, the nav bar at a narrow width
- [x] 6.7 Resize the terminal while a document is open and confirm it re-wraps and stays within the frame

## 7. Afterwards

- [x] 7.1 Do not move `copy-change-name-and-path` to `tea.SetClipboard`. It is OSC 52 and fire-and-forget, so a refused write reports success, which one of that feature's requirements forbids. The reasoning is in this change's design
- [x] 7.2 Note that light and dark adaptation is now possible, via `tea.RequestBackgroundColor()` and `lipgloss.LightDark`, and belongs in its own change
