## Context

OpenSpec 1.10 resolves a root before it does anything else. specgetty has never
had that step, because until stores existed the answer was always "the directory
holding `openspec/`". The rules below are read from openspec 1.10.0's
`core/root-selection.js`, `core/project-config.js` and `core/store/foundation.js`,
not inferred from behaviour.

```
  resolve(startDir):
    walk up for the nearest QUALIFYING openspec/ dir
      qualifying = holds specs/ or changes/   OR   holds a config file
      |
      +-- holds specs/ or changes/   -> THIS is the root          source: nearest
      |     a store: key here is IGNORED (openspec warns on stderr)
      +-- no content, store: <id>    -> registry lookup           source: declared
      +-- no content, no store:      -> THIS is the root          source: nearest
    nothing found -> global defaultStore in ~/.config/openspec/config.json
                                                                  source: global_default
```

Three facts settle most of the design. A store is always a local directory, so
nothing here touches the network. Several stores can live in one git working
copy, as the four under `nivis-openspec-stores/` do, so a store is not a repo.
And the mapping is many to one: `nivis-tunnel` and
`terraform-provider-nivis-tunnel` both declare `store: nivis-tunnel`.

## Decisions

### Identity and content become two fields, not two maps

`ProjectMap` stays keyed by the root, because with stores as picker rows the
keys are already roots. Two annotations carry the rest:

```go
StoreID string   // non-empty when this root is a store
Origin  string   // the repo resolution started from; empty when it is the root
```

Every broken call site then has one unambiguous answer:

```
  root                              origin
  ────                              ──────
  ParseProjectInfo                  the header's path, when set
  ListOpenSpecContents              the repo config sub-tab
  doArchiveChange (cmd.Dir)         the second watched tree
  doDiscardChange, doExportChange
  the first watched tree
```

An alternative was a second map keyed by origin. It buys nothing: origin is
never a key, only a label and a watch target.

### The scanner is left alone

`isValidOpenSpecDir` requires a config file AND (`specs/` OR `changes/`). A
store-backed repo has the first and neither of the second, so it is already
excluded. A store created by `openspec store setup` gets `openspec/config.yaml`,
`openspec/specs/` and `openspec/changes/archive/` from birth, verified against
1.10.0 in a sandbox, so every store is already found.

```
  walk finds                                   walk skips
  ──────────                                   ──────────
  nivis-openspec-stores/nivis-tunnel     ⬡     gh.nivis-project/nivis-tunnel
  nivis-openspec-stores/nivis            ⬡     gh.nivis-project/terraform-provider-...
  nivis-openspec-stores/nivis-demos      ⬡
  nivis-openspec-stores/registry         ⬡
  gh.speclib/specgetty
```

That is the wanted row set exactly. What the scanner gains is a label, read from
`.openspec-store/store.yaml`, and the requirement that the current exclusion is
deliberate rather than incidental. The project cache stores paths, and those
paths are roots before and after, so its format does not move.

### The header gets a mark and nothing else

```
 /home/pim/gh.speclib/specgetty                        plain project
 Specs: 22  Changes: 0 active  Archived: 33

 /home/pim/gh.nivis-project/nivis-tunnel  ⬡            resolved through a store
 Specs: 6  Changes: 0 active  Archived: 12

 nivis-tunnel  ⬡                                       the store, from the picker
 Specs: 6  Changes: 0 active  Archived: 12
```

The header's job is to say where you are, and where you are is the repo you ran
`spg` in. What the mark adds is that the numbers beneath it were not read from
there. Putting the store's path on the second line was the obvious alternative
and was rejected: it costs a row on a header that is already four, and it answers
a question asked once per project rather than once per glance.

### The config tab splits, because a store-backed project has two of them

This is the part that would have been a silent regression. `nivis-tunnel`'s own
`openspec/config.yaml` carries thirty lines of repo-specific `context:`, `rules:`
and `operations:`, and the store's config says in as many words that the split is
intentional: "Each repo keeps its own openspec/config.yaml context, rules and
operations. What is here is what both sides share." Reading only the root's
config would drop the repo's half off the screen.

```
| [changes] [specs] [config]                            |
| repo · store · store details                          |
| ╭───────────────────────────────────────────────────╮ |
| │ schema: spec-driven                               │ |
| │ store: nivis-tunnel                               │ |
| │                                                   │ |
| │ context: |                                        │ |
| ╰───────────────────────────────────────────────────╯ |
```

The sub-tab row sits where the filename sits today, above the border. That
follows the rule `box-the-tab-content` settled: a line that names the content is
chrome, a line that reports on the content is content. Sub-tabs name; they stay
outside. A plain project has one configuration and keeps the single dimmed
filename, so nothing changes for it.

The third sub-tab is where the header's mark is explained, which is what makes
the mark affordable. Everything it reports is a local read:

```
  store id        .openspec-store/store.yaml
  root path       registry backend.local_path
  origin repo     resolution
  remote, branch  registry entry, and store.yaml's canonical remote
  git state       origin url, uncommitted changes, ahead/behind upstream
```

openspec's own `core/store/git.js` documents every one of those as local-only,
and is explicit that ahead/behind compares against the current local upstream
ref rather than the live remote. The details sub-tab says so, because a number
that looks like it was fetched and was not is worse than no number.

### Both watched trees, and why the second one is not paranoia

The store's tree is where changes and specs move. The repo's `openspec/` holds
one file, but that file is the pointer. Editing `store:` from `nivis-tunnel` to
`nivis` changes everything on screen without touching the store at all, and a
watcher on the store alone would show stale content with no sign anything had
happened. The second tree is one directory holding one file, so it costs one
inotify watch.

## Risks

### Reimplementing resolution means owning the edge cases

Parsing directly was chosen over shelling out to `openspec context --json`,
which would be authoritative but would cost a process per project in the picker.
The price is that these are now specgetty's rules to keep:

```
  planning shape beats a pointer       a repo with changes/ AND store: is its own root
  store: must be a string              a list or a map is malformed, not absent
  config.yml counts                    openspec probes .yaml then .yml
  id must match                        registry key vs .openspec-store/store.yaml id
  XDG_DATA_HOME wins when set          on every platform, not only Linux
```

Each is a scenario in `store-resolution`. The mitigation that matters is the
last line of the tasks: the resolver's answer is compared against
`openspec context --json` for a set of fixtures, so a drift in the beta shows up
as a failing test rather than as wrong content on screen. That comparison is a
test, never a runtime dependency.

### A pointer that does not resolve must not look like an empty project

Today's failure is `Specs: 0  Changes: 0 active`, which reads as a project with
nothing in it. Every unresolvable case has to be louder than that: unregistered
id, no registry at all, malformed `store` value, registered path gone, id
mismatch. This is the one behaviour where doing nothing is indistinguishable
from the bug being fixed, so it gets its own scenarios and its own tasks.

### Two ways into the same content

`spg` in `nivis-tunnel` and picking `nivis-tunnel` from the picker show the same
specs under different headers, one with an origin and one without. That is
truthful rather than confusing, but it does mean the header has two store states
to test, not one, and that the picker cursor has to find the store's row when the
open project has no row of its own.

### Rollback

Its own commit. A revert restores the empty-project behaviour for a store-backed
repo and touches nothing else, because the scanner, the cache format and the
plain-project paths are all unchanged by this work.
