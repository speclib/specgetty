---
# specgetty-4wjm
title: remove the unused unwrap helper
status: completed
type: task
priority: low
created_at: 2026-09-15T19:12:17Z
updated_at: 2026-09-15T20:07:40Z
---

`unwrap` in `src/ui/table.go` has zero callers.

Added during the generics extraction in the `project-picker` change, on the
assumption the picker would need to strip match reasons off a filtered row. It
never did: the picker works with `filtered[projectRow]` throughout, same as the
change list.

Dead on arrival, and it shows up as a 0% function in the coverage report, which
is how it was noticed.

Its counterpart `wrap` is used and stays.

## Planned

Covered by openspec change `cover-scanner-and-export` (tinychange). Planning
artifacts written 2026-09-15; run /mip:tinychange-apply to implement.

## Summary of Changes

Applied as openspec change `cover-scanner-and-export` (tinychange, skip_specs),
archived to `openspec/changes/archive/2026-09-15-cover-scanner-and-export/`.

### The bug, demonstrated before it was fixed

The regression test was written first and failed with exactly the predicted
message:

    Source not found: .../openspec/changes/archive/my-feature

`ChangeInfo` now carries `DirName`, the directory the change was read from,
which for an archived change keeps its `YYYY-MM-DD-` prefix while `Name` has it
stripped for display. `doExportChange` takes both and can no longer confuse
them: the signature is `(projectPath, dirName, semanticName, isArchived)`.

`exportSemanticName` was deleted along with the `regexp` import it needed. Once
the source path comes from `DirName`, the display name is already the semantic
name, so stripping it again was a no-op. It only ever existed because of the
assumption that caused the bug.

Verified against this repository rather than a fixture: exporting
`2026-03-31-fix-openspec-detection-false-positives` produces a zip rooted at
`fix-openspec-detection-false-positives/` with its five real files.

### Coverage

| Package | Before | After | Target |
| ------- | ------ | ----- | ------ |
| scanner | 53.3%  | 91.9% | 80% (met) |
| ui      | 68.1%  | 72.4% | 80% (short) |
| total   | 64.4%  | 73.8% | 70% (met) |

Floors raised to match. Overall and scanner now clear the gate the ship script
names; `ui` is the only one still short.

### Also

`unwrap` deleted from `src/ui/table.go`; it never had a caller.

### Unverified

`Walk` calls `log.Fatal` when a glob's parent directory cannot be read, which
would kill the process rather than report an error. Left alone: it is
pre-existing, out of scope here, and testing it would end the test binary.
`skip` panics on an empty string in the exclude list (`f[0:1]`), also
pre-existing and untested for the same reason it is unlikely: an empty exclude
entry has to be written by hand.
