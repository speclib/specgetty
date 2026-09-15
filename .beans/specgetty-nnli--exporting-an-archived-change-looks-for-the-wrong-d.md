---
# specgetty-nnli
title: exporting an archived change looks for the wrong directory
status: todo
type: bug
priority: normal
created_at: 2026-09-15T19:55:07Z
updated_at: 2026-09-15T19:55:07Z
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
