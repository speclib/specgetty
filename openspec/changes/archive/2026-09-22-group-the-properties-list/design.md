## Context

See proposal.md for why. What matters here is that the mechanism already exists
twice and the look is the only thing being invented.

The change list solved the non-selectable header two changes ago, and its
comment states the rule this change follows:

> The cursor never indexes these. It indexes changes, because the selection key,
> the search filter, `g`, `G` and all three actions already work that way and
> none of them should have to know a header exists. Only the arithmetic below
> maps between the two.

`propSections` already returns the rows in the order the groups want them: the
configuration, then the schemas, then the store. Only the store has to move, and
only in the list, not in what any row reports.

## Goals / Non-Goals

Goals: a list that says what its rows are; one arithmetic for the three lists
that have it; a width chosen on evidence rather than on a rule whose reason has
expired.

Non-Goals: the problems row, counts on headers, collapsible groups, a second
level of nesting, or any change to what a row's content reports.

## Decisions

### The look

Uppercase headers, rows indented by two, a blank line between groups, no counts:

```
 ╭──────────────────╮ ╭───────────────────────────────────────────────╮
 │ PROJECT          │ │ # /stores/nivis-tunnel/openspec/project.md    │
 │   config         │ │                                               │
 │ ▌ store          │ │ store: nivis-tunnel                           │
 │                  │ │ root:  /stores/nivis-tunnel                   │
 │ SCHEMAS          │ │ git:   clean, 3 commits ahead                 │
 │   spec-driven    │ │                                               │
 │   tinychange     │ │                                               │
 ╰──────────────────╯ ╰───────────────────────────────────────────────╯
```

Uppercase because that is what a group header already reads as in this
application: the change list draws `ACTIVE` and `ARCHIVED`. Rows stay lowercase
because a schema is named `spec-driven` by its own file, and `Config` beside
`spec-driven` would be two conventions in a list of five rows.

No counts, though the change list's headers carry them. A count earns its place
on a list of 47 changes under two headers. Across the 33 OpenSpec projects on
this machine, 25 use one schema and 8 use two, so every count this list could
draw would be `(1)` or `(2)`.

### The width is a floor, not a share

Three policies, measured against the panel:

| terminal | sized to labels | floor of 20 | three tenths |
| -------- | --------------- | ----------- | ------------ |
|       60 |              17 |          18 |           16 |
|      100 |              17 |          20 |           28 |
|      200 |              17 |          20 |           58 |

A share is wrong here for the reason the original rule was right about: the
labels are short and fixed, so at 200 columns a share hands 58 columns to
thirteen characters. A floor keeps the list the same width on every terminal
that can hold it, which also means the content grows with the terminal and the
eye does not have to find the divider again after a resize.

The floor gives way to the existing one-third cap on a narrow terminal, so the
two halves cannot collide.

### The row named `project` is named `config`

It shows the configuration file and the path it was read from. Under a header
that already says `PROJECT`, a row named `project` says nothing, and the group it
sits in now supplies the context the old label was carrying.

### An empty group keeps its header

`countSchemaUsage` adds the project default even when nothing uses it, so a
project with a recorded default always has at least one schema row. A project
with no default and no changes carrying one has none. The change list already
answered this case and gave the reason: a group keeps its header even when it
holds nothing, because "nothing" is an answer and an absent header cannot be
told apart from one that was filtered away.

### One arithmetic, extracted before it is used

Three lists, one shape: the cursor indexes items, the pane counts drawn lines,
and an item may occupy more than one line. The thinnest interface between them
is a slice of item indices, one entry per drawn line, with chrome marked:

```
  []int, one entry per drawn line, -1 where the line is a header or a spacer

  span(item)          first and last drawn line of an item
  fit(from, rows)     how many items fit in a pane of that many rows
  offsetFor(item, h)  the scroll offset that brings the whole item into view
```

Checked against each site before being written:

- `offsetFor` is `max(0, last-h+1)`, clamped down to `first` when the item is
  taller than the pane. For the change list, where every item is one line, that
  reduces to today's `at - bodyHeight + 1`.
- `fit` is `nodesInRows`, and reduces to the row count when items are one line
  each, which an existing test already asserts.
- The change list's spacer-skip stays where it is. It needs to know that a line
  is a spacer rather than which item it belongs to, which is the one thing this
  interface deliberately does not carry.

The two shipped sites move first, proven by the tests they already have, and the
new list is built on the helper. The other order leaves a third copy of the
arithmetic in the tree for the length of the change, and copies outlive changes.

### The problems row gets a seat and nothing else

`specgetty-p44t` will report the specs openspec cannot read, per project. That is
a third row under `PROJECT`, beside `config` and `store`. Grouping the list now
means that row is an addition rather than a redesign, which is the whole reason
to name the group `PROJECT` rather than `CONFIGURATION`.

## Risks / Trade-offs

- Two shipped renderers are refactored to reach one helper. Both have tests that
  assert their arithmetic directly, including the change list's off-by-one cases
  and the outline's wrapped-node scrolling, so the refactor is verifiable rather
  than hopeful. It is still the largest part of this change by risk.
- The list takes three columns from the content pane. Measured over the widths
  that matter, those columns never decide whether the longest line wraps.
- A list of seven drawn lines in a pane of twelve still leaves air. The specs tab
  has the same air for a small project, so this is consistent rather than new,
  and the problems row will take one of the empty lines.

## Migration Plan

None. Every row reports what it reported, in the same order, under a header that
did not exist.
