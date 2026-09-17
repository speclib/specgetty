---
# specgetty-pp5l
title: change view default filter configurable
status: completed
type: feature
priority: normal
created_at: 2026-09-15T20:28:24Z
updated_at: 2026-09-17T17:48:06Z
---


e.g. : mode: archived+active

## Planned (2026-09-17)

OpenSpec change `configurable-change-mode` (tinychange) holds the spec deltas
and tasks. `change-list-view` and `changes-tab` modified.

## Decisions

| Topic           | Decision                                              |
| --------------- | ----------------------------------------------------- |
| the word        | `active`, not `open`                                   |
| rename          | bundled into this change                               |
| config key      | `change_mode`, beside `change_fields`                  |
| CLI flag        | `--change-mode`, flag over config over default         |
| values          | `active`, `archived`, `active+archived`                |
| unknown value   | error naming the value and the valid ones              |
| project switch  | snap back to the configured default                    |

## Why the rename belongs here

The three states are spelled two ways today and both reach the screen: the nav
bar says `f mode:open` while the empty state says "No active changes". The spec
contradicts itself inside one requirement, titled "Active and archived changes
share one list" and then describing "open only" mode. The README does it across
three lines.

"active" wins because `openspec` itself says active, the scanner already says
`ActiveChanges` with only the UI layer having invented "open", and "open" is
overloaded: the detail header currently reads `my-change (open)` about a change
that is open on screen at that moment.

The rename is bundled because a config value is a contract. Renaming a Go
constant is free; renaming a value in someone's config file is not. This is the
last moment the words are cheap to change.

## Principle worth keeping

The config value is exactly the string the nav bar shows. What you see is what
you write. That kills this class of drift rather than correcting one instance
of it.

## Surface

24 uses of `modeOpen` across five files, three on-screen strings
(`listModeNames`, the state column, the open-change header), two specs and three
README lines. All string-level; no logic changes beyond reading the default from
configuration.

## Summary of Changes

Shipped as `configurable-change-mode`, archived to
`openspec/changes/archive/2026-09-17-configurable-change-mode/`. Commit f96e09d.
26 tasks. `change-list-view` and `changes-tab` modified.

`change_mode` in the config, or `--change-mode`, sets which changes the list
starts on and returns to on a project switch. Values are `active`, `archived`
and `active+archived`, which are exactly the strings the nav bar shows.

## The rename it was really about

The three states were spelled two ways and both reached the screen: the nav bar
said `mode:open` while the empty state said "No active changes". 24 uses of
`modeOpen`, three on-screen strings, two specs and three README lines now all
say active.

"open" survives only in its other meaning, the change you have drilled into,
which is why it had to go: `my-change (open)` was being printed about a change
that was open on screen at that moment.

## Two things worth keeping

**A test that would have passed either way.** The project-switch test uses
`modeBoth` as the default deliberately. With `active` it could not distinguish
"returned to the configured default" from "reset to the old hardcoded value".
Reverting `m.listMode = m.defaultMode` to `modeActive` fails it, which was
checked.

**The nav bar test now derives its expectation** from `listModeNames` rather
than spelling the word out, so the next rename cannot leave it asserting a word
the UI no longer uses.

## Coverage, and why it moved so much

The ratchet caught `src` dropping from 29.4% to 28.2%: the new flag added
uncovered statements inside `main()`, which is untestable as a whole. Rather
than pad the number, `expandScanDirs` and `resolveStartupPath` were extracted
out of `main()` and tested. `src` went to 42.3%, and `main()` is shorter.

Totals: `src` 29.4% to 42.3%, `ui` 79.6% to 79.7%, overall 79.4% to 79.9%.

## Left for a terminal

Running `spg --change-mode=active+archived` and confirming the list opens with
both and the nav bar reads `mode:active+archived`. The flag's rejection path was
checked from the command line: `--change-mode=open` now reports
`unknown mode "open" ... valid modes are: active, archived, active+archived`.
