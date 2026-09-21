---
# specgetty-cn8h
title: picker lists the repos you work in, not the stores behind them
status: in-progress
type: feature
priority: high
created_at: 2026-09-21T14:18:49Z
updated_at: 2026-09-21T14:21:56Z
---

`follow-the-store` (specgetty-8quc) chose store rows only: a repo that declares a `store:` and keeps no content of its own is not listed, and the store it points at is listed instead.

Against a real tree that is the wrong way round. Five working repos under gh.nivis-project (nivis, nivis-demos, nivis-tunnel, registry, terraform-provider-nivis-tunnel) are absent from the picker, and the four store directories stand in for them.

It also created a functional gap: the repo's own `openspec/config.yaml`, which carries that repo's context and rules, is only reachable when the project is opened from the repo. Opening the store row has no origin, so that sub-tab does not exist and those lines cannot be reached from the picker at all.

Second, smaller fault: a directory carrying `.openspec-store/store.yaml` is named by the id in that file whether or not the registry resolves that id to it. gh.nivis-project/ospecs claims id `nivis` while the registered `nivis` is elsewhere, so an unregistered leftover wears the real store's name.

Decisions taken: the picker lists repos and plain projects and drops store rows; a directory is only treated as a store when the registry actually resolves to it.

## OpenSpec change

`list-the-repos-not-the-stores` (`openspec/changes/list-the-repos-not-the-stores/`),
validated strict. Not yet implemented.
