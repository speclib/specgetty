---
# specgetty-m4vt
title: reorganize config tab
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T14:36:40Z
updated_at: 2026-09-21T16:22:39Z
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
