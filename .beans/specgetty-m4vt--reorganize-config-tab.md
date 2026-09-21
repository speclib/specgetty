---
# specgetty-m4vt
title: reorganize config tab
status: completed
type: task
priority: normal
created_at: 2026-09-21T14:36:40Z
updated_at: 2026-09-21T16:38:09Z
---

## currently: 
- HORIZONTAL: changes|specs|config
  - HORIZONTAL: repo|store|store details

## wanted: 
- HORIZONTAL: changes|specs|settings
  - VERTICAL:
    - project: syntax highlighted contents of config.yml
    - schema: single view with [which schema is set] and [details of schema configuration]
    - store: local or contents of what now is shown in store details

## OpenSpec change

`properties-tab` (`openspec/changes/properties-tab/`), validated strict. Not
yet implemented.

Decisions taken during exploration:

- tab renamed `properties`; vertical list beside the content, reusing the specs
  tab's split
- sections: `project`, one row per schema in use, `store`
- `project` shows the resolved root's configuration, which for a store-backed
  repo is the store's. The pointing repo's `context`, `rules` and
  `operations` are inert, proven in a sandbox; openspec calls them
  `inertPointerDeclarations` and warns about none of them, so `store` names
  them instead
- schema definitions located by `openspec schema which <name> --json`, then
  parsed locally; one call per used schema, concurrent, on tab open, kept for
  the session, dropped on project open only
- an open change names its schema; the change-list column stays in
  [[specgetty-vru8]]

## Summary of Changes

Shipped as `properties-tab` (commit 8f7aab5), archived to
`openspec/changes/archive/2026-09-21-properties-tab/`.

- the config tab is `properties`, drawn as a row list beside its content using
  the specs tab's split; `focusSpecsList`/`focusSpecsContent` became a
  generic pair now that two tabs are splits
- rows: `project` (the one configuration that applies, headed by the file it
  came from), one per workflow schema in use, and `store`
- `scanner.ResolveSchema` locates a definition with
  `openspec schema which <name> --json` and parses `schema.yaml` locally.
  One call per used schema, concurrent, on tab open, kept for the session,
  dropped on project open
- `ChangeInfo.Schema` reads each change's `.openspec.yaml`; an open change
  names its schema, and changes recording none are counted apart
- the store row names the keys a declaring repo carries in vain

## Corrected along the way

`follow-the-store` believed a declaring repo's own `context`, `rules` and
`operations` still applied. They do not: OpenSpec reads that file for
`store:` alone. Proven in a sandbox and against the source, where the term is
`inertPointerDeclarations`. The `openspec-scanning` requirement that endorsed
it is corrected, and the test built on the old premise was replaced with one
asserting the opposite.

## Notes

- `config-tab-display` keeps its directory name: OpenSpec's RENAMED works on
  requirements and has no capability-level form. Recorded in its Purpose.
- a real bug the timeout test caught: killing a process on a deadline does not
  close an output pipe a child inherited, so `Output()` waited for the child.
  `cmd.WaitDelay` is what makes the bound real.
- coverage floors raised: scanner 94.5 to 94.8, ui 82.5 to 83.3, total 85.0 to
  86.0
