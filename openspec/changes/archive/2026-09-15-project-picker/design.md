## Context

`specgetty-jdif` asks for the zoom vocabulary to go, the single project view to
become the startup view, and the project list to become a cached picker popup.
`specgetty-edee` and `specgetty-oarh` are scrapped into it.

The previous change, `change-list-full-view`, deliberately said nothing about
what sits above the change list, so that this one could define it without
contradicting an archived spec. That held: the change list and change specs need
no modification here.

### What the measurements say

Taken on the real configuration (`$HOME/c*`, `$HOME/tc*`, `$HOME/gh.*`,
`$HOME/mipnix`, `$HOME/secondbrain`), 18 projects found:

| Step                                        | Cost           |
| ------------------------------------------- | -------------- |
| Directory walk, cold                        | ~3.9s          |
| Directory walk, warm                        | 697ms          |
| Listing openspec contents (all 18)          | 55ms           |
| Reading every artifact (3522 files, 18MB)   | 101ms          |
| Resolving the project at the cwd            | ~1ms           |

Total wall-clock matches the walk almost exactly, because the consumer keeps up
with the producer. The walk is the entire cost.

Note that `ProjectStatus.ScanTime` and the logged `scanDuration` measure only
`ListOpenSpecContents`, not `ParseProjectInfo`. Anyone reading those numbers as
"the parse cost" will be misled, which is why the artifact read above was
measured separately.

## Goals / Non-Goals

**Goals**

- Startup that costs a millisecond instead of seconds
- One project list, reachable from anywhere, that remembers what it found
- Retire the zoom vocabulary and the split layout it belonged to

**Non-Goals**

- OpenSpec stores. The registry at
  `~/.local/share/openspec/stores/registry.yaml` makes `--view=store` cheap to
  add later, and a store root already contains an `openspec/` directory so the
  scanner finds it today as an ordinary project. Deliberately out of scope.
- Configurable picker columns. The change list has `--change-fields` because
  columns were asked for; nobody has asked to configure the picker's. Hard-code
  them until a second request justifies the generalisation.
- Reviving the split layout. It may return one day; git remembers it.

## Decisions

### The picker is an overlay, not a level

```
                    ┌─────────────────────────────────┐
          p ───────▶│  project picker                 │
                    │  table, / filter, r refresh     │
                    │  esc closes, enter switches     │
                    └──────────────┬──────────────────┘
                                   │ always lands here
                                   ▼
   level 0   project view   ────────────────  esc does nothing: this is the floor
                │  enter
                ▼
   level 1   change view                      esc returns to the change list
```

Depth shrinks from three levels to two. The picker is orthogonal to depth: it
opens from either level and always returns to the project view on the changes
tab, because an open change cannot outlive the project it belongs to.

**Alternative rejected**: keeping the project list as level 0 and adding the
picker as well. That leaves two ways to choose a project and keeps the split
layout alive for one of them.

### Cache paths, never statistics

The walk is 697ms to 3.9s and discovers something that changes when you create a
project. Reading every artifact is 101ms and produces numbers that change every
time you save a file. So the cache holds paths only:

```yaml
# $XDG_CACHE_HOME/specgetty/projects.yaml
version: 1
scanned_at: 2026-09-15T18:04:03Z
scandirs:
  include: [...]      # the question these paths answer
  exclude: [...]
paths:
  - /home/pim/gh.speclib/specgetty
```

`scandirs` is recorded because editing `include:` in the config means the cached
list answers a question you are no longer asking. A mismatch invalidates.

Deleted projects need no walk to notice: `stat` each cached path on load, which
is microseconds, and drop the missing ones. New projects appear only after an
explicit refresh, which is exactly the job `r` does.

Because statistics are never cached, every number the picker shows is current as
of the moment it opened. Opening costs roughly 100 to 250ms for 18 projects,
which is why rows are parsed eagerly rather than lazily. If that ever becomes
slow, parsing only the rows on screen is the escape hatch; it is not worth the
complexity now.

### Three key-intercepting modes must be ordered

`Update()` already intercepts keys for the confirm modals, before the search
prompt check. The picker is a third interceptor with its own cursor, its own
prompt and its own table. The order is:

```
  1. confirm modals   (archive / discard / export)   modal, blocking
  2. picker           (when open)                    modal, blocking
  3. search prompt    (when focused)                 captures runes only
  4. ordinary keys
```

A confirm modal outranks the picker because it is asking a question that must be
answered. The picker outranks the change list search because it is an overlay
drawn on top of it.

### There is always a base view, sometimes empty

`--view=all` opens the picker with nothing behind it, and declining the
"no project found" prompt leaves nothing to show. Rather than make `esc`
sometimes quit the application, which is a surprising key to overload, the
project view has an empty state that names the picker key.

This keeps one uniform model: a base view always exists, and the picker is
always dismissible.

### Query grammar is reused, not reinvented

The picker filters with the grammar `change-search` already defines: bare terms
fuzzy on the name, `'` literal on the name, `:` literal in contents. For the
picker, `:` searches file paths as well as file contents, because both are
already in memory after a parse and a fourth sigil would not earn its place.

The matcher parts that are already row-agnostic (`parseQuery`, smart case,
`fitCell`, `layoutFields`) are reused directly. The parts bound to `changeRow`
(`filterRows`, `knownFields`, `renderChangeTable`, `indexOfKey`) need a seam. A
row needs three things: a stable key, searchable text, and field accessors.

Generalising from one example is speculation; the picker is the second consumer,
which makes it evidence. Doing the extraction here rather than earlier was
deliberate.

## Risks / Trade-offs

- **The watcher gains a new way to leak.** It currently starts on entering a
  project and stops on leaving. With the project view as startup it runs
  essentially always, and switching projects through the picker must stop the
  old one before starting the new one. The original inotify design already
  listed a goroutine leak as its main risk, and this adds a path to it.
- **First run is still slow.** No cache, no project at the working directory,
  means a prompt and then the full walk. Only the first run pays it.
- **A stale cache hides new projects until refreshed.** That is the deliberate
  trade for an instant picker, and `r` is the remedy. The alternative, walking on
  every open, is the cost this change exists to remove.
- **Removing `--zoom` breaks muscle memory and any shell alias.** Accepted: it
  names a concept that no longer exists, and leaving an alias behind would keep
  the vocabulary alive in `--help`.

## Open Questions

- `overview-tab` describes a tab replaced by the persistent header back in
  `overview-as-persistent-header`, archived 2026-04-20. It is stale, but it is
  unrelated to this change and folding unrelated spec debt into a feature change
  would be tidying, not scoping. Left for its own cleanup.
