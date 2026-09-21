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

There are three levels, and the project picker opens over any of them. Which
second level `enter` reaches depends on the tab you are standing on.

```
   [ project picker ]   p opens it, esc closes it
           |            enter switches project
           v
   project view         the startup view; esc does nothing here
       | enter          tabs: changes | specs | properties
       |
       +--> change view       from the changes tab
       |      |               artifact sub-tabs: proposal | design | tasks | specs
       |      +--> change spec view    from that change's specs sub-tab
       |                               its deltas, marked with what they do
       +--> spec view         from the specs tab
                              an outline beside one requirement or scenario
```

Active and archived changes are shown together in one list, grouped, with the
active group first. There is no separate archive tab and no filter to choose
between them.

### Keys everywhere

| Key                        | Action                  |
| -------------------------- | ----------------------- |
| `j`/`k` or `<up>`/`<down>` | Move the cursor         |
| `<enter>`                  | Go one level deeper     |
| `<esc>`                    | Go one level back       |
| `E`                        | Open the file in your editor |
| `p`                        | Open the project picker |
| `s`                        | Rescan the open project |
| `q` / `ctrl-C`             | Quit                    |

### Keys in a project

| Key                | Action                                           |
| ------------------ | ------------------------------------------------ |
| `<left>`/`<right>` | Switch tab                                       |
| `1` / `2` / `3`    | changes / specs / config                         |
| `<tab>`            | Switch focus between the panels (top level only) |

### Opening a file in your editor

`E` hands the file a pane is showing to your own editor. It is bound wherever a
pane shows exactly one file, and nowhere else:

| Pane                                        | What `E` opens        |
| ------------------------------------------- | --------------------- |
| a change's proposal, design or tasks        | that artifact's file  |
| a spec on the specs tab                     | that spec's `spec.md` |
| the project row of the properties tab       | the configuration file |
| a change's specs sub-tab                    | nothing: several files |
| the schema and store rows                   | nothing: no file      |

The nav bar lists `E` exactly where it would do something, so the pane itself
tells you whether the key applies.

The editor comes from `$VISUAL` when that is set, and from `$EDITOR` otherwise.
That is the convention: `EDITOR` names a line editor that works without a full
terminal and `VISUAL` the full-screen one, so with both set the first one here
is what you meant for an interactive edit.

A value carrying arguments works, because the value is split on whitespace:

```sh
export EDITOR="code -w"
export VISUAL="emacsclient -nw"
```

What that does not handle is a quoted path inside the variable, such as
`EDITOR='/opt/my editor/bin/ed'`. No shell is involved, so the space is read as
a separator. `$VISUAL` gives you a second place to put a value that works.

With neither set, `E` says so on the nav bar and does nothing. There is no
fallback to `vi`: guessing is worse than saying, and it guesses wrong on a
machine that has no `vi`.

The interface stops drawing while the editor has the terminal and resumes when
it exits, then reads the project again, so your edit is on screen without
pressing `s`. The file opens at its top; no position within it is passed,
because the syntax for that differs per editor and would work for some people
and silently do nothing for the rest.

### Keys in the specs tab

The specs tab has two halves and `tab` moves the keyboard between them.

| Key     | Action                                             |
| ------- | -------------------------------------------------- |
| `tab`   | Move the keyboard between the list and the content |
| `j`/`k` | Change spec, or scroll it, depending on the focus  |

The selected spec is highlighted while the list has the keys and dimmed while
the content has them, and the title shows a reading position only in the second
case. When the content has the keys it takes the same scroll keys as an open
change.

`enter` opens the selected spec at its own level, described below. A spec whose
headings do not follow the OpenSpec structure does not open: the status line
says so and the cursor stays where it was.

### Keys in a spec

A spec opens as an outline of its requirements and scenarios beside a card
showing whichever one the cursor is on. The card is the spec rewritten for
reading: one requirement or one scenario at a time, clauses laid out with their
keyword above the text, rather than the raw markdown.

