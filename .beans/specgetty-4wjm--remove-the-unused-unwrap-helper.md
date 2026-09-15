---
# specgetty-4wjm
title: remove the unused unwrap helper
status: todo
type: task
priority: low
created_at: 2026-09-15T19:12:17Z
updated_at: 2026-09-15T19:12:17Z
---

`unwrap` in `src/ui/table.go` has zero callers.

Added during the generics extraction in the `project-picker` change, on the
assumption the picker would need to strip match reasons off a filtered row. It
never did: the picker works with `filtered[projectRow]` throughout, same as the
change list.

Dead on arrival, and it shows up as a 0% function in the coverage report, which
is how it was noticed.

Its counterpart `wrap` is used and stays.
