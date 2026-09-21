## Context

See proposal.md for why. What matters here is that OpenSpec's parser is
available to read, so the grammar is not a guess.

The relevant files in openspec 1.10.0:

```
  dist/core/parsers/markdown-parser.js     sections by heading level, parseSpec
  dist/core/parsers/requirement-text.js    SCENARIO_HEADER, countScenarios
  dist/core/parsers/requirement-blocks.js  scenarioNameAt, requirement blocks
  dist/core/parsers/spec-structure.js      findMainSpecStructureIssues
  dist/core/parsers/code-fence.js          buildCodeFenceMask
  dist/core/validation/constants.js        the message wording
```

`spec-structure.js` is the closest thing to a specification of a main spec, and
its messages are written for a reader: "Main spec contains delta header ...
Delta headers are only valid inside openspec/changes/... and truncate the parsed
## Requirements section", and "Requirement header ... appears outside the main
## Requirements section. Main specs only parse requirements inside that section,
so this requirement is currently invisible to validate, list, and archive."

## Goals / Non-Goals

Goals: read what OpenSpec reads; never drop a line of a scenario; say why a file
cannot be structured, in terms that let the reader fix it.

Non-Goals: validating a spec in specgetty's own right, repairing one, reporting
across a project, or teaching the clause renderer new shapes.

## Decisions

### The grammar is transcribed, not invented

Every rule this change adds exists in openspec's source, and each is written in
the code as the regex it is. Transcribing them keeps two questions apart that
specgetty has been answering as one: whether a file is a spec, which OpenSpec
decides, and how to draw it, which specgetty decides.

The alternative was to shell out to `openspec show <spec> --json` and parse
nothing. It is authoritative and it was measured: the JSON carries `overview`,
`requirements[].text` and `scenarios[].rawText`, and it drops both the
`### Requirement:` name and the `#### Scenario:` name. I probed the output for
three names from `change-list-view` and none appear. An outline of unnamed nodes
is not an outline, so the file has to be read here either way. It also costs a
process per spec, and specgetty already knows what that feels like from
`openspec schema which`.

### A scenario is raw text, and clauses are a rendering of it

This is the shape of OpenSpec's own model and the fix for the reported bug at
once. A scenario node carries its content as an ordered list of parts, each
either a clause or a paragraph of prose:

```
  specPart{kind: partClause, keyword: "WHEN", text: "..."}
  specPart{kind: partProse,  text: "..."}
```

Ordered rather than "clauses plus a body", because the two interleave in real
files: this corpus has 68 horizontal rules and a number of `**Rationale**:`
blocks sitting inside scenario blocks. Clauses-then-prose would reorder a
scenario's content, and reordering a behaviour contract is worse than showing it
plainly.

A horizontal rule is dropped rather than kept as prose. It is a separator
between requirements in the file, and a lone `---` in a card is chrome with no
meaning.

### The clause test stays as it is

`- **WHEN**` and `- WHEN`, case-sensitive, as today. Of 5262 clause lines in the
specs that will still open, 5080 are in those two forms. The other 182, in 30
files, are `**When**` and `- **Given**`, and they will render as prose.

A better rule exists and is recorded here rather than taken: let the markup be
the signal instead of the case. A bold-wrapped keyword in any case is a clause,
because the author marked it up as a label; a bare keyword is a clause only in
ALL CAPS, because title case is a sentence. Every shape in the 544-spec survey
falls on one side of that with no ambiguous case. It is not taken now because
once nothing is dropped this is presentation only, and a presentation choice
made with the numbers in hand later is better than one bundled into a bug fix.

### The refusal descends instead of turning back

Today `enter` on a non-conforming spec leaves the cursor where it is and puts
one line on the nav bar. That line cannot hold what the reader now needs: three
reasons, the line each sits on, and what to do about them. So the spec level
gains a second state, and `enter` always descends:

```
  specs tab          markdown, readable, unchanged
       | enter
       v
  spec level --+-- it fits ------> outline | card
               |
               +-- it does not --> the reasons, and their lines
                                   E    open the file in your editor
                                   esc  back to the markdown
```

No fallback to the markdown inside the report: it is one `esc` away on the tab
above, where it already is. The alternative, refusing before descending, keeps
the model simpler and has nowhere to put the reasons.

### The reasons are written for someone who is about to fix the file

Each reason names the rule, the consequence, and the line. `E` shipped in the
previous change, so the report can hand the file straight to an editor. The
wording follows openspec's own messages, because the reader may well run
`openspec validate` next and two tools disagreeing about the same file is worse
than either being terse.

### The fixture corpus

The test that should have caught this walked this project's 24 specs, all
written in one shape by one author, which is why a parser whose whole job is
reading other people's files shipped unable to read them. `src/ui/testdata/specs/`
gets one file per shape found in the survey, each named for what it is:
a delta-headed main spec, a spec with no Purpose, bare uppercase clauses,
bold title-case clauses, a `#### ` header without the `Scenario:` prefix, a
`#### ` block with no content, a heading inside a fenced block, a scenario whose
content is prose. The corpus walk over the live specs stays, because it is the
thing that will notice the next shape.

## Risks / Trade-offs

- 136 of 544 live specs will refuse to open in the detail view, concentrated in
  whole projects (53 of 58 in one). The markdown on the tab above is untouched,
  so this costs structure, not access. It is also true, and no other tool in the
  chain can see those requirements either.
- A file that parses today and refuses after this change is a visible
  regression to anyone who has not read the report. The report is the mitigation
  and is why it is a panel rather than a line.
- Ordered parts make a scenario node more complex than clauses plus a body. The
  interleaving in the corpus is real, so the simpler model would be wrong
  rather than merely smaller.

## Migration Plan

None. Every file that opens today and still fits opens unchanged; the 254 blank
cards gain their content; the files that stop opening gain an explanation.
