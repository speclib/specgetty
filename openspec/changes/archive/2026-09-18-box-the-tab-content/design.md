## Context

The mockup shows one of four content areas. The other three each ask the same
question differently, and the answers are what most of this change is.

```
changes   table, plus a search prompt when a filter is on
specs     two panes side by side
config    a filename, a blank line, then a scrolling document
L2        a change header, a sub-tab row, then a scrolling document
```

## Decisions

### What is chrome and what is content

The rule that settles all four: **a line that names the content is chrome; a
line that reports on the content is content.**

```
chrome, above the border          content, inside it
──────────────────────           ──────────────────
the tab bar                       the table
the config tab's filename         the search prompt and its count
the change header of an open      the artifact document
  change, and its sub-tab row
```

The filename says what you are looking at, the same job the tab chips do. The
search prompt says how much of it you are being shown, which is only meaningful
against the rows beside it. Drawing the filename inside the box would make the
box contain two unrelated things; drawing the prompt outside it would separate a
count from what it counts.

### The specs tab keeps one border for now

It is two panes, so it could take one border around both or one each. This
change gives it one, which is exactly what it looks like today plus an edge:

```
╭──────────────────────────────────────╮
│ aggregate-task-stats   # aggregate…  │
│ archive-from-ui                      │
│ change-list-view       ## Purpose    │
╰──────────────────────────────────────╯
  the split stays implied by whitespace
```

Splitting it into two borders is `light-the-focused-pane`, because the moment
there are two borders they need a colour rule, and that rule is the focus
signal. Shipping the split without the colouring would mean choosing a colour
we know to be wrong. One border here is a coherent resting state, not a
placeholder.

### The panel border is left alone

An earlier reading of this had the panel border going permanently dim, on the
grounds that once the content has a border of its own the outer one is
redundant. That reasoning assumed the two borders would compete.

They do not, because the rule that governs them is containment: a border is lit
when the keyboard is inside it. The panel encloses the content, so it stays lit
whenever the content has the keyboard, which is exactly what it does today. The
finer signal goes underneath without contradicting it.

Nothing in this change touches the panel border. The rule is written down in
`light-the-focused-pane`, where it first has two borders to govern.

## Risks

### Two rows is the price, and it is paid where it hurts most

```
terminal   content rows   with the border   change rows
  60x20         12              10           11 -> 9
  92x30         22              20           21 -> 19
  92x50         42              40           41 -> 39
```

A border costs a fixed two rows, so it is a fifth of the budget at the minimum
size and a twentieth at a large one. This is not recoverable by arithmetic; it
is what the change buys and what it costs.

### Four columns, and a table that already drops columns

The content width goes from `m.width - 4` to `m.width - 8`. The change table
drops its rightmost column when the flexible column would fall under twelve,
which `widen-table-column-gaps` measured at a 28-column table. At the 60-column
minimum the table now gets 52, so nothing is lost that was not already lost.
Worth a test rather than the arithmetic above, because the arithmetic above has
been wrong before.

### The assertion that will not catch it

The picker misalignment of 15 September is the standing example: a test that
every line has the same width passed both with and without the bug, because
lipgloss pads a wrapped remainder back out to full width. Only the line count
caught it. A nested box is exactly that shape of change, so the assertions here
are line counts and wrap columns, not widths.

### Rollback

Its own commit above `unify-focus-state`. A before-and-after capture of every
view at three terminal widths is taken in the same session and diffed, so the
review reads a diff rather than comparing screenshots from memory.
