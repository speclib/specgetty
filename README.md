# specgetty

Do you work with multiple projects that use OpenSpec for managing specifications?

Have you ever lost track of which OpenSpec projects exist on your local machine
or what state they're in?

`spg` (specgetty) is a text-mode UI tool to find and report the status of
OpenSpec projects on your local machine.

## Source-mode installation

```bash
go install github.com/mipmip/specgetty@master
```

## Configuration

Copy [config.yml](src/config.yml) to `~/.config/specgetty/config.yml` and edit to your needs.

The config path follows the XDG Base Directory Specification. If `$XDG_CONFIG_HOME` is set, the config is read from `$XDG_CONFIG_HOME/specgetty/config.yml`.

## Running

```bash
spg [ <directories...> ]
```

`spg` opens the OpenSpec project you are standing in. It resolves it from the
working directory and starts immediately, without scanning your disk. If there
is no project there, it offers the project picker.

| Flag                    | Effect                                              |
| ----------------------- | --------------------------------------------------- |
| `--view=single`         | Open the project at the working directory (default) |
| `--view=all`            | Open at the project picker                          |
| `--path <dir>`          | Open that project; implies `--view=single`          |
| `--change-fields=<a,b>` | Choose the change list columns                      |

If one/more directories are specified as `<directories>`, then this will override the
`scandirs.include` from your config file.

## UI

There are two levels, and the project picker opens over either of them.

```
   [ project picker ]   p opens it, esc closes it
           |            enter switches project
           v
   project view         the startup view; esc does nothing here
       | enter          tabs: changes | specs | config
       v
   change view          artifact sub-tabs: proposal | design | tasks | specs
```

Active and archived changes share one list. The `f` key cycles which of them it
shows, so there is no separate archive tab.

### Keys everywhere

| Key                        | Action                  |
| -------------------------- | ----------------------- |
| `j`/`k` or `<up>`/`<down>` | Move the cursor         |
| `<enter>`                  | Go one level deeper     |
| `<esc>`                    | Go one level back       |
| `p`                        | Open the project picker |
| `s`                        | Rescan the open project |
| `l`                        | Toggle the log panel    |
| `q` / `ctrl-C`             | Quit                    |

### Keys in a project

| Key                | Action                                           |
| ------------------ | ------------------------------------------------ |
| `<left>`/`<right>` | Switch tab                                       |
| `1` / `2` / `3`    | changes / specs / config                         |
| `<tab>`            | Switch focus between the panels (top level only) |

### Keys in the specs tab

The specs tab has two halves and `tab` moves the keyboard between them, and on
through the log panel when it is open.

| Key     | Action                                             |
| ------- | -------------------------------------------------- |
| `tab`   | Move the keyboard between the list and the content |
| `j`/`k` | Change spec, or scroll it, depending on the focus  |

The selected spec is highlighted while the list has the keys and dimmed while
the content has them, and the title shows a reading position only in the second
case. When the content has the keys it takes the same scroll keys as an open
change.

### Keys in the change list

| Key   | Action                                                |
| ----- | ----------------------------------------------------- |
| `/`   | Filter the list (see below)                           |
| `f`   | Cycle active / archived / both                        |
| `a`   | Archive the selected change                           |
| `d`   | Discard the selected change                           |
| `e`   | Export the selected change as a zip                   |

### Keys in an open change

| Key                     | Action                       |
| ----------------------- | ---------------------------- |
| `<left>`/`<right>`      | Switch artifact sub-tab      |
| `j`/`k` or arrows       | Scroll the document one row  |
| `<pgdn>`/`<pgup>`       | Scroll a full page           |
| `ctrl-f`/`ctrl-b`       | Scroll a full page           |
| `ctrl-d`/`ctrl-u`       | Scroll half a page           |
| `gg` / `G`              | Jump to the start or the end |

The panel title shows how far down the document you are. No percentage means
the whole thing fits, so there is nothing below. The config tab scrolls with the
same keys.

Sub-tabs stay inside the change: pressing `<right>` on the last one does not
spill over into the project tab bar.

## The project picker

Press `p` from anywhere. It lists every OpenSpec project it has found, with its
spec, change and task counts, and `enter` switches to the highlighted one. You
land on that project's change list.

| Key       | Action                             |
| --------- | ---------------------------------- |
| `<enter>` | Switch to the highlighted project  |
| `<esc>`   | Close, keeping the current project |
| `/`       | Filter the list                    |
| `r`       | Look for projects again            |
| `j`/`k`   | Move the cursor                    |

The filter uses the same grammar as the change list, with one addition: `:`
searches file paths as well as file contents, so you can look for a project by
a filename it contains.

### Projects that keep their specs in a store

OpenSpec lets a repo keep no specs and no changes of its own and name a store
instead, so that several repos can plan together:

```
  openspec/config.yaml          the store, registered on this machine
    schema: spec-driven
    store: nivis-tunnel  ---->  ~/openspec-stores/nivis-tunnel/openspec/
    context: |                    specs/
      ...                         changes/
```

`spg` follows that pointer. Such a repo opens on the store's specs and changes,
its header carries a `store` mark, and the config tab gains sub-tabs: the repo's
own configuration, the store's shared one, and a report on the store itself.

The picker lists the repos you work in. A registered store is not listed on its
own: it is where content lives rather than where work happens, and every repo
reading from it already stands for it. The `store` column names the store a row
reads from, so two repos sharing one are visibly sharing rather than looking
like duplicates. To open a store directly, run `spg` inside it or pass
`--path`.

Resolution reads local files only. Nothing here clones or fetches, and the
store's git state is reported from local refs, so an `ahead` count is against
the last upstream ref you fetched.

### Why `r` exists

Walking your disk to find projects is the slow part of specgetty, measured at
about four seconds cold for eighteen projects. So `spg` never does it at
startup, and the picker remembers what it found in
`~/.cache/specgetty/projects.yaml`.

That cache holds paths only. Every count you see in the picker is read fresh
from disk, so the numbers are never stale. What can go stale is the list itself:
a project you created since the last scan will not appear until you press `r`.
Editing `scandirs` in your config invalidates the cache on its own, and projects
you have deleted disappear without a rescan.

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

| Field      | Shows                                            |
| ---------- | ------------------------------------------------ |
| `name`     | The change name. Takes the leftover width        |
| `tasks`    | Task progress as `done/total`                    |
| `specs`    | How many specs the change touches                |
| `archived` | Whether the change is active or archived         |
| `date`     | The archive date, for archived changes           |

The default is `name,tasks,specs`. On a narrow terminal, columns are dropped
from the right rather than squeezing the name past readability.

## Which changes the list starts on

`f` cycles between active, archived and both. Where it starts, and where it
returns to when you switch project, comes from `change_mode` in the config file,
and `--change-mode` overrides that:

```bash
spg --change-mode=active+archived
```

| Mode              | Shows                         |
| ----------------- | ----------------------------- |
| `active`          | active changes only (default) |
| `archived`        | archived changes only         |
| `active+archived` | both, each row labelled       |

These are the same names the nav bar shows, so what you see is what you write.

## Development

```bash
make lint
```
