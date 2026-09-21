## Context

See proposal.md for why. What matters here is the order, because three of the
four gaps have to be closed before the key can go and one of them is not in this
change.

```
  specgetty-7lc7  ──▶  retire-the-scan-key
  reports a watcher      closes gaps 1, 2 and 4, then removes the key
  that failed to start
```

`s` is the only recovery for a watcher that died silently. Removing it first
leaves a failure mode with no workaround and no message, which is strictly worse
than the key existing.

## Goals / Non-Goals

Goals: make the automatic path cover what the key covered; remove the key.

Non-Goals: reporting a failed watcher, which is `specgetty-7lc7`. Removing the
manual re-read, which stays on the picker. Any change to how a scan itself
works.

## Decisions

### The registry is watched only for a project that declares a store

The registry decides resolution only for a project whose configuration names a
store. For every other project it is inert: nothing in it can change what that
project reads. Watching it unconditionally would add an inotify watch per
session for no behaviour, and on Linux inotify watches are the resource this
application is most likely to run out of.

A project whose declaration could not be resolved still watches it, because
registering the missing store is exactly the event worth noticing, and that is
the case a user is most likely to be in while looking at the screen.

### A pending flag, not a bigger buffer

The dropped signal comes from a non-blocking send into a channel of one:

```go
select {
case w.events <- struct{}{}:
default:          // a signal arriving while one is unconsumed is discarded
}
```

The non-blocking send is right. A watcher that blocks on a busy consumer stops
reading filesystem events, and then it misses changes rather than merely
delaying them.

Enlarging the buffer was considered and rejected. It converts one lost update
into N queued scans of the same project, which is worse: the scans are redundant
by construction, since each reads the whole project anyway.

The fix is a flag. A signal that cannot be delivered sets "something happened",
and the consumer scans again when its current scan lands and the flag is set.
One further scan however many signals were dropped, which is the correct number,
because a scan reads everything.

Where the flag lives is an implementation question rather than a design one. In
the watcher it keeps the concern with the thing that dropped the signal; in the
model it is visible next to `m.scanning`, which is the state it has to coordinate
with. The tasks take the second, because the coordination is the hard part and
splitting it across a package boundary is what makes such bugs.

### Git state moves to tab entry rather than to the watcher

Watching `.git/` was considered. It is the obvious symmetry with watching the
registry and it is a bad idea: `.git/` churns constantly, a fetch rewrites refs
in bulk, and index operations touch files many times a second. It would turn a
quiet project into a rescan loop, and a rescan re-reads every spec and change.

Reading on tab entry costs nothing until the tab is opened and is exactly
current when it is. `schema-inspection` already established the pattern and the
argument for it, including that the cost is a subprocess and that nothing is
read twice for the same project.

The cost of the difference is that git state can be stale while the tab is open
and something commits elsewhere. That is one subprocess per tab entry against a
rescan loop, and the state was already stale between scans before this change.

### The picker keeps the manual path

`choosePickerProject` already calls `doScanSingle` on the project it opens, and
its comment says why: "Opening a project is the moment the rest is worth
reading, which for a store means its git state." That is `s` under another key,
and nothing here touches it.

This matters beyond mechanics. A manual re-read is a trust affordance, and
removing every way to force one would be a worse outcome than the key being
redundant. What changes is that it stops being the answer to a defect and goes
back to being what it is: opening a project.

## Risks / Trade-offs

- A user with `s` in muscle memory presses it and nothing happens. → The nav bar
  stops advertising it in the same commit, and the README names the picker. A
  key that silently does nothing is the normal outcome for every unbound key in
  the application.
- The registry watch adds a watched directory for store-backed projects. → One
  directory, and only for projects that declare a store.
- Git state on tab entry is a subprocess on a keypress. → `schema-inspection`
  already spends about a second there for schema definitions, and this is
  smaller.
- Closing gap 4 is a concurrency fix, which is where subtle bugs live. → It is
  one flag with one reader, and the scenario that exercises it is a write during
  a scan, which is testable without timing by driving the messages directly.

## Migration Plan

None for data. For the user, the README gains the picker where it named the key.
