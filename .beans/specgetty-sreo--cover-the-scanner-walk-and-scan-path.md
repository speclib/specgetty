---
# specgetty-sreo
title: cover the scanner walk and scan path
status: completed
type: task
priority: normal
created_at: 2026-09-15T19:12:17Z
updated_at: 2026-09-15T20:07:40Z
---

The scan path is the least covered core code in the repo: 53.3% against the
80% the ship gate names for core packages. Every zero is a function that
touches the filesystem, and the pattern for testing those already exists.

## Uncovered

| Function                        | Coverage |
| ------------------------------- | -------- |
| `Walk` (find.go:116)            | 0%       |
| `walkone` (find.go:34)          | 0%       |
| `ListOpenSpecContents`          | 0%       |
| `ScanPaths` (added 2026-09-15)  | 0%       |
| `Scan`                          | 0%       |
| `DumpConfig`                    | 0%       |

## The pattern to copy

`src/ui/loadprojects_test.go` builds real OpenSpec projects under `t.TempDir()`
and runs a full walk over them. Its `makeProject` helper is the thing to lift.

Watch for this: `isValidOpenSpecDir` requires a `config.yaml` or `project.md`
marker alongside `specs/` or `changes/`. A bare specs tree is not detected as a
project, which cost time before the guard was read.

## Worth covering deliberately

- `exclude` matching, which has two modes: a leading `/` compares the full path,
  otherwise only the basename (`find.go:14`)
- `followsymlinks` both ways
- A directory that errors mid-walk, since `ErrorCallback` decides between
  `SkipNode` and `Halt`
- `ScanPaths` against a path that has since been deleted

## Why it matters now

`specgetty-p4e2` (bubbletea v2) is deliberately waiting for the 70/80 gate,
because that migration is broad and mechanical and the suite is what makes it
safe.

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
