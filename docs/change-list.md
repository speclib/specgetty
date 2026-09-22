# The change list

Active and archived changes in one grouped list: searching it, choosing its
columns, and the three actions a change offers.

- [Searching](#searching-the-change-list)
- [Columns](#change-list-columns)
- [How it is grouped](#how-the-change-list-is-grouped)
- [Exporting a change](#exporting-a-change)

![The grouped change list, and a change opened at its artifacts](../demo/recordings/changes.gif)

## Searching the change list

Press `/` to filter. The list narrows as you type, and `<up>`/`<down>` (or
`ctrl-p`/`ctrl-n`) still move the cursor while the prompt has focus, so you can
type, arrow down and press `<enter>` to open a change without leaving the
prompt. `<esc>` clears the filter.

| Query      | Matches                                                  |
| ---------- | -------------------------------------------------------- |
| `expzip`   | change names, fuzzy, best match first                    |
| `'export`  | change names, literal substring                          |
| `:export`  | the text inside proposal, design, tasks and spec files   |

A lowercase query ignores case. A query containing any uppercase letter is
matched case-sensitively, the same rule `rg` and `fzf` use.

When a change matches on file contents rather than its name, the row names the
files that matched, so you can see why it is in the list.

The filter survives opening a change and returning, and survives the automatic
rescan that fires when a file changes on disk. It clears when you select a
different project.

## Change list columns

The change list is a table. Which columns it shows, and in what order, comes
from `change_fields` in the config file, and `--change-fields` overrides that:

```bash
spg --change-fields=name,tasks,specs,archived
```

| Field      | Shows                                                                                          |
|:-----------|:-----------------------------------------------------------------------------------------------|
| `name`     | The change name. Takes the leftover width                                                      |
| `tasks`    | Task progress as `done/total`                                                                  |
| `specs`    | How many specs the change touches                                                              |
| `archived` | Whether the change is active or archived. Redundant against the group header, so not a default |
| `date`     | The archive date, for archived changes                                                         |

The default is `name,tasks,specs,date`. On a narrow terminal, columns are
dropped from the right rather than squeezing the name past readability.

## How the change list is grouped

Active and archived changes are grouped in one list, with the active group first
and each group counting its rows:

```
 name                              tasks    specs  archived
 ACTIVE (1)
 group-the-change-list             1/39     3
 ARCHIVED (36)
 list-the-repos-not-the-stores     28/28    3      2026-09-21
 properties-tab                    47/47    5      2026-09-21
```

Both groups are always shown, with a count of zero when one is empty, so
"nothing in flight" is an answer rather than an absence. The archive date is a
column, blank on an active change, and the archived group is ordered newest
first while the active group is ordered by name. Neither order is selectable.

Earlier versions put a filter in front of the list, cycled with `f` and
configured with `change_mode`. The group a row sits in says what the filter
said, so the key, the setting and the `--change-mode` option are all gone. A
configuration still carrying `change_mode` is reported at startup.

## Exporting a change

![Exporting a change as a zip](../demo/recordings/export.gif)

`e` asks where to put the zip. The prompt opens on `export_dir` from your
config, or your home directory when that is unset, and `tab` completes a path
against the filesystem:

You type a directory; the filename is generated from the change name and the
date. Editing the directory redirects that one export and leaves your config
alone. A directory that does not exist is refused rather than created, and an
existing file is confirmed before it is replaced.
