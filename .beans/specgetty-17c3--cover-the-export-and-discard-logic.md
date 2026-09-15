---
# specgetty-17c3
title: cover the export and discard logic
status: todo
type: task
priority: normal
created_at: 2026-09-15T19:12:17Z
updated_at: 2026-09-15T19:12:17Z
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
