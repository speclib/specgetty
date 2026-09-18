## Why

`design/specgetty-layout-improvement.png`. The tab content is drawn into the
panel with nothing around it, so the tab bar floats above a region with no edge:

```
| [changes] [specs] [config]          |     | [changes] [specs] [config]          |
| name              tasks specs       |     | ╭─────────────────────────────────╮ |
| fix-openspec...   6/6   1           |     | │ name            tasks     specs │ |
| replace-git...    24/24 1           |     | │ fix-openspec... 6/6       1     │ |
                                            | ╰─────────────────────────────────╯ |
  today                                       the mockup
```

With an edge, the chips read as tabs belonging to the thing below them, and the
content has a boundary of its own rather than borrowing the panel's.

## What Changes

- The content of each tab is drawn inside its own border, directly under the tab
  bar with no blank row between them
- The config tab's source filename stays above that border, with the tab bar.
  It names what is in the box, which makes it the same kind of label the chips
  are, not part of the content
- The changes tab's search prompt goes inside the border. It reports how many
  rows the box is showing, so it belongs to the box
- An open change gets the same treatment: its artifact document is boxed, its
  change header and sub-tab row stay above

Not in scope: the specs tab splitting into one box per pane, and the border
carrying the focus signal. That is `light-the-focused-pane`, which is separate
so that either can be taken back without the other.

Not in scope: the outer panel border, which keeps the colour and the meaning it
has today.

## Capabilities

### Modified Capabilities
- `panel-layout`: the panel says where its content begins relative to its border.
  It now also says that the content sits in a border of its own, and which parts
  of a tab are content rather than chrome

## Impact

- `src/ui/ui.go`: `renderDetailPanel` wraps the tab content; the height budget
  loses two rows and the width four columns
- The content width goes from `m.width - 4` to `m.width - 8`. At the 60-column
  minimum the change table has 52 columns, where it drops a column at 28, so no
  column is lost that was not already lost
- The visible row count falls by two everywhere. At 20 rows that is eleven
  changes down to nine, which is the real cost of this change

## Rollback

Its own commit. `git revert` takes the borders out and leaves `unify-focus-state`
in place. The revert also moves this change back into `openspec/changes/`, where
`openspec list` will show it as active again.
