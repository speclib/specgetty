---
# specgetty-p4e2
title: upgrade bubbletea to v2
status: draft
type: task
priority: normal
created_at: 2026-09-15T17:49:32Z
updated_at: 2026-09-15T19:13:01Z
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
