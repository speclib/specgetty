---
# specgetty-8quc
title: 'support new openspec beta functionality: stores'
status: completed
type: epic
priority: normal
created_at: 2026-09-15T15:27:04Z
updated_at: 2026-09-18T16:37:56Z
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

## Summary of Changes

Shipped as `follow-the-store` (commit 1ff4785), archived to
`openspec/changes/archive/2026-09-18-follow-the-store/`.

All three bullets are covered:

- **Autodetect.** A new resolver in `src/scanner/store.go` reimplements
  OpenSpec's own root rules: walk up to the nearest qualifying `openspec/`,
  prefer local content over a pointer, otherwise follow `store: <id>` through
  `$XDG_DATA_HOME/openspec/stores/registry.yaml`. No network, and no call out
  to the `openspec` binary.
- **Show the contents.** `ProjectMap` is keyed by the resolved root, so reads
  and the archive, discard and export actions all land in the store. Discard
  used to create an unused `changes/discarded/` in the pointing repo.
- **Visual sign.** A `store` chip in the project header, and nothing more.
  The config tab explains it: sub-tabs for the repo's own config, the store's
  shared config, and a store report (id, root, origin, registered remote and
  branch, plus local git state read without fetching).

Also: the picker names a store by its declared id and marks the row, repos that
only point at a store stay out of the list, the watcher covers both trees,
`config.yml` is read alongside `config.yaml`, and an unfollowable declaration
is reported where the project would otherwise just look empty.

A parity test puts the resolver and `openspec context --json` to the same
fixtures, so a drift in the beta fails a test rather than showing wrong content.

Coverage floors raised: src 42.3 to 85.9, scanner 91.9 to 93.5, ui 79.7 to 81.7,
total 79.7 to 84.0. `git` was added to the nix check inputs, since the store
git-state paths are only exercised against a real repository.

Follow-ups: [[specgetty-a2zv]] (which repos use a store) and [[specgetty-lexl]]
(`--store` and `defaultStore`).
