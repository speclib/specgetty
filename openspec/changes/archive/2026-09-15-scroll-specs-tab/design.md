## Context

The archived design of `scroll-markdown-documents` recorded the answer to this
question rather than leaving it open:

> When the specs tab content pane becomes scrollable, `tab` toggles focus
> between the spec list and the content, and the focused pane gets the active
> border treatment `renderPanel` already applies.

The first half holds. The second half does not survive contact with the code.
`renderPanel` applies its active border to a whole panel, and since
`project-picker` removed the split layout there is only one panel. The two
halves of the specs tab live inside it, so there is no border to colour.

## Decisions

### Focus is shown by which cursor is lit

The spec list already highlights its cursor with `selectedStyle`. Focus is shown
by that highlight moving between the halves rather than by a border:

```
  list focused                      content focused
  ┌───────────┬──────────────┐      ┌───────────┬──────────────┐
  │ >spec-a<  │ # Spec A     │      │  spec-a   │ # Spec A     │
  │  spec-b   │ content...   │      │  spec-b   │ content...   │
  └───────────┴──────────────┘      └───────────┴──────────────┘
     ^ lit cursor                      ^ dimmed      title shows 42%
```

With the content focused the list cursor dims to show it is still the selected
spec but no longer the thing the keys move, and the panel title carries the
scroll percentage, which only a focused document produces. Two signals, no new
chrome, no column of width spent on a border.

### tab cycles rather than toggles

`tab` toggled the log panel before this change. On the specs tab it now cycles
spec list, spec content, and the log when the log is open. Making it toggle only
the specs halves would strand the log panel on this one tab.

Everywhere else `tab` keeps its current meaning.

### The document is keyed by spec name

`currentDocument` gains a specs case keyed by project and spec name, so moving
the spec cursor is moving to a different document and starts at the top, by the
same rule that already governs artifact sub-tabs. No separate reset handler.

Note that the cursor can only move while the list is focused, and the content
can only scroll while the content is focused, so the two never fight.

### The content half wraps to its own width

`docRegion` returns the full panel width today. The specs content half is
roughly 70% of it, so the specs case has to compute the same `listWidth` split
that `renderSpecsTab` uses. If the two disagree the rows re-wrap in the box and
the reported position stops matching the screen, which is the failure this whole
line of work exists to avoid. The split is therefore computed in one place and
used by both.

## Risks / Trade-offs

- **`tab` means something new on one tab.** It cycles three ways there and two
  ways elsewhere. The alternative was a second key for focus, spending a
  binding on something `tab` already means everywhere else in the tool.
- **A dimmed cursor is a quieter signal than a border.** It is paired with the
  percentage in the title, which appears only when the content is focused, so
  there are two independent cues.

## Open Questions

- The change list and the picker could take the same focus treatment if their
  panes ever grow a second half. Nothing needs it today.
