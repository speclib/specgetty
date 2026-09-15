---
# specgetty-nnli
title: exporting an archived change looks for the wrong directory
status: completed
type: bug
priority: normal
created_at: 2026-09-15T19:55:07Z
updated_at: 2026-09-15T20:07:40Z
---

## Found while planning coverage (2026-09-15)

Exporting an **archived** change has never worked.

`src/scanner/scan.go` stores the archived name with its date prefix already
stripped:

    displayName = dirName[11:]      // "2026-04-02-my-feature" -> "my-feature"
    ci := parseChangeDir(filepath.Join(archiveDir, dirName), displayName)

`doExportChange` then builds the source path from that stripped name:

    srcDir = filepath.Join(projectPath, "openspec", "changes", "archive", changeName)

which is a directory that does not exist. Verified on this repository:
`openspec/changes/archive/fix-openspec-detection-false-positives` is absent;
the real directory is `2026-03-31-fix-openspec-detection-false-positives`.

Present since the feature shipped in `f2f6105`: that commit already read
`archived[m.archiveCursor].Name`, which was the stripped name then too.

`exportSemanticName` exists only to strip a date prefix from the name, so it is
a no-op in practice and is the fossil of the assumption that caused the bug.

## Not a spec problem

`openspec/specs/export-change/spec.md` already requires the correct behaviour:

> **WHEN** exporting an archived change with directory name `2026-04-02-my-feature`
> **THEN** the zip SHALL contain a root folder `my-feature/` (date prefix
> stripped) with all files and subdirectories from
> `openspec/changes/archive/2026-04-02-my-feature/`

The spec is right and the code does not match it, so the fix carries no spec
delta.

## Where it is tracked

Task group 1 of the openspec change `cover-scanner-and-export`, together with
the test that demonstrates it (task 2.3, written to fail first).

Nobody noticed for the same reason the change exists: the export logic had no
tests at all.

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
