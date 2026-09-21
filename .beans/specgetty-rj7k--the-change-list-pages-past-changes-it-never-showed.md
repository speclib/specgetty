---
# specgetty-rj7k
title: The change list pages past changes it never showed
status: todo
type: bug
created_at: 2026-09-21T17:36:36Z
updated_at: 2026-09-21T17:36:36Z
---

The change list's page key moves the cursor further than the pane shows.

`listPage()` returns the number of *lines* the body has, and `moveListCursor`
spends that on *changes*. The group headers occupy lines that hold no change,
so a page steps past changes that were never on screen.

```
 80x24              lines in the body   changes shown   a page moves
 today                     13                 11          13 changes
 with the group gap        13                 10          13 changes
```

Two changes skipped today, three once `space-the-change-groups` lands.

`document-viewer` already states the rule this breaks:

> Moving by a page SHALL move by the number of rows the surface is currently
> showing

and the change `open-a-spec-in-detail` amends that requirement to say what it
means over items whose rows are not all one row tall: a page moves the cursor by
as many items as fill one pane of rows. The change list needs the same walk,
counting a group header and a group gap as the rows they occupy.

## Notes

- Blocked by nothing, but the fix is cleaner once `open-a-spec-in-detail` has
  written the page walk, since both surfaces want the same one.
- `change-list-view` has a requirement `The change list can be paged and jumped
  through` which is where the scenario belongs.
