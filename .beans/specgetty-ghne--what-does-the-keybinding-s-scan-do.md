---
# specgetty-ghne
title: what does the keybinding s (scan) do?
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T16:56:34Z
updated_at: 2026-09-21T21:51:59Z
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
