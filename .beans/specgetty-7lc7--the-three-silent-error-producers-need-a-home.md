---
# specgetty-7lc7
title: the three silent error producers need a home
status: todo
type: task
priority: normal
created_at: 2026-09-21T17:53:56Z
updated_at: 2026-09-21T17:53:56Z
---

Deferred from `remove-the-log-panel`.

Three things log an error and nothing shows it now that the panel is gone:

- `ui.go` a watcher that failed to start, which means auto-rescan is dead and
  the list silently stops updating
- `find.go` a scan directory that cannot be read
- `find.go` a store registry that cannot be read

They were already as good as invisible, buried under thirty telemetry lines in a
panel nobody opened, so little changed in practice. But a watcher dying without
a word is the sharpest of the three: the UI simply stops responding to the
filesystem.

`statusMsg` (the transient nav bar line) and the error modal both already exist
as homes.
