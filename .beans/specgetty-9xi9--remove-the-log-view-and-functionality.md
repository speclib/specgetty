---
# specgetty-9xi9
title: remove the log view and functionality
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T16:51:18Z
updated_at: 2026-09-21T17:49:24Z
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