```
 +----------------------------+ +-------------------------------------+
 | Purpose                    | |                                     |
 | A change opens at its own  | |    Scenario: Ascend from a change   |
 |   navigation level         | |                                     |
 |   Descend from the list    | |    WHEN                             |
 | > Ascend from a change     | |       a change is open and the      |
 |   An open change names its | |       user presses esc              |
 |     schema                 | |                                     |
 | Artifact sub-navigation    | |    THEN                             |
 |   Right arrow at the last  | |       the change list SHALL be      |
 |     sub-tab                | |       displayed with the same       |
 |   Left arrow at the first  | |       change under the cursor       |
 +----------------------------+ +-------------------------------------+
```

| Key                        | Action                                        |
| -------------------------- | --------------------------------------------- |
| `<tab>`                    | Move the keyboard between the outline and the card |
| `j`/`k` or `<up>`/`<down>` | Move one node, or scroll the card             |
| `^f`/`^b`                  | Page the outline, or the card                 |
| `^d`/`^u`                  | Half a page                                   |
| `gg` / `G`                 | First and last                                |
| `<esc>`                    | Back to the specs tab, same spec selected     |

A label too long for the outline wraps rather than being cut, and the cursor
still moves one requirement or scenario per keystroke however many rows it
occupies. The cursor stays on the same node across a rescan, so editing the
file above it does not move it.

A scenario's content is shown whatever shape it is written in. Clauses written
as the OpenSpec template shows them get their keyword on a row of its own;
anything else is shown as the prose it is, in the place it was written. Nothing
is left out.

### Keys in a change's specs

A change carries a spec delta per capability it touches, and the specs sub-tab
lists them all as one document. `enter` opens them as an outline instead, marked
with what the change does to each requirement:

```
 +----------------------------+ +-------------------------------------+
 | change-search              | | [diff] old  new                     |
 | + The prompt names the     | |                                     |
 |     matchers it accepts    | |    Scenario: The legend             |
 |     The prompt is opened   | |                                     |
 |   ~ The first character    | |    WHEN                             |
 |   + Too narrow to say it   | |       the prompt is focused and     |
 | ~ A failed name search     | |       empty has been typed          |
 | - The old grammar          | |                                     |
 +----------------------------+ +-------------------------------------+
```

| Mark | On a requirement      | On a scenario                          |
| ---- | --------------------- | -------------------------------------- |
| `+`  | the change adds it    | the change introduces it                |
| `~`  | the change modifies it | the change edits it                    |
| `-`  | the change removes it | never: a delta cannot drop one          |
|      |                       | unchanged, and drawn back to say so     |

The marks on scenarios are the useful half. A modified requirement is restated
in full even to change one sentence, and more than half of what it restates is
usually text it does not touch. The outline says which half is which without
your reading it.

| Key                | Action                                              |
| ------------------ | --------------------------------------------------- |
| `<left>`/`<right>` | Move between the difference, the original and the new |
| `<esc>`            | Back to the change, specs sub-tab selected           |

Everything else is the same as a spec: `tab`, `j`/`k`, the paging keys and `E`
all behave as they do one level across.

The three-way choice appears only where there is something to compare, which is
a requirement the change modifies and the scenarios inside it. The difference is
what opens, and each node opens on its own difference rather than keeping the
last choice.

Only a change that has not been archived offers it. An archived change's deltas
have already been applied to the project's specs, so those specs are the result
rather than the original, and presenting them as the text the change modified
would be a lie. Such a change shows what it proposed and nothing beside it.

### When a spec opens as a report instead

`enter` always descends. If the file cannot be read as a spec, the view says why
instead of showing an outline:

```
 +--------------------------------------------------------------+
 |                                                              |
 |    specgetty cannot read airplane-mode as a spec             |
 |                                                              |
 |    These are the rules openspec itself reads a spec by, so   |
 |    validate, list and archive cannot see this file either.   |
 |                                                              |
 |    There is no ## Purpose section. Every spec needs one.     |
 |                                                              |
 |    line 3                                                    |
 |    ## ADDED Requirements is a delta header. It belongs in    |
 |    a change, under openspec/changes/<name>/specs/. A main    |
 |    spec keeps its requirements under ## Requirements...      |
 |                                                              |
 +--------------------------------------------------------------+
  esc back to specs   E edit   jk scroll
```

Every reason is listed, with the line it sits on. `E` opens the file in your
editor, and `esc` returns to the specs tab where the whole file is readable as
markdown. Fix the file and the next scan turns the report into an outline
without you leaving the view.

