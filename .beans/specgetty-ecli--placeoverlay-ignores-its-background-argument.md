---
# specgetty-ecli
title: placeOverlay ignores its background argument
status: in-progress
type: bug
priority: low
created_at: 2026-09-15T20:46:59Z
updated_at: 2026-09-18T09:21:44Z
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
