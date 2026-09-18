---
# specgetty-ecli
title: placeOverlay ignores its background argument
status: completed
type: bug
priority: low
created_at: 2026-09-15T20:46:59Z
updated_at: 2026-09-18T09:28:05Z
---

`placeOverlay` takes a `background` argument and never uses it.

    func placeOverlay(width, height int, modal, background string) string {
        return lipgloss.Place(
            width, height,
            lipgloss.Center, lipgloss.Center,
            modal,
            lipgloss.WithWhitespaceBackground(lipgloss.NoColor{}),
        )
    }

`lipgloss.Place` centres the modal on blank space. The view passed in as
`background` is discarded, so every modal replaces the whole screen rather than
being drawn over it. Present since the initial commit (10b2062).

## Why it is worth deciding rather than leaving

`View()` carries a comment describing a stacking order that cannot happen:

    // The picker is an overlay over the current view. Confirmation modals are
    // drawn after it, so they sit on top.

There is no "on top". Drawing the picker and then a modal discards the picker
entirely. That comment was written during the `project-picker` change on the
assumption that `placeOverlay` composites, without checking that it does.

Found while covering `View()` (specgetty-2pub): a test asserting the picker was
still visible behind a confirmation modal failed, and the rendered output was a
blank screen with one modal on it.

## The decision

Either behaviour is defensible. A full-screen modal is a normal design, and
specgetty has always looked like this.

- **If it is intended**: delete the unused parameter and fix the comment in
  `View()`, so the code stops promising something it does not do.
- **If it is not**: composite the modal onto the background, which is what the
  name and the signature say. `lipgloss.Place` cannot do this; it needs a
  line-by-line overwrite of the background with the modal's box.

The second is more work and changes how every modal looks. The first is a few
minutes and makes the code honest.

## Not urgent

Nothing is broken for the user today. This is code that lies about itself, which
matters mainly because the next person to touch overlays will believe the
comment.


## Summary of Changes

Shipped in `0218bf4`, OpenSpec change `modal-takes-the-whole-frame`.

The decision went to the first option: the full-frame modal is intended. It is
how specgetty has looked since the first commit, and compositing would change
how every modal reads for no reported complaint.

- `placeOverlay` lost the `background` parameter it discarded, and is now
  `modalFrame`, which is what it builds.
- The comment in `View()` that promised a stacking order is gone. The order
  there is precedence, not layering, and only one modal can be up at a time
  because whatever holds the keyboard refuses the keys that would raise another.
- The two test comments that repeated the same claim were corrected.

The decision is recorded in `openspec/specs/modal-presentation/spec.md` rather
than in a commit message, and pinned by
`TestAModalReplacesTheFrameRatherThanCoveringIt`. That test was checked against
a working composite implementation, which it fails.

## Correction to this report

The example given in conversation was wrong: pressing `a` with the picker open
does nothing, because the picker refuses every key it does not bind. The
observable case is pressing `a` from the change list, which leaves the
confirmation alone on a blank frame.
