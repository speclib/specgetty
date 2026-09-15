---
# specgetty-sreo
title: cover the scanner walk and scan path
status: todo
type: task
priority: normal
created_at: 2026-09-15T19:12:17Z
updated_at: 2026-09-15T19:55:07Z
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
