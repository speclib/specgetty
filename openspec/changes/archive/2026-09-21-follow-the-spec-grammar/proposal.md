## Why

Bean `specgetty-vfnn`. Scenarios are shown empty. In
`/home/pim/mipnix/openspec/specs/airplane-mode/spec.md` the detail view opens
and every scenario card is a bare title.

The cause is one line. A scenario is the only node kind with no fallback, so a
line the clause matcher does not recognise is dropped rather than kept:

```go
if cur.kind == nodeScenario {
    if kw, rest, ok := clauseOf(line); ok { ... }   // recognised
    if pending != nil { ... }                       // its continuation
    continue                                        // everything else: gone
}
cur.body += line + "\n"                             // what every other node does
```

Measured over the 544 live specs on this machine: 27 open with at least one
blank card, 254 cards in total, and **every one of the 254 is content the parser
dropped**. None is a scenario that is genuinely empty in the file. Ten of those
27 specs pass `openspec validate --strict`.

That last number is the finding. It means the shapes being dropped are not
malformed, so widening the clause matcher until the reported file works would
only move the boundary to the next author's habit.

### There is an official grammar, and it says something different

Transcribed from openspec 1.10.0, `dist/core/parsers/`:

| rule                   | pattern                                                |
| ---------------------- | ------------------------------------------------------ |
| the requirements home  | `/^##\s+Requirements\s*$/i`                            |
| purpose                | a section titled `Purpose`, case-insensitive           |
| a requirement          | `/^###\s+Requirement:\s*(.+)\s*$/i`                    |
| a scenario             | `/^####\s+/`, any text, with content under it          |
| a delta header         | an ERROR in a main spec: it truncates `## Requirements` |
| code fences            | masked everywhere                                      |

And the whole of how a scenario's content is modelled:

```js
if (scenarioSection.content.trim()) {
    scenarios.push({ rawText: scenarioSection.content });
}
```

A scenario is `{ rawText }`. There is no GIVEN/WHEN/THEN in OpenSpec's data
model at all. `openspec show <spec> --json` confirms it: the scenario objects
carry `rawText` and nothing else. The `- **WHEN**` form appears in the template,
in the schema instruction and in a validator error message, which is three
pieces of guidance and no grammar. Nothing will ever reject a spec for writing
`GIVEN` bare, so specgetty cannot rely on the convention.

So specgetty is strict where OpenSpec is silent, and lax where OpenSpec is
strict. Against the official grammar, 136 of the 544 live specs are invalid,
102 of them because a main spec carries a delta header. Those files are
invisible to `openspec list`, `archive` and `validate` too, and specgetty opens
them today and shows half of what they contain.

## What Changes

- The structure specgetty reads becomes the structure OpenSpec defines: the
  `## Requirements` section, `### Requirement:` inside it, any non-fenced
  `#### ` header with content under it, code fences masked, and a delta header
  in a main spec treated as the error OpenSpec calls it.

- A scenario's content becomes raw text that is never dropped. Laying a clause
  out under its keyword becomes presentation applied when the text is
  recognisable, with the prose shown as prose when it is not. This is what fixes
  all 254 blank cards, and it fixes them for shapes nobody has written yet.

- A file that does not fit the grammar no longer refuses at the door. `enter`
  descends into a report naming each reason and the line it is on, because a
  one-line nav bar message cannot hold three reasons and the fix for them. The
  whole file stays readable as markdown one level up, which is where it already
  is, so nothing becomes unreachable.

- Adopted with the grammar, each worth its own scenario: a `#### ` header need
  not read `Scenario:` (24 such headers in this corpus are invisible today), a
  `#### ` block with no content under it is not a scenario, and a heading inside
  a fenced code block is not a heading (2 in this corpus, and specgetty treats
  one as a node boundary today).

- The corpus test gains a fixtures directory holding one specimen per shape
  found in the wild. The test that should have caught this walked only this
  project's own 24 specs, every one written in one shape by one author.

Not in scope, and the reason for each:

- **Widening which clause shapes are laid out.** Once nothing is dropped this is
  presentation only. Of 5262 clause lines in the specs that will still open,
  5080 are in the two forms already handled; the remaining 182 will render as
  prose, correctly, in 30 files. Revisit it with the numbers rather than now.
- **Repairing a broken spec file, and reporting problems per project.** A
  mechanical fix would repair 130 of the 136 invalid files, and whole projects
  are affected rather than scattered files. That deserves its own proposal: it
  would be the first time specgetty writes anything to a spec file.

## Capabilities

### Modified Capabilities
- `spec-detail-view`: the grammar it reads, a scenario as raw text, and the
  refusal becoming a report inside the view
- `specs-tab`: `enter` descends whether or not the file fits, so the reason no
  longer goes on the nav bar

## Impact

- `src/ui/specparse.go`: rewritten against the grammar; returns either a tree or
  the reasons it is not one, each with a line number
- `src/ui/specdetail.go`: the report, and a card that renders raw text
- `src/ui/ui.go`: the spec level gains a second state
- `src/ui/testdata/specs/`: new, one fixture per shape in the wild
- `README.md`: what opens, what does not, and why

## Rollback

Its own commit. Reverting returns the parser to dropping unrecognised lines,
which is the bug this change exists to remove, and returns the refusal to the
nav bar.
