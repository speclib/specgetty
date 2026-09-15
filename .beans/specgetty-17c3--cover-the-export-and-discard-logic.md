---
# specgetty-17c3
title: cover the export and discard logic
status: completed
type: task
priority: normal
created_at: 2026-09-15T19:12:17Z
updated_at: 2026-09-15T20:07:40Z
---

`export-change` was archived as a capability on 2026-09-15 with five
requirements, and its logic has no tests at all.

## Uncovered

| Function             | Coverage | Shape                          |
| -------------------- | -------- | ------------------------------ |
| `exportSemanticName` | 0%       | pure, string in and string out |
| `exportDestPath`     | 0%       | pure, depends on $HOME and today's date |
| `doExportChange`     | 0%       | writes a zip to disk           |
| `doDiscardChange`    | 0%       | moves a directory              |
| `doArchiveChange`    | 11%      | shells out to the openspec CLI |

The first two are pure functions and should have been tested when they were
written. They carry the whole naming rule, including stripping the `YYYY-MM-DD-`
prefix from an archived change so the zip is not named twice over.

## Reachable

`doExportChange` and `doDiscardChange` both operate on directories and can be
pointed at a `t.TempDir()`. Reading the zip back and asserting its entries is
the real test: that the internal structure is preserved and the archived date
prefix is stripped from the root folder inside it.

`doArchiveChange` is the awkward one, since it runs the `openspec` binary. It
already reports a clean error when the binary is missing, and that branch is
worth covering even if the success path is not.

## Requirements to test against

See `openspec/specs/export-change/spec.md`: export keybinding, confirmation
modal, zip creation, result feedback, nav bar hint. The keybinding and modal
behaviour are covered by the UI tests already; the zip creation and naming are
not.

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
