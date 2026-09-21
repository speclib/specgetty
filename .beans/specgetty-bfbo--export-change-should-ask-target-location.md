---
# specgetty-bfbo
title: export change should ask target location
status: completed
type: task
priority: low
created_at: 2026-09-21T17:04:59Z
updated_at: 2026-09-21T18:08:39Z
---

currently it only allows saving in home

## OpenSpec change

`choose-the-export-directory` (`openspec/changes/choose-the-export-directory/`),
validated strict. Not yet implemented.

Decisions taken during exploration:

- both: an `export_dir` setting, and a prompt that opens on it and can be
  edited for one export
- a directory is typed, never a filename. The generated name strips an archived
  change's date prefix and stamps the export date, and retyping would lose both
- `tab` completes against the filesystem, using the text input's own suggestion
  support rather than new machinery
- a directory that does not exist is refused with a message, leaving the typed
  text in place. Nothing is created
- an existing file is confirmed before replacement. The current spec says
  overwrite silently, which was defensible for a generated name in your own home
  and is not once the path is typed
- `edit_command` is deleted, and reported through the retired-key mechanism

## Found while exploring

`edit_command: code %WORKING_DIRECTORY` ships in `src/config.yml`, parses
into `Config.EditCommand`, and nothing anywhere reads that field. It is not in
the README either.

## Deferred

The retired-setting requirement lives in `change-list-view` because
`change_mode` was introduced there. It now serves a second key with nothing to
do with the change list. Generalised in place; its home is task 6.2.

## Summary of Changes

Shipped as `choose-the-export-directory` (commit 7220428).

- `export_dir` in the config, and a prompt that opens on it. Editing redirects
  one export and leaves the config alone
- a directory is typed; the filename stays generated
- `tab` completes against the filesystem via the text input's own suggestions
- a missing directory is refused with a message, the typed text stays, nothing
  is created
- an existing file is confirmed; declining returns to the prompt with the
  directory still in it
- `edit_command` deleted and retired through the existing mechanism

## Found while verifying

Task 2.7's frame assertion caught a violation older than this change: modal
widths were fixed, so the export prompt asked for 66 columns and the startup
question for 62, both wider than the 60-column minimum `modal-presentation`
already required them to respect. Every modal is now clamped to the terminal,
with a test over all eight states.

Coverage floors raised: scanner 95.1 to 95.4, ui 85.3 to 85.8, total 87.5 to
88.0.