The rules are OpenSpec's own, transcribed from its parser rather than invented
here, so a file specgetty will not structure is one `openspec validate` also
rejects:

| rule                                                  | why it matters             |
| ----------------------------------------------------- | -------------------------- |
| requirements live under one `## Requirements` heading | openspec parses only that section |
| a `## Purpose` section                                | required of every spec     |
| `### Requirement: <name>` inside that section         | outside it, it is invisible |
| a `#### ` heading with content under it               | that is a scenario         |
| no `## ADDED Requirements` and the like               | delta headers belong in a change |

What OpenSpec does not define is how a scenario's content is written. Its own
model for a scenario is raw text, so specgetty accepts whatever is there rather
than requiring a convention no tool enforces.

### Keys in the change list

| Key   | Action                                                |
| ----- | ----------------------------------------------------- |
| `/`   | Filter the list (see below)                           |
| `f`   | Cycle active / archived / both                        |
| `a`   | Archive the selected change                           |
| `d`   | Discard the selected change                           |
| `e`   | Export the selected change as a zip                   |

### Exporting a change

`e` asks where to put the zip. The prompt opens on `export_dir` from your
config, or your home directory when that is unset, and `tab` completes a path
against the filesystem:

```
 ╭────────────────────────────────────────────╮
 │  Export "follow-the-store"                 │
 │                                            │
 │  → ~/Downloads/                            │
 │    follow-the-store-2026-09-21.zip         │
 │                                            │
 │  tab completes   ⏎ export   esc cancel     │
 ╰────────────────────────────────────────────╯
```

You type a directory; the filename is generated from the change name and the
date. Editing the directory redirects that one export and leaves your config
alone. A directory that does not exist is refused rather than created, and an
existing file is confirmed before it is replaced.

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
the whole thing fits, so there is nothing below.

### The same keys everywhere

`<pgdn>`/`<pgup>`, `ctrl-f`/`ctrl-b`, `ctrl-d`/`ctrl-u`, `gg` and `G` act on
whatever holds the keyboard, by the same rule `j` and `k` follow:

| The keyboard is on     | These keys move                       |
| ---------------------- | ------------------------------------- |
| the change list        | the selected change                   |
| a list beside a pane   | the selection in that list            |
| a document             | the document                          |
| the project picker     | the selected project                  |

A page is however many rows the thing is showing, so the keys mean the same at
any terminal size. On the specs and properties tabs, `<tab>` moves the keyboard
between the list and the content beside it.

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

Pressing `/` names the three matchers until you start typing, and a name search
that finds nothing suggests the same term as a `:` search.

## The properties tab

The third tab reports what the project is and where its parts come from, as a
list of rows beside their content:

```
 ╭─────────────╮ ╭────────────────────────────────────────╮
 │ project     │ │ schema: spec-driven                    │
 │ spec-driven │ │ default: yes                           │
 │ tinychange  │ │ changes: 26                            │
 │ store       │ │ source: package                        │
 ╰─────────────╯ │ artifacts:                             │
                 │   proposal -> proposal.md              │
                 ╰────────────────────────────────────────╯
```

| Row       | Shows                                                          |
| --------- | -------------------------------------------------------------- |
| `project` | the one configuration that applies, and the file it came from   |
| a schema  | one row per workflow schema the project's changes record        |
| `store`   | where the content comes from, or `local`                        |

`tab` moves between the list and the content, `j`/`k` move down the list.

A schema's definition is located by running `openspec schema which`, which costs
about a second, so nothing is read until you open the tab and nothing is read
twice for the same project. Editing a schema while `spg` is open will not show
up until you open the project again. Without the `openspec` CLI the rows still
report which schemas are in use and how many changes are on each, and say why
the rest is missing.

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
its header carries a `store` mark, and the properties tab reports where the
content came from.

Only the store's configuration applies. OpenSpec reads a pointing repo's file
for its `store:` key and takes everything else from the store, so a `context:`,
`rules:` or `operations:` block in the repo has no effect. The properties tab
names those keys where it finds them, because OpenSpec itself does not.

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
| `archived` | Whether the change is active or archived. Redundant against the group header, so not a default |
| `date`     | The archive date, for archived changes           |

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

## Development

```bash
make lint
```
