## Why

Bean `specgetty-84n9`. Every list in the application is walked one row at a time.

```
  list                rows   pgdown                   G
  ──────────────────  ────   ──────────────────────   ──────────────────────
  changes table         38   nothing                  nothing
  specs list            24   scrolls the DOCUMENT     scrolls the DOCUMENT
  properties list        4   scrolls the DOCUMENT     scrolls the DOCUMENT
  project picker         n   nothing                  moves to the last row
```

The change list of this project is thirty-eight rows and the only way down it is
`j`. The two side lists are worse than inert: they send those keys to the
document beside them, which is what `page-from-either-half` did an hour ago on a
misreading of the first report of this.

## What Changes

- The paging and jump keys act on whatever holds the keyboard. One rule, no
  exceptions:

```
  the keyboard is on     pgdn / pgup / ctrl+f / ctrl+b / ctrl+d / ctrl+u / gg / G
  ──────────────────     ──────────────────────────────────────────────────────
  the changes table      move the row cursor
  a side list            move the list cursor
  a document             scroll it
  the log panel          scroll it
  the project picker     move its cursor
```

- A page is what the list can show, so `pgdown` lands about where the eye would.
  Half-page keys move half of it. `gg` and `G` go to the first and last row
- The project picker gains the paging keys. It has had `gg` and `G` since it
  shipped and nothing between them
- `page-from-either-half` is undone. Paging a document again requires the
  keyboard to be on it, which is how the specs tab has always worked and is what
  makes one rule possible

Not in scope: `j` and `k`, which already follow the keyboard and are what this
change makes the rest of the keys agree with.

## Capabilities

### Modified Capabilities
- `document-viewer`: the requirement that the paging keys ignore focus is
  removed. It was added an hour ago and is the thing being corrected
- `change-list-view`: the change list can be paged and jumped through
- `specs-tab`: so can the spec list, and its own paging requirement is removed
- `config-tab-display`: so can the properties rows, and the same removal
- `project-picker`: the picker gains the paging keys beside its jumps

## Impact

- `src/ui/ui.go`: six key handlers gain a branch for the list that holds the
  keyboard, and the picker's own key block gains the paging keys
- `src/ui/docview.go`: `docDisplayed` goes, having existed for one commit

## Rollback

Its own commit. Reverting restores both the missing list paging and the wrong
document paging, since this change replaces one with the other.
