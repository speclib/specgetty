# Keys

Every key specgetty binds, by where the keyboard is. The rule underneath all of
it: a key acts on whatever holds the keyboard, and a page is however many rows
that thing is showing, so the keys mean the same at any terminal size.

- [Everywhere](#keys-everywhere)
- [In a project](#keys-in-a-project)
- [In the specs tab](#keys-in-the-specs-tab)
- [In the change list](#keys-in-the-change-list)
- [In an open change](#keys-in-an-open-change)
- [The same keys everywhere](#the-same-keys-everywhere)

The two levels below a project have their keys where their behaviour is
described: [a spec](reading-specs.md#keys-in-a-spec) and
[a change's specs](reading-specs.md#keys-in-a-changes-specs). `E` has a page of
its own, [the editor](editor.md).

## Keys everywhere

| Key                        | Action                       |
|:---------------------------|:-----------------------------|
| `j`/`k` or `<up>`/`<down>` | Move the cursor              |
| `<enter>`                  | Go one level deeper          |
| `<esc>`                    | Go one level back            |
| `E`                        | Open the file in your editor |
| `p`                        | Open the project picker      |
| `q` / `ctrl-C`             | Quit                         |

There is no rescan key. The open project is watched and re-read whenever
anything under its `openspec/` changes, and the two things that live outside
that tree are covered too: the store registry is watched for a project that
declares a store, and a store's git state is re-read whenever you enter the
properties tab. To force a read anyway, press `p` and `<enter>` on the project
you are already in, which re-resolves and re-reads it.

## Keys in a project

| Key                | Action                                           |
| ------------------ | ------------------------------------------------ |
| `<left>`/`<right>` | Switch tab                                       |
| `1` / `2` / `3`    | changes / specs / config                         |
| `<tab>`            | Switch focus between the panels (top level only) |

## Keys in the specs tab

The specs tab has two halves and `tab` moves the keyboard between them.

| Key     | Action                                             |
| ------- | -------------------------------------------------- |
| `tab`   | Move the keyboard between the list and the content |
| `j`/`k` | Change spec, or scroll it, depending on the focus  |

The selected spec is highlighted while the list has the keys and dimmed while
the content has them, and the title shows a reading position only in the second
case. When the content has the keys it takes the same scroll keys as an open
change.

`enter` opens the selected spec at [its own level](reading-specs.md#keys-in-a-spec).
A spec whose headings do not follow the OpenSpec grammar opens too, as
[a report saying why](reading-specs.md#when-a-spec-opens-as-a-report-instead)
rather than as an outline.

![The specs tab, and a spec opened in detail](../demo/recordings/specs.gif)

## Keys in the change list

| Key   | Action                                                |
| ----- | ----------------------------------------------------- |
| `/`   | Filter the list                                       |
| `a`   | Archive the selected change                           |
| `d`   | Discard the selected change                           |
| `e`   | Export the selected change as a zip                   |

## Keys in an open change

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

## The same keys everywhere

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
