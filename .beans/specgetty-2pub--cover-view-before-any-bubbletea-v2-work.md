---
# specgetty-2pub
title: cover View() before any bubbletea v2 work
status: completed
type: task
priority: normal
created_at: 2026-09-15T20:43:19Z
updated_at: 2026-09-15T20:47:21Z
---

`View()` has 0% coverage and holds 77 of the 97 statements needed to bring
`src/ui` from 72.4% to the 80% the ship gate names for core packages.

It also matters more than that number suggests: bubbletea v2 changes its
signature from `View() string` to `View() tea.View`, and moves `WithAltScreen()`
into a field on that struct. It is the single most affected function in the
`specgetty-p4e2` migration and currently nothing asserts what it produces.

## Why structural assertions, not snapshots

A golden-file comparison would break on every styling tweak and teach people to
regenerate it without reading. Assertions about structure survive cosmetic
change: a modal appears when its state is set, the size guard fires below 60x20,
overlays stack in the documented order.

## The precedent

On 2026-09-15 the picker's highlighted row was visibly misaligned while the whole
suite was green. A test asserting every line had the same width passed both with
and without the bug, because lipgloss pads the wrapped remainder; only counting
lines caught it. Tests over rendered output have to assert the thing that
actually breaks.

## Branches to cover

- Not yet sized: width or height still zero
- Too small: under 60 wide or 20 high
- The ordinary view, with and without the log panel
- Picker overlay
- Startup prompt (no project found)
- Scanning modal, error modal
- Archive, discard and export: confirming, running, result ok, result failed
- Archive and discard confirming, in their incomplete-tasks form
- Overlay precedence: a confirmation modal draws over the picker

## Summary of Changes

`View()` is now 100% covered. `src/ui` went from 72.4% to 78.8%, and the overall
total from 73.8% to 78.8%. Floors raised to match.

Test-only change: `src/ui/view_test.go`, no production code touched.

## What is asserted

Structure, not pixels. A golden-file comparison would break on every styling
tweak and teach whoever hits it to regenerate without reading.

- The two early returns: not yet sized, and below the documented 60x20 minimum,
  including that 60x20 itself renders
- The frame is exactly the terminal height and never wider, at three sizes and
  again with every overlay up at once
- The log panel appears only when visible
- Each overlay: picker, startup prompt, scanning, error
- All three action state machines through every state: confirming (both the
  plain and the incomplete-tasks form), running, result succeeded, result failed
- The export confirmation names its destination path, which the spec requires
- A confirmation modal takes precedence over the picker

## Found while doing it

`placeOverlay` ignores its `background` argument, so overlays replace the screen
rather than compositing onto it. A comment in `View()` describes a stacking
order that therefore cannot happen, and I wrote that comment during the
project-picker change without checking. Recorded as specgetty-ecli.

The test deliberately does not assert that the picker stays visible behind a
modal, and says why, so the current behaviour is not silently frozen as
intended behaviour.

## A verification mistake worth remembering

I first checked this work by counting `--- PASS` lines, which cannot see a
failure. One test was wrong (it assumed the plain confirmation branch while the
fixture had 1 of 4 tasks done) and I reported success before noticing. Grep for
FAIL, not for PASS.

## Where this leaves specgetty-p4e2

The bubbletea v2 bean is blocked on 70% overall and 80% core. Overall is met at
78.8%, scanner at 91.9%. `ui` at 78.8% is 1.2 points short, roughly 15
statements, and what remains is mostly untested branches of `update` plus `Run`
and `Init`, which need a TTY.

More to the point, the two functions v2 rewrites hardest were `View` (77
uncovered statements, now zero) and `update` (114 uncovered). Half of that risk
is now covered.
