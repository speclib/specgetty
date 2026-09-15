---
# specgetty-6t74
title: project picker highlighted row is misaligned
status: completed
type: bug
priority: normal
created_at: 2026-09-15T17:47:19Z
updated_at: 2026-09-15T17:50:16Z
---

The selected row in the project picker does not line up with the other rows.

## Cause

renderPicker builds the overlay with:

    lipgloss.NewStyle().Border(...).Padding(0, 1).Width(inner).Render(content)

and passes the same `inner` to renderTable. But lipgloss Style.Width includes
horizontal padding and excludes the border, so the usable content area is
inner-2, not inner. Every table row is therefore 2 columns too wide and lipgloss
wraps it.

On an unstyled row the two wrapped columns are spaces and invisible. On the
highlighted row they carry the selected background colour onto the next line,
so only that row looks misaligned.

Measured: box Width(88) with Padding(0,1) and a border renders 90 columns, so
content is 86. renderTable was given 88. The overlay rendered 8 lines where the
layout calls for 6, which is the wrap.

## Summary of Changes

Two border-width defects, both found by measuring rather than by eye.

### The reported one: picker rows

`renderPicker` built the overlay with `Padding(0, 1).Width(inner)` and passed
the same `inner` to `renderTable`. lipgloss `Style.Width` covers content plus
horizontal padding and excludes the border, so the usable content area was
`inner-2`. Every row overflowed by two columns and lipgloss wrapped it.

On an unstyled row the wrapped remainder is spaces and invisible. On the
highlighted row it carries the selected background colour, which is why only
that row looked wrong.

Fixed by separating the two numbers: `styleWidth = width-2` for the box,
`inner = styleWidth-2` for the content.

### The one found next to it: panel top border

`renderPanel` assembles its top border by hand and was one column short:
`width-len(title)-2` where the arithmetic wants `-1`. Three columns go to the
two corners and the segment before the title, not two. The top-right corner sat
one column inside the box edge.

Pre-existing, unchanged since before the picker work (verified against commit
91d777d). Fixed here because it is the same class of defect and one character.

### Regression tests

In `src/ui/layout_test.go`. Both were checked by reintroducing each bug and
confirming the test fails, because a regression test that passes either way is
worthless.

Worth knowing: "all lines have the same width" does NOT catch an overflowing
row, because lipgloss pads the wrapped remainder so the widths still match. The
detector is the line COUNT. That is recorded in the test file so nobody leans on
the weaker check.

Styles are also stripped when no terminal is attached, which hides exactly this
class of defect, so the tests force a colour profile.
