## Why

The padding is already there. It stops one line early.

`renderDetailPanel` wraps the project header in `Padding(1, 1)` and then writes
everything below it into the same builder with no style at all. So the header
sits one column in from the border and the eight things under it sit on it:

```
| /home/pim/gh.speclib/specgetty                  |
|changes  specs  config                           |
|name                               tasks   specs |
|fix-openspec-detection-false-positives  6/6    1 |
```

On the right it is worse than untidy. Every wrapped line of every markdown
document ends flush against the border, so a word and the frame share a column
and the border stops reading as a frame:

```
|This is too loose, any nested directory with that name triggers a match, inside|
|test fixtures, etc.                                                            |
```

The project picker already does this correctly. It is not a view to fix, it is
the reference to copy.

## What Changes

- The main panel insets its content by one column on each side, which covers
  every view drawn inside it: both tab bars, the change table, the search
  prompt, the specs split, the config pane, the markdown documents, the change
  header, the placeholder, and the log panel
- The project header gives up the horizontal half of its own padding, keeping
  the blank line above and below. Otherwise it sits two columns in while
  everything under it sits one
- The hardcoded two-space indent on "No project selected. Press p to pick one."
  goes away. It was compensating for the missing inset
- The nav bar indents its text by one column and keeps its background spanning
  the full width

Not in scope: vertical padding. The top and bottom already have room, and the
bean was narrowed to left and right for that reason.

Not in scope: the gap between table columns. That is real and it is being fixed,
but as `widen-table-column-gaps`, so that either can be reverted without the
other.

## Capabilities

### New Capabilities
- `panel-layout`: where the content of the main panel begins and ends relative
  to its border, stated once rather than repeated across the eight capabilities
  whose views happen to be drawn there

### Modified Capabilities
- None. No requirement in any existing capability says where its view starts, so
  there is nothing to restate. This change adds a rule that was never written
  down, which is why eleven views quietly disagreed about it

## Impact

- `src/ui/ui.go`: `renderPanel` gains the inset; `renderDetailPanel` loses the
  horizontal half of the header's padding and the compensating indent;
  `renderNavBar` indents its text
- The content width shrinks by two columns everywhere. Markdown wraps two
  columns earlier, the change table's flexible column gives up two, and the
  specs split divides a smaller number
- Depends on `unify-content-width` landing first. That change collapses the
  three copies of `m.width - 2` into one method; without it this change has to
  edit the same arithmetic in three files and keep them in step by hand

## Rollback

Shipped as its own commit on top of v0.5.0, so `git revert` takes the look back
out and leaves `unify-content-width` in place. The revert also moves this change
out of `openspec/changes/archive/` and back into `openspec/changes/`, where
`openspec list` will show it as active again. That is correct, and it is worth
knowing before it happens.
