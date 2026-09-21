---
# specgetty-5awz
title: ansi.Wrap puts a break character one column past the pane
status: todo
type: bug
priority: normal
created_at: 2026-09-21T22:27:55Z
updated_at: 2026-09-21T22:27:55Z
---

`ansi.Wrap` emits a break character even when it does not fit, so a wrapped row
comes out one column wider than the width it was given. The box that draws it
then re-wraps that row, and the reader gets an orphan row that nothing asked for.

Found while shipping `highlight-the-spec-keywords`, where a test of mine
asserted no row exceeds its pane and failed for a reason that change did not
cause. Confirmed identical before and after it.

## The rule

A break character is placed past the limit rather than moved to the next row:

```go
ansi.Wrap("aaaa bbbb cccc ddddd, eeee", 20, "/,;:")
//   row 0: "aaaa bbbb cccc ddddd," width 21, limit 20
//   row 1: "eeee"
```

Two sources, and the second survives fixing the first:

- the breakpoints this application passes, `wrapBreakpoints = "/,;:"` in
  `src/ui/ui.go`, used by `wrapStyled`
- the hyphen, which `ansi.Wrap` breaks on natively whatever is passed

## What the reader sees

The over-wide row does not reach the terminal: the box re-wraps it. So the cost
is not a lost character, it is a row that splits in the wrong place.

```
  the renderer produced            what the box drew
  ---------------------            -----------------
  width 53:                        ╭──────────────────────────────────╮
  "  - THEN the content SHALL      │   - THEN the content SHALL say    │
   say that it is being read,"     │ read,                             │
  width 39:                        │ and the interface SHALL stay      │
  "and the interface SHALL         ╰──────────────────────────────────╯
   stay responsive"
                                     ^ an orphan row, and the
                                       continuation loses its place
```

## The part that is not cosmetic

`renderMarkdownLines` records a `rowStart` and `rowEnd` per source line, and the
box then draws a different number of rows than it counted. With one overflowing
line in a document, every row index below it is short by one.

Verified: content handed to the box as 4 rows was drawn as 5. Not verified, and
worth checking as part of the fix: whether the task checkbox cursor, which reads
that mapping, highlights the wrong rows below an overflowing line. The highlight
is applied before the box re-wraps, so it may follow its text and be fine.

## How often

Over the 606 live specs on this machine, rendered at a range of pane widths:

| pane width | rows drawn | rows one column over |
| ---------- | ---------- | -------------------- |
|         40 |    100,610 |            357 (0.35%) |
|         52 |     84,735 |            251 (0.30%) |
|         60 |     78,727 |            192 (0.24%) |
|         76 |     68,436 |            152 (0.22%) |
|        100 |     60,169 |             21 (0.03%) |

Worst where it matters most, in a narrow pane.

## Fixes, measured

Dropping the breakpoints removes 96% of it:

| breakpoints | rows   | overflowing |
| ----------- | ------ | ----------- |
| `"/,;:"`    | 84,735 |         251 |
| `""`        | 84,800 |           9 |

The cost is 65 extra rows across the whole corpus, and paths break at the column
rather than after a slash:

```
  "/,;:"  ["root: /home/pim/gh.nivis-" "project/nivis-openspec-stores/" "nivis"]
  ""      ["root: /home/pim/gh.nivis-" "project/nivis-openspec-" "stores/nivis"]
```

Barely a difference on real content, which suggests the breakpoints are not
earning what they cost.

The remaining 9 are all hyphens and need something else:

```
  "- WHEN text contains a -----BEGIN … PRIVATE KEY-----"     54 in a 52 pane
  "- WHEN a spec runs beans update --help, beans list "      53
  "- WHEN the first line of the input is not exactly ---"    55
```

Three ways to finish it:

1. A post-check in `wrapStyled`: after wrapping, any row wider than the width
   has its trailing break character moved to the next row. Local, precise, and
   covers both causes.
2. Wrap at `width` and let a hard truncate catch the remainder. Loses a
   character, which is worse than an odd row.
3. Report it upstream to `charmbracelet/x/ansi`. Worth doing regardless of what
   is done here, since the behaviour looks unintended: the function is given a
   width and returns something wider.

I would take 1 and also do 3.

## Where

- `src/ui/ui.go`: `wrapStyled`, and `wrapBreakpoints`
- `src/ui/ui.go`: `renderMarkdownLines`, which is what records the row mapping
- the same `wrapStyled` draws YAML, so the properties tab has it too
