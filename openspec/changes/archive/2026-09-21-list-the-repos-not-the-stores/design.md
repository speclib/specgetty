## Context

`follow-the-store` kept a property the walk already had: a repo with a
configuration and no content failed `isValidOpenSpecDir`, so it was never
listed. That looked like a free win. It was the wrong half of the pair.

```
  before this change                    after
  ─────────────────                     ─────
  listed:   the store                   listed:   every repo with a config
  skipped:  every repo reading it       skipped:  every registered store

  gh.nivis-project/nivis-tunnel   ✗      gh.nivis-project/nivis-tunnel   ✓
  gh.nivis-project/tf-provider..  ✗      gh.nivis-project/tf-provider..  ✓
  nivis-openspec-stores/nivis-..  ✓      nivis-openspec-stores/nivis-..  ✗
```

## Decisions

### The rule inverts rather than widening

Listing both was the other candidate, and it doubles the list without answering
which row to pick. A person opens `nivis-tunnel` to work on nivis-tunnel; the
store is where the file happens to sit. One row per workplace is the list they
want, and the store's name on that row is the thing they need to know about it.

The cost is real and accepted: a registered store no one points at is no longer
reachable from the picker. `spg` inside it and `--path` still open it, and
whether that deserves a row is a question with its own answer, not a side effect
of this one.

### A store is what the registry says it is

The identity file alone was never a sound test. `gh.nivis-project/ospecs` holds
`id: nivis` and is not the registered `nivis`:

```
  registry:  nivis -> nivis-openspec-stores/nivis
  on disk:   nivis-openspec-stores/nivis/.openspec-store/store.yaml   id: nivis   <- the store
             gh.nivis-project/ospecs/.openspec-store/store.yaml       id: nivis   <- a leftover
```

Both were shown, both called `nivis`, told apart only by a parenthesised parent.
The declaration path already demanded that the registry key and the metadata id
agree, and refused with `store_identity` when they did not. This applies the
same test to a directory met head on, so the two paths stop disagreeing about
what a store is.

### Discovery resolves, which makes the store's statistics the row's

A row for `nivis-tunnel` reporting zero specs would be worse than no row. So
discovery follows the declaration and reports the store's counts, and the file
listing the `:` search reads comes from the store too, or searching inside a
project would miss the specs the project shows.

That makes the same store parsed once per pointing repo, and two repos sharing a
store is exactly the case this change makes ordinary. The scan memoises by root,
so the nivis store's ninety archived changes are read once per scan rather than
once per row.

Git state stays out of this path. It was never affordable per discovered
project and still is not; it is read when a project is opened, as before.

## Risks

### The same specs on more than one row

`nivis-tunnel` and `terraform-provider-nivis-tunnel` both read store
`nivis-tunnel`, so two rows carry six specs and eight archived changes each.
That is not a defect, it is the truth about the tree, and it is why the store
column exists rather than being dropped once stores stopped having rows of their
own.

```
  project                          store           specs  archived
  nivis-tunnel                     nivis-tunnel    6      8
  terraform-provider-nivis-tunnel  nivis-tunnel    6      8
```

### A project opened by path may have no row

The picker's cursor keys on the open project's root. Opening a registered store
with `--path` gives a root that is deliberately not in the list, so the lookup
must miss quietly rather than selecting the wrong row or failing. That is a
scenario, not a comment.

### Registry reads during a walk

The exclusion asks a registry question of every candidate directory. The
registry is read once per scan and held for its duration; the per-candidate cost
is the stat for `.openspec-store/store.yaml` that already happens, and a map
lookup. A machine with no registry pays one failed open and excludes nothing.

### Rollback

Its own commit. Reverting restores store rows and drops the repo rows, which is
what ships today, and touches nothing outside discovery and the picker column.
