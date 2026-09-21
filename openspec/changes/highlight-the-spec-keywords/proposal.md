## Why

Bean `specgetty-lpg4`. A spec is written in a small vocabulary that carries all
of its meaning, and specgetty draws almost none of it. `SHALL` is the most
common word in the corpus on this machine, 7997 uses across 544 specs, and it is
rendered as ordinary prose. Clause keywords are coloured on the detail card and
nowhere else: on the specs tab, `- **WHEN**` is merely bold, because the markdown
renderer sees two asterisks and not a keyword.

Measured on one real spec at card width, 534 drawn rows: 16 normative keywords
and 21 clause keywords per 100 rows. A screenful holds five or six of each. Dense
enough to be worth drawing, and never competing for the same position, because
clause keywords open a line and normative ones sit inside a sentence.

### The vocabulary is already decided, by a sibling project

`openspec.nvim` has a `spec-highlighting` capability that gives the same content
capture groups for an editor:

| element                                  | capture group  |
| ---------------------------------------- | -------------- |
| `SHALL`, `MUST`, `SHALL NOT`, `MUST NOT` | `@keyword`     |
| `WHEN` at the start of a condition line   | `@conditional` |
| `THEN` at the start of an assertion line  | `@property`    |
| `AND` at the start of a continuation line | `@operator`    |

Adopting it means a spec reads the same in neovim and in specgetty. It also
settles two questions this change would otherwise have to guess at. The normative
set is four words, not the ten of RFC 2119: openspec's own schema instruction
tells authors to "use SHALL/MUST for normative requirements (avoid should/may)",
and the corpus agrees, with 26 uses of `SHOULD` and 40 of `MAY` against 7997 of
`SHALL`. And clause keywords count only at the start of a line, which disposes of
the 40 mid-sentence uppercase `AND`s in the corpus without a special case.

`GIVEN` is the one word the nvim capability does not name, and it appears 735
times. It is a precondition, so it joins `WHEN`.

### The renderer cannot compose two styles, and already gets it wrong

`renderInlineMarkdown` applies bold, then code spans, then italics, each wrapping
the result of the last. Wrapping does not compose, because the inner style's
reset ends the outer one:

```
  "A scenario's clause is one source bullet, `- **WHEN** the user...`."
  → "...\x1b[36m- \x1b[1mWHEN\x1b[m the user...\x1b[m."
                              ^^^^^^ ends the cyan; the rest of the span is plain
```

That is live today in this project's own `spec-detail-view/spec.md`, and in 20
places across the corpus, every one of them a spec quoting the clause format. A
keyword pass layered on top would be a fourth wrapper over the same mechanism,
and it needs the one thing wrapping cannot give: leave a keyword alone when it is
inside a code span, which is where 44 of the corpus's keyword hits live, most of
them specs naming `SHALL` as a word rather than using it as one.

## What Changes

- Keywords are drawn wherever markdown is rendered: the specs tab, the detail
  card, a change's spec deltas, and proposals and designs. One renderer draws
  them all, so one rule covers them.

- The vocabulary and the styles:

| keyword                                | style             | why                        |
| --------------------------------------- | ----------------- | -------------------------- |
| `GIVEN`, `WHEN`                         | bold yellow       | what a clause keyword already looks like |
| `THEN`                                  | bold blue         | the assertion is the payload of a clause |
| `AND`                                   | yellow, not bold  | it continues the clause above it, so it is quieter than what it continues |
| `SHALL` `MUST` `SHALL NOT` `MUST NOT`   | bold magenta      | binding, and unlike the others it appears inside a sentence |

  Clause keywords are recognised at the start of a line only, with or without a
  bullet and with or without bold marks around them. Normative keywords are
  recognised anywhere in prose. Neither is recognised inside a code span, and
  both are uppercase only: lowercase `shall` appears zero times in the corpus, so
  the rule costs nothing and keeps ordinary English out of it.

- `renderInlineMarkdown` becomes a single tokenizing pass that emits each segment
  with exactly one style. This is what makes the keywords possible, and it fixes
  the bug above rather than building on it.

Not in scope:

- **`SHOULD`, `MAY`, `REQUIRED`, `RECOMMENDED`, `OPTIONAL`.** Left unhighlighted
  deliberately, matching the nvim capability and openspec's own advice to avoid
  them. Highlighting a word the format discourages would read as endorsement.
- **Splitting a heading into its prefix and its name.** The nvim capability draws
  `### Requirement:` and the name that follows it as two groups. specgetty draws
  a heading as one thing, and changing that is a separate question from
  keywords.

## Capabilities

### Modified Capabilities
- `specs-tab`: the inline renderer draws keywords, and one style per segment
  rather than styles wrapped around each other
- `spec-detail-view`: the card's clause keywords are drawn by role, and normative
  keywords inside clause text are drawn too

## Impact

- `src/ui/markdown.go` or `src/ui/ui.go`: `renderInlineMarkdown` rewritten as one
  pass over the line
- `src/ui/ui.go`: three clause styles where there was one, and a normative style
- `src/ui/specdetail.go`: the card asks the same renderer for its clause text
- `README.md`: what the colours mean

## Rollback

Its own commit. Reverting removes the keyword styles and returns the renderer to
wrapping styles around each other, which is the bug this change also fixes, so a
revert reintroduces it.
