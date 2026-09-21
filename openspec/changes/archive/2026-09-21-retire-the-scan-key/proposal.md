## Why

Bean `specgetty-ghne`, which asks what `s` does. Tracing it turned into a better
question: what does it do that the filesystem watcher does not?

```
  s ─▶ rescanCurrent ─▶ doScanSingle(startDirOf(key)) ─▶ ScanResolved
                                                          ├ ResolveRoot      config + registry
                                                          ├ ListOpenSpecContents
                                                          ├ ParseProjectInfo
                                                          └ ReadStoreGit     withGit: true
```

The watcher covers `<root>/openspec/**` and `<origin>/openspec/**`, and it
handles directories created after it started, which is the gap you would expect
and is already closed. What is left is four things, all outside those trees by
construction:

```
  1  the store registry       under the data directory, never watched. A store
                              registered, removed or repointed while specgetty
                              runs changes where content is read from
  2  the store's git state    read from `.git/`, not from `openspec/`. A commit
                              or fetch in another terminal leaves the ahead and
                              behind counts stale on the properties tab
  3  a watcher that never     `startWatcher` is allowed to fail. Auto-rescan is
     started                  then dead and nothing says so
  4  a dropped signal         `events` holds one, and the send is
                              `select { case ...: default: }`. A burst settling
                              while the previous signal is unconsumed is
                              discarded
```

So `s` is not a refresh key. It is the manual recovery for four defects, three of
which are fixable and one of which is already bean `specgetty-7lc7`. A tool that
watches the filesystem and also asks the user to press a key when it looks wrong
is telling on itself.

## What Changes

The automatic path is made to cover all four, and then the key goes.

- The store registry is watched for a project that declares a store, so
  registering, removing or repointing one is noticed the way editing the
  declaration already is. A project that declares no store watches nothing extra:
  the registry cannot change what it resolves to.

- A filesystem signal arriving while a scan is in flight is remembered rather
  than dropped. The `default:` arm exists so the watcher never blocks, which is
  right; what is wrong is that the signal is then lost. A pending flag, scanned
  again when the current scan lands, keeps the non-blocking send and the update.

- The store's git state is read when the properties tab is entered, not only
  when the project is scanned. This is the rule `schema-inspection` already
  follows for schema definitions: nothing is read until the tab is asked for.
  Git state is the other thing on that tab that costs subprocesses, and it is
  the only part of a project that changes without `openspec/` changing.

- **REMOVED**: the `s` key, its four nav bar hints, and its README row.

- `p` and `enter` on the open project stay as they are, and remain a full
  re-resolution. Anyone who wants to force a read still can; it is two keys
  rather than one, and it is no longer the fix for anything.

Not in scope, and the reason:

- **Reporting a watcher that failed to start.** That is bean `specgetty-7lc7`,
  written before this one, which calls it "the sharpest of the three: the UI
  simply stops responding to the filesystem". It is a prerequisite rather than a
  part: removing `s` removes the workaround for a failure nobody is told about,
  so the telling has to exist first.

- **Dropping the manual rescan entirely.** It is not dropped, it moves. The
  picker already rescans the project it opens, and that path is unchanged.

## Capabilities

### Modified Capabilities
- `project-view`: the rescan key is removed
- `fs-watch`: the registry is watched where it can change resolution, and a
  signal arriving during a scan is not lost
- `config-tab-display`: the store's git state is read on entering the tab

## Impact

- `src/ui/ui.go`: the `s` case, four nav bar hints, and a pending-scan flag
- `src/ui/watchdirs`: `watchDirs` gains the registry's directory for a project
  that declares a store
- `src/watcher/watcher.go`: the dropped-signal fix, if it belongs there rather
  than in the model
- `src/ui/properties.go`: git state read on entering the tab
- `README.md`: the key row, and the sentence in the editor section that tells
  the reader their edit appears "without pressing `s`"

## Rollback

Its own commit. Reverting restores the key and the four gaps it was covering,
which is what ships today.
