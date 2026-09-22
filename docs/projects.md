# Projects

Working across every OpenSpec project on your machine: finding them, switching
between them, and reading what each one is.

- [The project picker](#the-project-picker)
- [The properties tab](#the-properties-tab)
- [Projects that keep their specs in a store](#projects-that-keep-their-specs-in-a-store)
- [Why `r` exists](#why-r-exists)

## The project picker

Press `p` from anywhere. It lists every OpenSpec project it has found, with its
spec, change and task counts, and `enter` switches to the highlighted one. You
land on that project's change list.

![The project picker](../demo/recordings/picker.gif)

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

![The properties tab](../demo/recordings/properties.gif)

The rows are grouped under headers the cursor skips over. A group keeps its
header when it holds nothing, so a project with no recorded schema still says
so.

| Row        | Shows                                                         |
|:-----------|:--------------------------------------------------------------|
| `config`   | the one configuration that applies, and the file it came from |
| `store`    | where the content comes from, or `local`                      |
| a schema   | one row per workflow schema the project's changes record      |

The list takes twenty columns wherever the panel can spare them, and no more:
its labels are short and fixed, so every column a wider terminal adds goes to
the content beside it and the divider stays where your eye left it.

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
