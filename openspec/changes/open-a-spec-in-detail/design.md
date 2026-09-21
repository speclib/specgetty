## Context

See proposal.md for why. What matters here is the shape of the code this lands
in.

The specs tab is a split: `specsSplit` divides the panel, `renderSpecsTab` draws
the list half, and the content half is `m.docViewport.View()` fed by
`currentDoc`. `docRegion` computes the content width from the same `specsSplit`,
and its doc comment says why the two must agree: a document wrapped to a
different width than the pane it is drawn in re-wraps inside the box and the
reported position stops matching the screen.

Navigation is a linear depth counter, `levelProject` and `levelChange`. Six
places branch on it together with a tab name: `docRegion`, `docActive`,
`listPage`, `moveListCursor`, `gotoListEnd` and `splitTab`, plus `renderNavBar`.

Nothing in the codebase parses markdown structure today. `renderMarkdownLines`
styles one line at a time and reports where each source line landed.

The corpus, measured across the 24 live specs of this project:

```
                       count     median     p90     max
  requirements           140      36 col   53 col  63 col
  scenarios              428      28 col   43 col  63 col
  nodes per spec          24 median, 57 max (change-list-view)
  headings outside the four known shapes: 0
  requirements with no text above the first scenario: 0
  requirements with no scenario: 0
  scenarios with no keyword bullet: 0
  keywords in use: GIVEN, WHEN, THEN, AND
```

## Goals / Non-Goals

Goals: read one spec by its structure; keep every existing surface working
unchanged; leave the parser able to grow a delta level later.

Non-Goals: the specs of an open change, which are deltas and need a level this
change does not build. Collapsing and expanding the outline. Any editing.

## Decisions

### A third level, not a deeper tree in the existing pane

The outline replaces the spec list rather than nesting inside it. The
alternative, one collapsible tree of spec, requirement and scenario in the
existing left half, needs no new level, but the overview and the detail then
compete for one column and the at-a-glance list of 24 specs is lost. A level
also reuses a rule the application already has: `enter` descends, `esc` ascends.

The cost is honest: a third arm in six switches and the nav bar. Add the arms
rather than inventing a surface abstraction in the same change. The arms are
where the duplication is already visible, and a rewrite of the layout arithmetic
under a new abstraction is its own change with its own risk.

### The outline wraps, and colour carries the level

At three tenths on an 80-column terminal the outline gets 22 columns, and:

```
  outline width    requirement titles fit    scenario titles fit
    19 (70 term)       25/140  (18%)             80/428  (19%)
    22 (80 term)       31/140  (22%)            131/428  (31%)
    28 (100 term)      43/140  (31%)            220/428  (51%)
    34 (120 term)      61/140  (44%)            328/428  (77%)
```

Four rows in five would be a stub at the common width. Truncating and letting
the card carry the full title was considered and rejected: an outline whose rows
cannot be read is a worse spec list than the flat one that ships today.

Wrapping brings a consequence. With requirements flush left and scenarios
indented two, a wrapped requirement row sits at the indent a scenario title uses,
so indentation can no longer carry the level. Colour does: requirements in
`sectionHeaderStyle`, the bold yellow already used for spec names in
`renderChangeSpecs`, scenarios in the normal style. This is the convention the
change list uses to separate its group headers from its rows.

A node then occupies one to three rows. That mapping already exists in the
codebase for a different cursor: `sourceLine{rowStart, rowEnd}` with
`highlightCursorLine` and `scrollCursorIntoView`. The outline gets the same
shape rather than a new one.

### Four tenths, not three

Wrapping makes any width work, so the split is a trade between outline rows and
card columns:

```
  split  term  outline  card  outline rows (median/max)
   30%     80       22    49        43 / 111
   40%     80       30    41        34 /  88
   40%    100       38    53        29 /  73
   40%    120       46    65        26 /  64
   50%     80       38    33        29 /  73   card too narrow
```

Four tenths. Five tenths buys five rows of outline and costs the card eight
columns, which it cannot spare once padding and the clause indent come out.

The split is its own function with two callers, `renderSpecDetail` and
`docRegion`, for the reason `specsSplit`'s comment already gives.

### What each node's card holds

Purpose shows its prose. A requirement shows the text between its heading and
its first scenario, and nothing more. A scenario shows its card.

Stacking a requirement's scenarios under its text was considered, since it would
give one node from which the whole requirement reads. It was rejected because it
makes one node in three unbounded in height while the other two are cards, and
the scenarios are one keypress away in the outline either way.

### A scenario's card is a re-layout

The source is one bullet per clause. The card puts the keyword on a row of its
own and indents the clause under it, so this is a second renderer over the parsed
clauses, not a style table over `renderMarkdownLines`. Inline styling, bold and
backticks, is shared with the markdown renderer, which is why the backtick
highlight lands in `renderInlineMarkdown` and reaches every document.

Padding is a formula rather than a set of width tiers: tiers reflow the card at a
threshold, a formula degrades smoothly, and padding is the only decoration with
columns to give back.

### Parsed when opened, held while open

`View` runs on every message, spinner ticks and watcher events included, so
parsing per frame would do real work for nothing. Parse on `enter`, hold the tree
for the open spec, drop it on `esc`. Nothing is parsed in a session that never
opens a spec, which is what parse on demand was asked for.

### A spec that does not fit does not open

The strict rule is safe on this corpus and this corpus is not the point:
specgetty reads other people's projects. So the failure has to be visible and
harmless. `enter` reports on the nav bar through `statusMsg`, which takes the
bar's row and is cleared by the next keystroke. The model comment on that field
already argues the case: a modal is the wrong weight for something instant and
harmless.

The markdown view on the tab is untouched, so a spec that cannot be outlined is
still fully readable. That is what makes the strict rule affordable.

The parse fails when the file yields no requirements, or when a requirement has
no scenarios. Those are the two conditions under which the outline would have
nothing to navigate.

### A page is a screenful of rows, walked over nodes

`document-viewer` already requires that a page move by the number of rows the
surface is showing. Over an outline whose nodes are not all one row tall, moving
the cursor by that many nodes would overshoot by roughly the average node height.

So the page walks: accumulate node heights forward until one pane of rows is
covered, and land the cursor there. This satisfies the existing requirement
literally, needs no row-based list scrolling, and reduces to today's behaviour
when every node is one row tall. `document-viewer` gains a scenario saying so,
because the existing wording under-determines it and the next variable-height
surface should not have to guess.

### The cursor is kept by node, not by index

The watcher rewrites files while they are read. The outline cursor is held by the
node's path through the tree, so an edit above it does not move the reader, and a
node that disappears lands the cursor in range rather than out of it. The change
list already does this with `selectedKey`.

## Risks / Trade-offs

- A third arm in six switches plus the nav bar is duplication the code already
  carries once. → Add the arms, and let a surface abstraction be its own change
  when a fourth level asks for it.
- The strict parse locks out a project whose specs are shaped differently. → The
  failure is a one-line report and the markdown view is unchanged, so nothing
  becomes unreadable. If real projects turn out to be looser, the rule loosens
  without touching the rest.
- A 57-node spec wraps to 88 rows at 80 columns, so the outline scrolls. → It is
  a list with the ordinary list keys, which is what every other list here does.
- Parsing at `enter` puts work on a keypress. → The largest spec in this corpus
  is 12.7 KB and the parse is a single pass over its lines.

## Migration Plan

None. The change is additive: without pressing `enter` on the specs tab, nothing
behaves differently except that backticked spans are now highlighted.
