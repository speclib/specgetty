## Why

specgetty is on bubbletea v1.3.10, bubbles v1.0.0 and lipgloss v1.1.0. The v2
line is released and is where the work goes now.

Two concrete gains, beyond staying current.

**Light terminals.** The UI hardcodes 24 colour references across 11 ANSI
values, chosen for a dark terminal. `dimStyle` is `250`, a light grey, and it
renders every hint, every empty state, every match label and the nav bar
backdrop. On a light background it is close to invisible. v2 is the first
version that can do anything about this: `tea.RequestBackgroundColor()` returns
a message with `IsDark()`, and `lipgloss.LightDark(isDark)` picks between two
colours. There is no v1 equivalent worth having.

**Tests stop lying by default.** v1's `Style.Render` checks for a terminal and
strips styling when there is none, which quietly turns styling and alignment
tests into assertions about nothing. That is not hypothetical: on 2026-09-15 a
test asserting every line had the same width passed both with and without the
bug it guarded. v2 moved that gate to the output writer, so `Render` always
emits and the `withColor` test helper can be deleted rather than ported.

## What Changes

- Move to `charm.land/bubbletea/v2`, `charm.land/bubbles/v2` and
  `charm.land/lipgloss/v2`
- `View() string` becomes `View() tea.View`, and `tea.WithAltScreen()` becomes a
  field on the value it returns
- `case tea.KeyMsg` becomes `case tea.KeyPressMsg`
- The viewport API renames land: `YOffset` becomes a method, `ViewDown` and
  `ViewUp` become `PageDown` and `PageUp`, `LineDown` and `LineUp` become
  `ScrollDown` and `ScrollUp`
- `lipgloss.SetColorProfile` is gone; the `withColor` test helper goes with it
- Behaviour is preserved throughout. Nothing a user can observe changes.

Not in scope: the light and dark colour work the upgrade makes possible. That is
a separate change with its own decisions about which pairs of colours to use,
and bundling it here would make a behaviour-preserving migration into something
that has to be reviewed twice.

## Capabilities

No capability changes. Every requirement in `openspec/specs/` describes
behaviour that is identical before and after, so the change declares
`skip_specs: true` rather than restating them.

## Impact

- `src/ui/ui.go`: import paths, the key message type switch, `View`'s signature
  and `Run`'s program construction
- `src/ui/docview.go`, `src/ui/picker.go`: viewport method renames
- All six test files carrying `tea.KeyMsg{}` literals
- `src/ui/layout_test.go`, `docview_test.go`, `specstab_test.go`: the `withColor`
  helper is deleted
- `go.mod`: three direct dependencies replaced; `termenv` drops out
- `flake.nix` and `package.nix`: `vendorHash` changes and must change in both
