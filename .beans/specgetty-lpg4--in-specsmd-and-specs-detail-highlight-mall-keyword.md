---
# specgetty-lpg4
title: in specs.md and specs detail highlight mall keywords
status: completed
type: task
priority: normal
created_at: 2026-09-21T21:53:22Z
updated_at: 2026-09-21T22:22:10Z
---

- WHEN
- THEN
- GIVEN
- SHALL
- SHALL NOT
etc...


## OpenSpec change

`highlight-the-spec-keywords`, validated strict. 28 tasks.

Decisions taken from this bean's exploration:

- The vocabulary is adopted from `openspec.nvim`'s `spec-highlighting`
  capability, so a spec reads the same in neovim and in specgetty. Four binding
  words (`SHALL`, `SHALL NOT`, `MUST`, `MUST NOT`) and three clause roles.
  `GIVEN` is the one addition: that capability does not name it and it appears
  735 times locally, so it joins `WHEN` as a condition.
- `SHOULD`, `MAY`, `REQUIRED`, `RECOMMENDED` and `OPTIONAL` stay unhighlighted.
  openspec's own schema instruction tells authors to avoid them, and the corpus
  agrees: 26 and 40 uses against 7997 of `SHALL`.
- Clause keywords at the start of a line only, which disposes of the 40
  mid-sentence uppercase `AND`s without a special case. Normative keywords
  anywhere. Upper case only, which costs nothing: lower-case `shall` appears zero
  times in 544 specs.
- Styles: condition bold yellow, assertion bold blue, continuation yellow
  unbolded, binding bold magenta. Three roles from two hues and a weight, so
  green and red keep meaning only added and removed where a spec delta shows both
  beside its clauses. Magenta and blue are reuses, of the store chip and the
  panel name line, both in a different region and treatment; the design names
  them so they do not read as accidents.
- Every markdown surface, since one renderer draws them all.

The change also fixes a live bug it cannot be built on top of.
`renderInlineMarkdown` wraps one style around another, and the inner style's
reset ends the outer one:

```
  "...one source bullet, `- **WHEN** the user...`."
  → "...\x1b[36m- \x1b[1mWHEN\x1b[m the user...\x1b[m."
                              ^ ends the cyan, so the rest of the span is plain
```

Reachable in 20 places in the corpus, including this project's own
`spec-detail-view/spec.md` line 119. The renderer becomes a single pass that
emits each segment with exactly one style, which is also what makes "never
highlight a keyword inside a code span" fall out rather than be special-cased.

Measured before choosing: 16 normative and 21 clause keywords per 100 drawn rows
of a real spec, so five or six of each on a screenful, and the two never compete
for a position because one opens a line and the other sits in a sentence.


## Summary of Changes

Shipped as OpenSpec change `highlight-the-spec-keywords`, commit `0060205`.
28 tasks.

The vocabulary a spec is written in is now drawn wherever markdown is rendered:
the specs tab, the detail card, a change's spec deltas, and the proposal and
design documents. `GIVEN` and `WHEN` open a condition, `THEN` the assertion,
`AND` continues the clause above it, and `SHALL`, `SHALL NOT`, `MUST` and
`MUST NOT` bind. The words and the roles are openspec.nvim's, so a spec reads
the same in the editor and here.

Three roles from two hues and a weight rather than three hues, because green and
red mean added and removed in a spec delta, which shows both on the same screen
as its clauses. The two reuses, blue from the panel name line and magenta from
the store chip, are named in the design so they do not later read as accidents.

The change also fixed the bug it could not be built on top of. The renderer
applied bold, then code spans, then italics, each wrapping the last, and the end
of an inner style ends the outer one, so a spec quoting `- **WHEN** ...` was
drawn half in colour and half in plain text. Reachable in 20 places in the local
corpus, including this project's own `spec-detail-view/spec.md`. One line is now
drawn in a single pass with exactly one style per segment, which is also what
makes "never draw a keyword inside a code span" fall out rather than be a case
of its own.

Three things the tasks did not anticipate, each recorded in `tasks.md`:

- Task 4.6 needed a parser change to be true at all. The card joined a
  scenario's consecutive source lines into one paragraph, as markdown says to,
  so in a scenario written without bullets only the first keyword sat at a line
  start. A prose line opening with a clause keyword now starts a paragraph of
  its own. The lines stay prose and the clause parser stays as narrow as
  `follow-the-spec-grammar` left it; they just keep the boundaries the author
  wrote. The 30 valid specs in the corpus that write that shape get their
  structure back visually without the parser widening.
- The first attempt at revert check 5.1 was a no-op. Swapping the order of the
  code-span and bold branches changes nothing, because a position is either a
  backtick or an asterisk. What carries the fix is consuming a span whole, so
  the revert now reproduces the original nesting and fails all three tests.
- Task 5.4 asked for an assertion that would have failed for a reason this
  change did not cause. `ansi.Wrap` lets a breakpoint character sit one column
  past the limit, identically before and after this change. The assertion this
  change owns is that drawing a keyword moves no column: a styled row measures
  exactly as its text does, checked over every live spec.

Coverage: ui 89.9% (floor raised 88.3 to 89.6), total 90.8% (89.8 to 90.5). Gate
green. Both revert checks bit, with the build confirmed first.
