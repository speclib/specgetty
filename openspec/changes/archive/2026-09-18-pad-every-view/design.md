## Context

Three places could hold the inset, and they are not interchangeable.

```
A. on the panel                B. in each renderer         C. outside the border
   renderPanel pads               eleven views indent         the box itself is
   renderDetailPanel gets         themselves                  drawn narrower and
   two columns less                                           placed with a margin
```

B gives per-view control and costs eleven edits, each of which a future view has
to remember to repeat. It is also where the width arithmetic lives, and width
arithmetic in this codebase has drawn blood twice: the picker misalignment on
15 September, and the lipgloss v2 change from content width to total width.

C moves the gap outside the frame, which is not what the bean asks for.

A is one edit, every view now and later, and the arithmetic stays in the one
place that already owns it.

## Decisions

### The selection bar gets a gutter on both sides

This is the visible consequence of A and the only thing here that is a matter of
taste rather than correctness.

```
today                              with the inset
|### selected row ###########|     | ## selected row ######## |
 highlight runs border to border     one blank column each side
```

The second is exactly what the project picker draws today, and the picker is the
view that was reviewed and signed off after its alignment fix. So this does not
invent a look, it spreads the one already accepted.

Keeping the full-bleed highlight while insetting the text is possible but forces
option B: the padding would have to be painted per cell, inside the highlight,
in eleven places. Four times the work to preserve a detail that currently makes
the change list the odd one out.

### Column one aligns with everything else, it does not get a second gutter

"The table columns get the same left gutter" reads two ways.

```
A. the panel's inset serves as        B. every column gets its own,
   column one's gutter                   column one included

| name              tasks  specs |     |  name              tasks  specs |
| survive-bad-scan  6/6    1     |     |  survive-bad-scan  6/6    1     |
  aligned with the tab bar                two in, while the tab bar,
  and the header above                    header and markdown sit at one
```

A is taken. One vertical line runs down the whole panel, and the table reads as
panel content rather than as a nested block. B is defensible but it breaks the
alignment with everything above the table, which is the alignment this change
exists to create.

### The nav bar indents its text and keeps its fill

The nav bar is a full-width painted strip, not a box.

```
C. text indented, paint full bleed      D. the whole strip inset

### q quit  jk/up-down navigate ####     ### q quit  jk navigate ###
^ paint from column zero                ^ unpainted notch at each end
```

C is taken. The strip's job is to mark the bottom edge of the screen, and D
leaves a one-column hole at each end of it, which reads as a rendering fault
rather than as spacing.

### The header gives up its own padding rather than the panel giving way

`Padding(1, 1)` on the header block becomes `Padding(1, 0)`. The vertical half
is what gives the header its breathing room above and below, and the bean was
narrowed to left and right precisely because that spacing is already right.

This was not predicted from reading the code. It showed up the first time the
inset was rendered: the header sat two columns in and everything under it one.

## Risks

### The three copies of the content width

`renderFrame`, `docRegion` and `recalcLayout` each derive the panel's content
width by writing `m.width - 2`. Both `docRegion` and `specsSplit` carry comments
warning that these must agree or the viewport's reported scroll position stops
matching the screen. Neither comment stopped the duplication.

Moving the inset means moving all three together. This change therefore depends
on `unify-content-width` landing first, which collapses them into one method.
Without that, this change adds a fourth place to get the number wrong.

### The assertion that will not catch it

The picker misalignment is the cautionary tale and it is directly relevant. A
test asserting every line had the same width passed both with and without the
bug, because lipgloss pads a wrapped remainder back out to full width. Only the
line count caught it.

So the assertions that bite here are: the frame is still exactly `m.height`
lines with each view up, a known-width document wraps at the expected column
asserted on its text rather than on its rendered width, and `docRegion` agrees
with what `renderFrame` passes down.

### Rollback

Its own commit on top of v0.5.0. `git revert` takes the look out and leaves the
`unify-content-width` deduplication in, which is the point of shipping them
separately. A baseline capture of twenty-seven views at three terminal widths
was taken before any of this landed and sits in the session scratchpad as
`views-before-v0.5.0.txt`, so the diff can be read rather than squinted at.
