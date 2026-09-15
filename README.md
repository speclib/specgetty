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

If one/more directories are specified as `<directories>`, then this will override the
`scandirs.include` from your config file.

## UI

The UI has three levels. `<enter>` goes one level deeper, `<esc>` comes back.

```
  projects            the project list beside a detail panel
     | enter
  project             tab bar: changes | specs | config
     | enter          the changes tab is a full-width table
  change              artifact sub-tabs: proposal | design | tasks | specs
```

Active and archived changes share one list. The `f` key cycles which of them it
shows, so there is no separate archive tab.

### Keys everywhere

| Key                        | Action                                    |
| -------------------------- | ----------------------------------------- |
| `j`/`k` or `<up>`/`<down>` | Move the cursor                           |
| `<enter>`                  | Go one level deeper                       |
| `<esc>`                    | Go one level back                         |
| `s`                        | Rescan                                    |
| `l`                        | Toggle the log panel                      |
| `gg` / `G`                 | Jump to first / last                      |
| `q` / `ctrl-C`             | Quit                                      |

### Keys in a project

| Key                | Action                                           |
| ------------------ | ------------------------------------------------ |
| `<left>`/`<right>` | Switch tab                                       |
| `1` / `2` / `3`    | changes / specs / config                         |
| `<tab>`            | Switch focus between the panels (top level only) |

### Keys in the change list

| Key   | Action                                                |
| ----- | ----------------------------------------------------- |
| `/`   | Filter the list (see below)                           |
| `f`   | Cycle open / archived / both                          |
| `a`   | Archive the selected change                           |
| `d`   | Discard the selected change                           |
| `e`   | Export the selected change as a zip                   |

### Keys in an open change

| Key                | Action                  |
| ------------------ | ----------------------- |
| `<left>`/`<right>` | Switch artifact sub-tab |

Sub-tabs stay inside the change: pressing `<right>` on the last one does not
spill over into the project tab bar.

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
| `archived` | Whether the change is open or archived           |
| `date`     | The archive date, for archived changes           |

The default is `name,tasks,specs`. On a narrow terminal, columns are dropped
from the right rather than squeezing the name past readability.

## Development

```bash
make lint
```
