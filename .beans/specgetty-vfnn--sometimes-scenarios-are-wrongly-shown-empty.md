---
# specgetty-vfnn
title: sometimes scenarios are wrongly shown empty
status: completed
type: bug
priority: high
created_at: 2026-09-21T20:19:39Z
updated_at: 2026-09-21T21:18:26Z
---

check /home/pim/mipnix/openspec/specs/airplane-mode/spec.md

the scenario's are rendered empty. Maybe caused by `GIVEN`


## OpenSpec change

`follow-the-spec-grammar`, validated strict. 38 tasks.

## Diagnosis

Not `GIVEN` itself. `clauseOf` requires a `- ` bullet, and matches keywords
case-sensitively. The file writes clauses bare:

```
GIVEN airplane mode is off
WHEN the user clicks the airplane mode toggle
THEN all radios SHALL be blocked via `rfkill block all`
```

The defect is that a scenario is the only node kind with no fallback. Every
other node does `cur.body += line`; a scenario ends its branch with a bare
`continue`, so any line the clause matcher does not recognise is dropped.

Measured over the 544 live specs under /home/pim: 27 specs open with at least
one blank card, 254 cards in total, and every one of the 254 is dropped content.
None is a scenario that is genuinely empty in the file. Ten of the 27 pass
`openspec validate --strict`.

## Why widening the matcher was the wrong fix

OpenSpec has an official grammar, in code at
`openspec-1.10.0/lib/openspec/dist/core/parsers/`. A scenario in its data model
is `{ rawText }`. There is no GIVEN/WHEN/THEN in it at all, and
`openspec show <spec> --json` confirms it. The `- **WHEN**` form appears in the
template, in the schema instruction and in a validator error message: three
pieces of guidance and no grammar. So the shapes being dropped are not
malformed, and no validator will ever push authors toward one shape.

Meanwhile specgetty was lax where OpenSpec is strict: 136 of the 544 live specs
are invalid by the official grammar, 102 because a main spec carries a delta
header, and specgetty opened them and showed part of what they held.

The change therefore does both: strict to the grammar where a rule exists,
lossless where there is only a habit.

## What the change does

- reads the grammar OpenSpec defines, transcribed from its parser: the
  `## Requirements` section, `### Requirement:` inside it, any non-fenced
  `####` heading with content, fences masked, a delta header in a main spec an
  error
- carries a scenario's content in full, as ordered parts, so nothing is dropped
  whatever shape it is in; the clause layout becomes presentation
- descends into a report naming every reason and its line, instead of turning
  back with one line on the nav bar
- adds a fixture corpus, one file per shape found in the wild. The test that
  should have caught this walked only this project's own 24 specs, every one
  written in one shape by one author

Deferred to `specgetty-p44t`: a per-project problems report and a mechanical
repair, which would reach 130 of the 136 invalid files.


## Summary of Changes

Shipped as OpenSpec change `follow-the-spec-grammar`, commit `8706633`. 38 tasks.

Both halves of the diagnosis, since either alone leaves the bug alive somewhere.

**Nothing is dropped.** A scenario now carries its content as ordered parts, each
a clause or a paragraph, so no line is lost for being written in a shape the
clause matcher does not know. Ordered rather than grouped, because the two
interleave: the corpus has horizontal rules and `**Rationale**:` paragraphs
sitting between clauses, and grouping would reorder a behaviour contract.

**Strict where OpenSpec has a rule.** The grammar is transcribed from
`openspec-1.10.0/dist/core/parsers/`, not invented: the `## Requirements`
section, `### Requirement:` inside it, any non-fenced `#### ` heading with
content, fences masked, a delta header in a main spec an error. A test writes
each of the 13 fixtures into a temporary project and runs
`openspec validate <spec> --type spec` against it, asserting the two tools agree
on whether the file is a spec. All 13 agree.

**The refusal became a report.** `enter` always descends. A file that does not
fit opens as a report naming every reason and the line it sits on, with `E` to
open the file in an editor and `esc` back to the markdown. Repair it and the next
scan turns the report into an outline without leaving the view.

Measured after the change, over every `openspec/specs/*/spec.md` under /home/pim:
605 specs, 469 structured, **0 blank scenarios**. Before it, 254 cards were blank.
The 136 that do not structure are files `openspec validate` also rejects; 102 of
them carry a delta header in a main spec.

Two things found while building it, neither in the plan:

- `document.path` styling ran after wrapping, so a backticked span straddling a
  wrap lost its colour and kept its marks on screen. Visible in the report,
  latent in the card since the previous change. Styled before wrapping now, with
  a test.
- The fence mask was blanking content as well as headings, which dropped a
  fenced code block out of the scenario that contained it. The mask now applies
  to heading detection only, which is what openspec does.

Two tasks were reworded rather than implemented as written, both noted in
`tasks.md`: 5.2 asked for `docActive()` to be inert in the report state, which is
wrong, because a report with several faults is longer than a 60x20 terminal and
has to scroll; and 1.4's rule was confirmed against openspec's own comment rather
than assumed.

The corpus test that missed this walked only this project's 24 specs, all one
shape by one author. `src/ui/testdata/specs/` now holds 13 specimens, each with a
comment naming the live project it came from.

Coverage: ui 88.5% (floor raised 87.6 to 88.3), total 90.0% (89.2 to 89.8). Gate
green. Both revert checks bit, with the build confirmed first.

Deferred to `specgetty-p44t`: the per-project problems report and the mechanical
repair, which would reach 130 of the 136 files that do not structure.
