---
# specgetty-cn8h
title: picker lists the repos you work in, not the stores behind them
status: completed
type: feature
priority: high
created_at: 2026-09-21T14:18:49Z
updated_at: 2026-09-21T14:44:13Z
---

`follow-the-store` (specgetty-8quc) chose store rows only: a repo that declares a `store:` and keeps no content of its own is not listed, and the store it points at is listed instead.

Against a real tree that is the wrong way round. Five working repos under gh.nivis-project (nivis, nivis-demos, nivis-tunnel, registry, terraform-provider-nivis-tunnel) are absent from the picker, and the four store directories stand in for them.

It also created a functional gap: the repo's own `openspec/config.yaml`, which carries that repo's context and rules, is only reachable when the project is opened from the repo. Opening the store row has no origin, so that sub-tab does not exist and those lines cannot be reached from the picker at all.

Second, smaller fault: a directory carrying `.openspec-store/store.yaml` is named by the id in that file whether or not the registry resolves that id to it. gh.nivis-project/ospecs claims id `nivis` while the registered `nivis` is elsewhere, so an unregistered leftover wears the real store's name.

Decisions taken: the picker lists repos and plain projects and drops store rows; a directory is only treated as a store when the registry actually resolves to it.

## OpenSpec change

`list-the-repos-not-the-stores` (`openspec/changes/list-the-repos-not-the-stores/`),
validated strict. Not yet implemented.

## Summary of Changes

Shipped as `list-the-repos-not-the-stores` (commit 6cb3984), archived to
`openspec/changes/archive/2026-09-21-list-the-repos-not-the-stores/`.

The rule inverted. A repo with an OpenSpec configuration is listed whether or
not it keeps content of its own; a registered store is not listed, because it is
where content lives rather than where work happens.

- `find.go`: a configuration is enough, content is not required, and a
  registered store is excluded. The registry is read once per walk and handed
  down rather than read per candidate.
- `store.go`: `StoreInfoWith` requires the registry to resolve the declared id
  to that same directory, so `ospecs` stopped wearing the `nivis` name. Paths
  are compared through symlinks, since the registry stores a canonical path.
- `scan.go`: discovery resolves, so a store-backed repo carries the store's
  specs, changes, tasks and file listing instead of the nothing it holds. A
  shared store is parsed once per scan. Projects are keyed by where the reading
  STARTED, which is what makes two repos on one store two rows; the root
  travels inside the info as `Root`.
- `ui`: `currentRoot()` reads `Info.Root` rather than the map key, so every
  action still targets the store. The picker's `kind` column became a `store`
  column naming the store a row reads from.

Verified against the real tree: all five repos listed, each naming its store and
carrying its counts; `nivis-tunnel` and `terraform-provider-nivis-tunnel` both
show 6 specs and 8 archived, visibly sharing; `ospecs` under its own name; the
four store directories gone; no row called `nivis` twice. A store opened by path
still works.

Coverage floors raised: scanner 93.5 to 94.5, ui 81.7 to 82.5, total 84.0 to
85.0. README gained a section on stores and what the picker lists.
