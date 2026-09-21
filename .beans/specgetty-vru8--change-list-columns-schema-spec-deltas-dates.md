---
# specgetty-vru8
title: 'change list columns: schema, spec deltas, dates'
status: draft
type: feature
priority: normal
created_at: 2026-09-15T13:45:03Z
updated_at: 2026-09-21T16:45:07Z
blocked_by:
    - specgetty-g4pa
---

Split out of specgetty-g4pa (change list full view).

The column data that needs new scanner capability, beyond what a scan already computes.

## Fields still to source

| Column              | Source                                     |
|---------------------|--------------------------------------------|
| schema              | <change>/.openspec.yaml -> schema:         |
| complete per schema | schema artifacts[].generates vs disk       |
| new specs           | '## ADDED Requirements' in delta spec      |
| spec modifications  | '## MODIFIED' / '## REMOVED'               |
| created             | .openspec.yaml -> created:                 |
| updated             | max mtime in change dir                    |

## Open question: how to resolve the schema

'complete according to schema' needs the schema definition, which lives outside
the project. 'openspec schema which spec-driven' resolves into the openspec
package directory (a nix store path here). specgetty's scan path touches nothing
but the filesystem today; only the archive action shells out to the CLI.

| Approach                                    | Cost                          | Fragility        |
|---------------------------------------------|-------------------------------|------------------|
| shell out to openspec per project           | forks a CLI per project scan  | version drift    |
| hardcode spec-driven's artifact set         | trivial                       | custom schemas   |
| resolve schema.yaml ourselves               | medium                        | duplicates rules |
| drop the column, show artifact count        | free                          | loses the signal |

## Known gap

Archived changes have no .openspec.yaml (verified on
openspec/changes/archive/2026-04-20-discard-change-action). Schema and created
date are unavailable for everything archived, whichever approach is picked. The
archive dir name prefix gives the archived date, not the created date.

## Depends on

The --change-fields mechanism and the column renderer land in the first change.
This bean adds fields behind that same mechanism.

## Update (2026-09-21)

The open question above is answered by `properties-tab`: shell out to
`openspec schema which <name> --json` for the path, then parse
`<path>/schema.yaml` locally. Not per project scan, which was the cost that
made it look unaffordable: once per schema per project, and only once the
properties tab is opened.

`scanner.ResolveSchema` is the reader, and `ChangeInfo.Schema` already
carries each change's recorded schema, so the `schema` column here is a
display field over data that now exists.

A sortable column would additionally need a sort key per field: the rendered
value sorts wrong, because `tasks` produces `9/12` and `10/12` and string
order puts ten before nine. Recorded here because `group-the-change-list`
considered sorting and deliberately left it to this bean.

## Note from group-the-change-list (2026-09-21)

Sorting was considered there and deliberately left here. The change list gained
grouping and two fixed orders (active by name, archived newest first) instead,
because what sorting was wanted for was seeing recent work first.

A sortable column needs a sort key per field, not a comparator on the rendered
value: `tasks` produces `9/12` and `10/12`, and string order puts ten before
nine. `fieldDef` carries only `value func(T) string` today, and the project
picker shares that type.
