---
# specgetty-vru8
title: 'change list columns: schema, spec deltas, dates'
status: draft
type: feature
priority: normal
created_at: 2026-09-15T13:45:03Z
updated_at: 2026-09-15T13:45:21Z
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
