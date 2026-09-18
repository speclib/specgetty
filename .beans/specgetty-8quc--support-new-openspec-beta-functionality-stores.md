---
# specgetty-8quc
title: 'support new openspec beta functionality: stores'
status: in-progress
type: epic
priority: normal
created_at: 2026-09-15T15:27:04Z
updated_at: 2026-09-18T16:14:38Z
---

openspec stores support should add the following functionality:

- when in a project dir automatically detect this project is using a store
- show autodetected store contents
- give some kind of visual sign this is project is using a store

## OpenSpec change

`follow-the-store` (`openspec/changes/follow-the-store/`) covers all three
bullets above: resolution through the registry, the store's content shown in
place, and a mark in the project header explained on the config tab.

Deferred to follow-up beans, recorded in that change's tasks.md:

- the reverse index: which repos point at a given store
- `--store <id>` and the machine-wide `defaultStore` resolution sources
