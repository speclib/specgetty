---
# specgetty-vfnn
title: sometimes scenarios are wrongly shown empty
status: in-progress
type: bug
priority: high
created_at: 2026-09-21T20:19:39Z
updated_at: 2026-09-21T20:55:10Z
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
