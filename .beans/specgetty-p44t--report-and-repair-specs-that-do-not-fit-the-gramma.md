---
# specgetty-p44t
title: report and repair specs that do not fit the grammar
status: todo
type: task
priority: normal
created_at: 2026-09-21T20:54:52Z
updated_at: 2026-09-21T20:54:52Z
---

Split out of `specgetty-vfnn` while exploring it. That bean's change,
`follow-the-spec-grammar`, makes specgetty read a spec by OpenSpec's own grammar
and report why a file cannot be structured. This bean is the half that acts on
what the report found.

## What the survey found

Over the 544 live specs on this machine, 136 (25%) are invalid by OpenSpec's own
grammar, and they are concentrated in whole projects rather than scattered:

| project           | invalid |
| ----------------- | ------- |
| voorzetramenshop  | 53/58   |
| mipnix            | 34/78   |
| mip.rs            | 25/30   |
| awsaccounts       |  9/12   |
| bmc               |  7/13   |
| openspec.nvim     |  6/6    |
| quiqr-desktop     |  1/34   |

The causes are one mistake cascading. 102 main specs carry a delta header
(`## ADDED Requirements`) where a main spec needs `## Requirements`, which also
means they have no requirements section and their requirements sit outside one.
133 have no `## Purpose`.

These files are invisible to `openspec list`, `archive` and `validate` as well,
not only to specgetty.

## Two halves

**Report per project.** A `problems` row on the properties tab, beside the
project, schema and store rows, counting the spec files that do not fit and why.
Free to compute: specgetty already reads every spec's content at scan time.
Following the properties tab's own rule that the store row is always present
"because a row that says so teaches the concept where an absent one teaches
nothing", the row is always there and reads none for a healthy project.

```
 +---------------+ +----------------------------------------------+
 | project       | |  specs                58                     |
 | spec-driven   | |  valid                 5                     |
 | problems  53  | |  invalid              53                     |
 | store         | |                                              |
 +---------------+ |  delta header in a main spec            53   |
                   |  no ## Purpose section                  52   |
                   |  requirements outside ## Requirements   51   |
                   |  a requirement with no scenario          2   |
                   +----------------------------------------------+
```

The picker already shows spec, change and task counts per project. A problem
count there would say which of twenty projects are healthy before any of them is
opened.

**Repair.** A mechanical fix reaches 130 of the 136:

| fix                                   | files |
| ------------------------------------- | ----- |
| header rename alone                   |     3 |
| header rename + a Purpose placeholder |   127 |
| needs real authoring                  |     6 |

There is precedent for the placeholder: openspec's own `archive` writes
`TBD - created by archiving change X. Update Purpose after archive.` when a
delta carries no Purpose.

The rename has a real boundary, discoverable from the file rather than guessed:

| delta headers present | files | mechanical                                    |
| --------------------- | ----- | --------------------------------------------- |
| ADDED only            |    89 | yes: rename to `## Requirements`              |
| ADDED + MODIFIED      |    12 | rename the first, delete the rest, or two     |
|                       |       | `## Requirements` sections result and         |
|                       |       | `findSection` reads only the first            |
| includes REMOVED      |     1 | no: merging a removed requirement back        |
|                       |       | resurrects behaviour someone deleted          |

## Open questions

- **Scope.** Per-spec from the refusal report, or project-wide from the problems
  row? Per-spec is confirmed, scoped to a file on screen, and a natural extension
  of a report that already says what is wrong. Project-wide is worth far more
  (53 files in one keystroke) and is far more dangerous. Recommendation: per-spec
  first, project-wide only once it has proven it never surprises.
- **Whether specgetty should write at all.** `E` now opens the file in the user's
  editor, so the report could simply say why and leave it there. The case against
  is 127 files needing the identical two-line edit, which is work a tool should do
  and a person should not.
- Confirmation, a report of exactly what changed, and git as the undo. specgetty
  writes one thing today, a task checkbox, through an atomic write that preserves
  mode; the same helper applies.
- Does the problems row duplicate `openspec validate`? Partly. The argument for
  it is that specgetty reports status across many projects and this is a status.

Explored 2026-09-21. Numbers from a survey of every `openspec/specs/*/spec.md`
under /home/pim, excluding archives.
