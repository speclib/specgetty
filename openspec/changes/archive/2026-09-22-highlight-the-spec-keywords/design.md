## Context

See proposal.md for why. What matters here is that the vocabulary was not
invented for this change and the renderer has to be rebuilt before it can carry
it.

`openspec.nvim/openspec/specs/spec-highlighting/spec.md` defines the capture
groups for the same content in an editor. This change is the same decision for a
terminal, so the two tools agree about what a keyword is.

`renderInlineMarkdown` is the one place that draws inline marks, and every
markdown surface goes through it: the specs tab, the detail card's prose and
clause text, a change's spec deltas, proposals and designs. One rewrite reaches
all of them, which is why the proposal does not list surfaces separately.

## Goals / Non-Goals

Goals: a spec's vocabulary visible as vocabulary; one renderer that composes;
the same answer as the editor about what a keyword is.

Non-Goals: a markdown renderer in general. Splitting headings into a prefix and a
name. Highlighting the words openspec tells authors to avoid. Configurable
colours.

## Decisions

### The vocabulary is adopted, not designed

Four normative words and three clause roles, from the nvim capability. The one
addition is `GIVEN`, which that capability does not name and which appears 735
times in the local corpus. It is a precondition, so it shares the condition role
with `WHEN` rather than earning a fourth.

Adopting rather than designing also settles the line-initial rule, which is what
keeps 40 mid-sentence uppercase `AND`s out of the highlighting without a special
case for them.

### The palette, and the two reuses in it

Every ANSI hue in this application already means something:

| hue      | what it means today                                  |
| -------- | ---------------------------------------------------- |
| red      | REMOVED, and a store that could not be followed      |
| green    | ADDED, and the active tab                            |
| yellow   | section headers, clause keywords, MODIFIED           |
| blue     | the panel name line                                  |
| magenta  | the store chip, as a background                      |
| cyan     | markdown headings, code spans, capability names      |

So:

| keyword                                | style             |
| --------------------------------------- | ----------------- |
| `GIVEN`, `WHEN`                         | bold yellow       |
| `THEN`                                  | bold blue         |
| `AND`                                   | yellow, not bold  |
| `SHALL` `SHALL NOT` `MUST` `MUST NOT`   | bold magenta      |

Three roles from two hues and a weight, which leaves green and red meaning only
"added" and "removed" where a spec delta shows both on the same screen as its
clauses. A third hue for `AND` would have had to be one of those.

Two deliberate reuses, named here so they do not read as accidents. Magenta is
the store chip's, but the chip is a background with black text in the header
line and a keyword is a foreground in the content pane. Blue is the panel name
line's, which appears once per screen and above the pane rather than inside it.

The alternative was the bright variants, 11, 12 and 13, which exist for exactly
this kind of emphasis. They are not taken because bright yellow on a light
background is close to invisible, and the tab bar already proves this application
is read on both.

### One pass, because wrapping does not compose

The renderer applies bold, then code, then italics, each wrapping the last. That
is broken, and demonstrably:

```
  mdBoldStyle.Render("prefix " + specKeywordStyle.Render("WHEN") + " suffix")
  → "\x1b[1mprefix \x1b[1;33mWHEN\x1b[m suffix\x1b[m"
                                    ^^^^^^ ends the bold; " suffix" is plain
```

Adding keywords as a fourth wrapper would multiply the same failure. So the
renderer becomes a single pass that scans the line once and emits each segment
with exactly one style.

Order of recognition inside that pass:

1. a code span, which binds tighter than everything else and whose contents are
   never looked at again. This is both what markdown says and what makes "never
   highlight a keyword inside a code span" fall out rather than be special-cased.
2. bold and italic marks
3. keywords, longest match first, so `SHALL NOT` is one keyword and never a drawn
   `SHALL` beside a plain `NOT`
4. everything else, as prose

An unclosed mark stays literal, which is the behaviour the current renderer has
and its tests already assert.

### Keywords are drawn before wrapping, and that is already true

The card styles a clause and then wraps it, a fix made in `open-a-spec-in-detail`
for spans that straddle a wrap. `ansi.Wrap` is ANSI-aware, so a `SHALL NOT` that
lands on a wrap boundary keeps its styling across the break. Nothing new is
needed here, but the ordering is load-bearing and a change to it would break this
quietly.

### The card keeps its own layout, and gains the renderer's keywords

A clause the card lays out already has its keyword on its own row, drawn by the
card. That keyword now carries its role's style instead of one style for all
four. The clause's text goes through the shared renderer, which is where a
`SHALL` inside it gets drawn.

A scenario the card could not lay out is drawn as prose, and that prose goes
through the same renderer, so its line-initial keywords are drawn even though the
card would not extract them. That is deliberate: `follow-the-spec-grammar` kept
the clause parser narrow, and 30 valid specs in the corpus write a shape it does
not extract. The parser answers "can I lay this out"; the renderer answers "is
this word a keyword". They are allowed to disagree, and the reader gains from it.

## Risks / Trade-offs

- The renderer is rewritten rather than extended, and everything markdown passes
  through it. The existing tests cover bold, code, italics and the unclosed
  cases; they are the floor, and the rewrite has to pass them unchanged.
- Four new styles on screens that already carry six. The density was measured
  before choosing: five or six keywords on a screenful, never two in the same
  position.
- Magenta and blue are reused. Both reuses are in a different region and a
  different treatment from their first use, but they are reuses.

## Migration Plan

None. Nothing is configurable, nothing is stored, and every document that renders
today renders tomorrow with more of it drawn.
