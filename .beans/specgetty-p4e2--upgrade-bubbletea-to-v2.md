---
# specgetty-p4e2
title: upgrade bubbletea to v2
status: completed
type: task
priority: normal
created_at: 2026-09-15T17:49:32Z
updated_at: 2026-09-16T13:42:08Z
blocked_by:
    - specgetty-sreo
    - specgetty-17c3
---

see https://github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md

## Deferred until the coverage gate is met (2026-09-15)

Explored and deliberately parked. The migration is broad and mechanical, and
what makes a broad mechanical refactor safe is the test suite it lands on. At
64.4% total it is not safe enough yet.

Blocked by reaching the target the ship gate already names: **70% overall, 80%
core** (scanner and ui). Current: total 64.4%, scanner 53.3%, ui 68.1%.

## What the upgrade would actually buy

**Light and dark terminals.** specgetty hardcodes 24 colour references across 11
ANSI values, tuned for a dark terminal. `dimStyle` is `250`, a light grey, and it
renders every hint, empty state, match label and the nav bar backdrop; on a
light background it is close to invisible. v2 is the first version that can fix
this: `tea.RequestBackgroundColor()` returns a `BackgroundColorMsg` with
`IsDark()`, and `lipgloss.LightDark(isDark)` picks between two colours. There is
no v1 equivalent worth having.

**Honest tests by default.** v1 `Style.Render` checks for a terminal and strips
styling when there is none, which silently turns styling and alignment tests
into no-ops. That trap was hit for real on 2026-09-15: an "all lines have the
same width" test passed both with and without the bug it guarded. v2 moved the
gate to the output writer, so `Render` always emits and the `withColor` helper
(10 call sites across 3 test files) can be deleted rather than ported.

**Being on the maintained line.** Weakest argument, usually the deciding one.

## What it would NOT buy

Mouse support is still not worth taking. v2 moves mouse config into
`View.MouseMode`, but any mouse reporting mode still takes native text selection
away from the terminal. That is terminal behaviour, not a bubbletea choice, so
the non-goal recorded in `scroll-markdown-documents` stands.

