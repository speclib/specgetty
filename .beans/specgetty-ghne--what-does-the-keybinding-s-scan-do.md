---
# specgetty-ghne
title: what does the keybinding s (scan) do?
status: completed
type: task
priority: normal
created_at: 2026-09-21T16:56:34Z
updated_at: 2026-09-21T21:57:43Z
blocked_by:
    - specgetty-7lc7
---

## Answer

`s` runs `rescanCurrent`, which re-reads the open project from the directory its
resolution started in, so the `store:` declaration is followed again rather than
assumed. It is scoped: one project, never a walk of the scan directories. It is
also the only scan that reads git state.

## What it catches that the watcher does not

The watcher covers `<root>/openspec/**` and `<origin>/openspec/**`, and it does
add directories created after it started. Four things sit outside that:

1. the store registry, under the data directory
2. the store git state, read from `.git/`
3. a watcher that failed to start, which is `specgetty-7lc7`
4. a signal dropped by the non-blocking send into a channel of one

So the key is not a refresh. It is the manual recovery for four defects.

## OpenSpec change

`retire-the-scan-key`, validated strict. Closes 1, 2 and 4, then removes the
key. Blocked by `specgetty-7lc7`: `s` is the only recovery for a watcher that
died silently, so removing it before the failure is reported would leave a
failure mode with no workaround and no message.

`p` then `enter` on the open project stays, and re-resolves and re-reads exactly
as `s` did. It is the manual path that remains.

## Summary of Changes

Shipped as `retire-the-scan-key` (commit `7362391`).

The automatic path was made to cover what the key covered, and then the key
went.

- **The store registry is watched** for a project that declares a store,
  resolved or not. It lives outside every `openspec/` tree, so registering,
  removing or repointing a store changed resolution with no watched file
  changing. A project declaring no store watches nothing extra.
- **A change arriving during a scan is no longer dropped.** The send stays
  non-blocking, a blocking watcher being one that stops reading events; a
  pending flag causes one further scan once the current one lands. One however
  many arrived, a scan reading everything.
- **A store's git state is read on entering the properties tab.** Watching
  `.git/` was rejected: a fetch rewrites refs in bulk and would turn a quiet
  project into a rescan loop.
- **`s` and its four nav bar hints are gone.** `p` then `enter` on the open
  project re-resolves and re-reads it, which is what the key did.

## Shipped ahead of its prerequisite

The bean was `blocked_by specgetty-7lc7`, which reports a watcher that failed to
start, and 7lc7 is still open. Shipped anyway on an explicit instruction naming
this change.

The consequence, recorded so it is not rediscovered: a watcher that fails to
start is now silent AND has no manual recovery. Before this, `s` was the
workaround. `specgetty-7lc7` is what closes it, and it is worth more now than it
was this morning.

## Notes

- Two tests passed under a deliberately broken implementation and were
  strengthened: both now assert a further scan positively rather than only
  asserting nothing is left pending.
- Two existing watch-set tests asserted the old two-tree answer and were updated
  to expect the registry, which is behaviour this change deliberately alters.
- Coverage: src/ui 89.5% (floor 88.3), total 90.5%.
