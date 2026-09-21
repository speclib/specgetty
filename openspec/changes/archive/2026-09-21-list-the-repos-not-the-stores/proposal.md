## Why

Bean `specgetty-cn8h`. `follow-the-store` chose store rows only, and against a
real tree it is the wrong way round:

```
  what the picker lists today          what the user works in
  ───────────────────────────          ──────────────────────
  nivis (nivis-openspec-stores)        gh.nivis-project/nivis
  nivis-demos                          gh.nivis-project/nivis-demos
  nivis-tunnel                         gh.nivis-project/nivis-tunnel
  registry                             gh.nivis-project/registry
  nivis (gh.nivis-project)   <- ospecs gh.nivis-project/terraform-provider-nivis-tunnel
```

Not one of the five repos is listed. Four directories nobody edits directly
stand in for them, and the names collide besides.

It also left a gap that was not visible when the rule was chosen. A repo's own
`openspec/config.yaml` carries that repo's `context:`, `rules:` and
`operations:`, which the store deliberately does not hold. That configuration is
reachable only when the project is opened from the repo, because only then is
there an origin. Opening the store row has none, so the repo sub-tab does not
exist and those lines cannot be reached from the picker at all.

Second fault, smaller and separate. A directory is treated as a store on the
strength of its `.openspec-store/store.yaml` alone. `gh.nivis-project/ospecs`
carries one claiming `id: nivis` while the registered `nivis` is
`nivis-openspec-stores/nivis`, so an unregistered leftover wears the registered
store's name and the two rows differ only by a parenthesised parent.

## What Changes

- The picker lists the directories a person works in: every repo with an
  OpenSpec configuration, whether it keeps content of its own or declares a
  store, plus every plain project
- A registered store is not listed. It is where content lives, not where work
  happens, and every repo that reads from one already stands for it
- A row that reads from a store says which one, so two repos showing the same
  specs explain themselves rather than looking like duplicates
- A row for a store-backed repo carries the store's statistics. A row reporting
  zero specs and zero changes for a repo whose specs are one directory away
  would be worse than no row
- A directory counts as a store only when the registry resolves its declared id
  to that directory. An unregistered leftover is named by its folder and carries
  no store marking

Not in scope: reaching a registered store from the picker. `spg` inside the
store's directory and `--path` both still open it, and this change accepts that
the picker is a list of workplaces rather than an index of everything OpenSpec
can see. If a store nobody points at turns out to need a row, that is a separate
decision with its own rule.

Not in scope: the reverse index. Naming a store's readers is `specgetty-a2zv`
and is unaffected either way.

## Capabilities

### Modified Capabilities
- `openspec-scanning`: the rule inverts. A repo declaring a store is listed; a
  registered store is not
- `project-picker`: rows are named by their directory, and a row reading from a
  store names it and carries its statistics
- `store-resolution`: a store is a directory the registry resolves to, not any
  directory carrying an identity file

## Impact

- `src/scanner/find.go`: the validity rule, and the registry consulted once per
  walk rather than once per candidate
- `src/scanner/store.go`: `storeInfoAt` gains the registry check
- `src/scanner/scan.go`: discovery resolves each path, so a pointer-only repo
  arrives with the store's statistics. The git state stays out of this path, as
  it was never affordable per discovered project, and is still read only when a
  project is opened
- The same store is parsed once per scan rather than once per repo pointing at
  it, since two repos sharing a store is the case this change makes common
- `src/ui/picker.go`: the column that said `store` now names the store a row
  reads from

## Rollback

Its own commit. Reverting restores store rows and drops the repo rows, which is
what ships today.