`WindowTitle` and the native `ProgressBar` field are a nice fit (project name in
the title bar, a real progress bar during the picker's multi-second walk) but
neither justifies the move.

## Migration surface, measured

| Area                              | Count |
| --------------------------------- | ----- |
| `tea.KeyMsg` literals in tests     | 82 across 5 files |
| Non-test bubbletea touchpoints     | 9 in ui.go |
| `docViewport.YOffset` uses         | 31, field becomes a method |

Other renames: `ViewDown`/`ViewUp` become `PageDown`/`PageUp`,
`LineDown`/`LineUp` become `ScrollDown`/`ScrollUp`, `HalfPageDown` stops
returning lines, `View() string` becomes `View() tea.View`, and
`tea.NewProgram(m, tea.WithAltScreen())` becomes `tea.NewProgram(m)` with
`v.AltScreen = true` set inside `View`.

## Trap worth knowing before starting

The v2 modules moved off GitHub import paths:

    github.com/charmbracelet/lipgloss/v2   ->  charm.land/lipgloss/v2
    github.com/charmbracelet/bubbletea/v2  ->  charm.land/bubbletea/v2
    github.com/charmbracelet/bubbles/v2    ->  charm.land/bubbles/v2

The upgrade guide does not mention this. `go get` on the github path fails with
"module declares its path as: charm.land/...". Worth a thought about depending
on a vanity domain.

Versions available at time of writing: bubbletea v2.0.9, bubbles v2.2.1,
lipgloss v2.0.6.

## Proposal written (2026-09-15)

`openspec/changes/upgrade-bubbletea-v2/` holds proposal, design and tasks
(spec-driven, `skip_specs: true` because behaviour is preserved). 31 tasks.

## What the second exploration changed

The first pass measured 82 key literals and treated the migration as large. The
second pass checked what production code actually uses, and the surface is much
smaller than the count implies:

- No binding for the space key, so v2's `String()` change for space is moot
- Production code never reads `msg.Type`, `msg.Runes` or `msg.Alt`, only
  `msg.String()`, so the struct-to-interface change touches one type switch
- No mouse or focus reporting, so the removed program options do not apply
- All 83 key literals are in tests

`View()` went from 0% to 100% coverage in `specgetty-2pub`, which removes the
single largest blind spot: it is the function whose signature changes.

## Blocking status

The bean was blocked on 70% overall and 80% core. Overall is met at 78.8% and
scanner at 91.9%. `ui` at 78.8% is 1.2 points short, roughly 15 statements, and
what remains is untested `update` branches plus `Run` and `Init`, which need a
TTY and cannot be covered here.

Whether that 1.2 points is worth waiting for is a judgement call. The risk it
was standing in for, `View` being untested, is gone.

## The thing to carry forward

Do not switch `copy-change-name-and-path` to `tea.SetClipboard` afterwards. It
is OSC 52 and fire-and-forget, so a refused write reports success, which that
feature's requirements forbid. Recorded in the change's design and as task 7.1.

## Summary of Changes

Shipped as openspec change `upgrade-bubbletea-v2`, archived to
`openspec/changes/archive/2026-09-16-upgrade-bubbletea-v2/`. Commit 56065aa.
31 tasks. `nix flake check` passes.

Now on `charm.land/bubbletea/v2` v2.0.9, `charm.land/bubbles/v2` v2.2.1 and
`charm.land/lipgloss/v2` v2.0.6. `termenv` dropped out entirely.

## Two things the plan did not foresee

**lipgloss v2 redefined Style.Width.** It now means the TOTAL rendered width,
with border and padding counted inside it; in v1 it was the content width and
the chrome was added outside. Every box in the application was affected. This is
the largest behavioural change in the upgrade and the guide does not mention it.

Three tests caught it, which is exactly why they were written:

- `TestPanelTopBorderMatchesTheBoxWidth`: the hand-built top border no longer
  matched its box
- `TestPickerDoesNotOverflowItsBox`: the picker came out two columns narrow
- `TestViewDrawsTheScanningModal`: the modal wrapped mid-phrase, because its
  content area had silently lost six columns

`renderPanel` now takes the total width and derives the content width from it,
the picker's split was inverted, and each fixed modal asks for its text width
plus `modalChrome`.

**The v2 modules require Go 1.25.** That raised this module's own directive, and
the nix build was on Go 1.24.10: `go.mod requires go >= 1.25.0`. Fixed with
`buildGo125Module`, which the already-pinned nixos-25.05 provides, so the
toolchain moved without moving nixpkgs.

## What went smoothly, and why

The production surface really was one type switch. Nothing reads `msg.Type`,
`msg.Runes` or `msg.Alt`, nothing binds space, and every key name v2 produces
("esc", "pgdown", "tab") matches what the handlers already switch on. The 94
key literals were all in tests.

`View()` being at 100% coverage before this started is what made the rendering
changes safe to make. It is also now asserted that `View().AltScreen` is set,
which is a smaller claim than "the alt screen works" but is the part a test can
make.

## Still unverified, and only you can do it

- The alternate screen actually being taken. The flag is asserted; whether the
  terminal honours it is not something any test here sees. If it is wrong,
  specgetty draws over the scrollback.
- The program starting and tearing down at all: `Run` has no coverage and needs
  a TTY.
- A live resize, the picker overlay, and the nav bar at a narrow width.

## Next

Light and dark adaptation is now possible via `tea.RequestBackgroundColor()` and
`lipgloss.LightDark`, and was deliberately left out of this change.

Do not move the clipboard feature to `tea.SetClipboard`: it is fire-and-forget
OSC 52, so a refused write would report success, which `copy-to-clipboard`
forbids.
