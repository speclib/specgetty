---
# specgetty-lpg4
title: in specs.md and specs detail highlight mall keywords
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T21:53:22Z
updated_at: 2026-09-21T22:11:18Z
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
