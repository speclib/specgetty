---
# specgetty-l5e6
title: keybindings for copying change-name and change-path to clipboard
status: draft
type: feature
priority: normal
created_at: 2026-09-15T20:19:06Z
updated_at: 2026-09-15T20:49:52Z
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
