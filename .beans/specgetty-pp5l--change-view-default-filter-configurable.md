---
# specgetty-pp5l
title: change view default filter configurable
status: todo
type: feature
priority: normal
created_at: 2026-09-15T20:28:24Z
updated_at: 2026-09-17T17:33:52Z
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
