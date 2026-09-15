---
# specgetty-l5e6
title: keybindings for copying change-name and change-path to clipboard
status: completed
type: feature
priority: normal
created_at: 2026-09-15T20:19:06Z
updated_at: 2026-09-15T21:02:02Z
---

two seperate keybindings

## Planned (2026-09-15)

OpenSpec change `copy-change-name-and-path` (tinychange) holds the spec delta
and tasks. New capability `copy-to-clipboard`.

## Decisions

| Topic      | Decision                                                      |
| ---------- | ------------------------------------------------------------- |
| mechanism  | `atotto/clipboard`, already in the module graph via `bubbles`  |
| OSC 52     | rejected: it only wins over SSH, which does not matter here    |
| v2         | do not wait; with SSH out, v2's clipboard offers nothing extra |
| `y`        | copy the change name                                           |
| `Y`        | copy the full absolute path                                    |
| scope      | the change list; the project picker stays out                  |
| feedback   | a status line, cleared by the next keypress                    |

Verified on this machine: `wl-copy` and `wl-paste` are both present and round
trip, and `atotto`'s `init()` checks `WAYLAND_DISPLAY` before falling back to
xclip and xsel, neither of which is installed here.

## The thing most likely to be got wrong

The path must be built from `ChangeInfo.DirName`, never from `Name`. For an
archived change `Name` has the `YYYY-MM-DD-` prefix stripped for display, while
the directory on disk keeps it. `doExportChange` made exactly this mistake and
shipped broken from f2f6105 until 2026-09-15.

It is worse here than it was there: a wrong path on the clipboard fails later
still, in a shell, far away from the tool that produced it.

## Note on the feedback mechanism

specgetty has no way to report anything except a modal you dismiss with a key.
That weight is wrong for a copy, and silence is worse, because a failed copy
would look exactly like a successful one until the paste. The status line is
specified inside this capability because this is the change that needs it; if a
second caller appears it should be lifted out on its own.

## Summary of Changes

Shipped as openspec change `copy-change-name-and-path`, archived to
`openspec/changes/archive/2026-09-15-copy-change-name-and-path/`. Commit 8044d3c.
The delta created `openspec/specs/copy-to-clipboard/` with 4 requirements.

`y` copies the change name, `Y` copies the absolute path of its directory.

## The trap, guarded three ways

The path is built from `DirName`, never `Name`. Checked by reverting the line
and confirming two tests fail:

- `TestCopyPathOfAnArchivedChangeResolves` stats the copied path against a real
  temporary project, so it fails for any reason the path is wrong, not only this
  one
- `TestChangeDirPathUsesDirNameNotName` tests the rule directly, without the key
  handling in the way

Also checked against this repository: the archived path resolves to
`openspec/changes/archive/2026-03-31-fix-openspec-detection-false-positives`,
prefix intact.

## The status line

New, and the first thing in specgetty that reports a result without a modal.
It replaces the nav bar while set and is cleared by the next keystroke, so
there is no `tea.Tick` and no re-render loop. Failure uses the same line, which
is the point: a copy that did not happen must not look like one that did.

Kept general rather than copy-specific. If a second caller appears it should be
lifted into its own capability.

## Testing note

`writeClipboard` is a package-level function variable so tests can replace it.
The suite never touches the real clipboard: clobbering what the developer had
copied is a bad neighbour, and depending on `wl-copy` being installed would fail
on machines that are fine.

## Left for you

Running `spg` and pressing both keys, then pasting. The clipboard reaches the
terminal and the window manager, and no test here can see either.

Coverage: `src/ui` 78.8% to 78.9%, total 78.8% to 79.0%.
