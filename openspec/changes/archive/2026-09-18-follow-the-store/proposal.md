## Why

Bean `specgetty-8quc`. OpenSpec 1.10 lets a repo keep no specs and no changes of
its own, and point at a store instead:

```
  gh.nivis-project/nivis-tunnel/          nivis-openspec-stores/nivis-tunnel/
    openspec/                               .openspec-store/store.yaml   id: nivis-tunnel
      config.yaml   <- the only file        openspec/config.yaml
        store: nivis-tunnel  ────┐          openspec/specs/     6 capabilities
        context, rules           │          openspec/changes/archive/
                                 │                ^
                                 └────────────────┘
                   via $XDG_DATA_HOME/openspec/stores/registry.yaml
```

`spg` in such a repo opens an empty project: `Specs: 0  Changes: 0 active`. The
content it should be showing sits two directories away and it never looks. The
picker has the mirror-image fault, listing the four stores under
`nivis-openspec-stores/` as if each were a project named after its folder, while
the repos people actually work in are absent.

Underneath both is one conflation. `projectPath` carries two jobs that a store
pulls apart:

```
        today                               with a store
  projectPath ─┬─ identity           identity -> the repo you stand in
               └─ content            content  -> the root the specs live in
```

## What Changes

- The root is resolved by openspec's rules, parsed directly, with no call out to
  the `openspec` binary. A repo whose `openspec/` has no `specs/` and no
  `changes/` but does carry `store: <id>` resolves through the registry to the
  store's root, and every read comes from there
- A store is labelled as a store. The picker names it by the `id` in its
  `.openspec-store/store.yaml` rather than by its folder, and marks it
- A repo that only points at a store stays out of the picker. It is an entry
  point, not a row. The scanner already excludes it, because it has neither
  `specs/` nor `changes/`, so this is a property to keep rather than build
- The project header gains one mark saying this project uses a store. No id, no
  path: those are details, and they go where details go
- The config tab grows sub-tabs. A store-backed project has two configs, the
  repo's and the store's, and a third sub-tab carries the store details
- Archive, discard and export target the resolved root. Today they would write
  into the pointing repo, which for discard means silently creating
  `openspec/changes/discarded/` in a repo that has no changes
- The watcher watches both trees: the store, where the content moves, and the
  repo's `openspec/`, where the pointer itself can move
- `openspec/config.yml` is read where `openspec/config.yaml` is read. openspec
  accepts both and specgetty knows only the first; the pointer lives in that
  file, so the gap stops being cosmetic here

Not in scope: the reverse index, which store a given repo points at as seen from
the store side. Naming the repos on a store's detail sub-tab is the obvious
follow-up and is nearly free, because the walk already passes every pointing
repo's `openspec/` directory. It is left out so this change is about following
one pointer, not about building an index. Deferred to its own bean.

Not in scope: `--store <id>` as a command-line flag, and the machine-wide
`defaultStore` from `~/.config/openspec/config.json`. Both are fallbacks that
nothing on this machine uses yet.

Not in scope: fetching, cloning or any network. A store is always a local
directory; `remote` in the registry records where it came from and is never
followed.

## Capabilities

### Added Capabilities
- `store-resolution`: how the root that holds specs and changes is resolved from
  a starting directory, what a store is, and the rule that every filesystem
  operation targets the resolved root rather than the directory the user stands
  in

### Modified Capabilities
- `openspec-scanning`: a discovered root that carries store metadata is reported
  as a store, with its declared id
- `project-picker`: a store row is named by its store id and marked as a store
- `project-view`: the header marks a project that reads from a store, and
  startup resolution follows a pointer
- `config-tab-display`: the tab becomes a set of sub-tabs, because a store-backed
  project has more than one configuration to show
- `fs-watch`: two trees are watched when they differ
- `archive-from-ui`: the command runs in the resolved root

## Impact

- `src/scanner/`: a new resolver reading `openspec/config.yaml`, the registry and
  `.openspec-store/store.yaml`. `ProjectInfo` gains the store id and the origin
  repo. `find.go` is untouched: its row set is already the wanted one
- `src/ui/ui.go`: the header mark, the config sub-tabs, and the three action
  commands changing which path they are given
- `src/watcher/`: a second tree
- The project cache stores paths, and those paths are roots both before and
  after, so its format does not change

## Rollback

Its own commit. Reverting restores the empty-project behaviour for a store-backed
repo, which is wrong but is what ships today, and takes nothing else with it. The
revert also moves this change back into `openspec/changes/`.
