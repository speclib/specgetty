---
# specgetty-9xi9
title: remove the log view and functionality
status: completed
type: task
priority: normal
created_at: 2026-09-21T16:51:18Z
updated_at: 2026-09-21T17:56:08Z
---

The log panel costs a quarter of the focus ring, a key and a clause in five
capabilities, and holds nothing worth reading: one scan writes 35 lines into it,
30 of them per-project telemetry.

## OpenSpec change

`remove-the-log-panel` (`openspec/changes/remove-the-log-panel/`), validated
strict. Not yet implemented.

Decisions taken during exploration:

- log output is discarded while the interface runs. Deleting the redirection
  instead would send it to stderr, which writes over the frame
- `--debug` keeps logging exactly as it does, being the one place the output is
  useful and correct
- the glob bug is fixed here rather than left behind

## Bug found and folded in

`Walk` adds each include to the walk list and expands globs afterwards, so the
literal pattern `gh.*` is walked as a directory and always fails. That is the
source of the three ERROR lines the panel shows on every scan, and it makes
`spg --ignore_dir_errors=false` fail outright for any config using a glob,
which the shipped default config does. The panel is what made it survivable
enough to go unnoticed.

## Deferred

The three real error producers (a watcher that failed to start, an unreadable
scan directory, an unreadable store registry) become silent. They were already
buried under thirty telemetry lines in a panel nobody opens, so little changes
in practice, but auto-rescan dying without a word deserves a home. Recorded as
task 5.2, to become its own bean.

## Summary of Changes

Shipped as `remove-the-log-panel` (commit 0aacf49).

- the panel, the `l` key, `focusLog`, the two-view panel switch,
  `logPanelHeight`, `viewFocused` and the log branch of `halfPage` are gone.
  The focus ring is three states and the panel border is always lit
- log output is discarded while the interface runs, replacing the redirection
  into the panel rather than deleting it. `--debug` prints exactly what it did
- a glob include is no longer walked as a path of its own

## The bug this removed

`Walk` appended each include before expanding globs, so `gh.*` was walked as a
directory and always failed. Three ERROR lines per scan on this author's config,
and `--ignore_dir_errors=false` failed outright for any config using a glob.

Fixing it uncovered `Walk`'s error branch: the spurious failure was the only
thing exercising it. A plain include that does not exist now covers it, which is
the case the flag was written for.

Coverage floors raised: scanner 94.9 to 95.1, ui 84.5 to 85.0, total 86.8 to
87.2. Follow-up bean [[specgetty-7lc7]] for the three error producers that are
now silent.
