---
# specgetty-84n9
title: lists cannot be paged or jumped through
status: completed
type: bug
priority: high
created_at: 2026-09-21T17:01:56Z
updated_at: 2026-09-21T17:07:25Z
---

The change list has 38 rows and moves one at a time: no pgdown, no gg, no G. The specs list has 24 and the same. Both side lists send those keys to the document beside them instead, which `page-from-either-half` (7b6bdcb) did on a misreading of the first report.

Decided: the keys act on whatever holds the keyboard. Lists become pageable, and paging a document again needs the keyboard on it. The picker gets the same keys; it has gg and G already but no paging.

## Summary of Changes

Shipped as `page-the-lists` (commit 92ad97a).

The keys act on whatever holds the keyboard, by the same rule `j` and `k`
follow. One rule, five surfaces: the change list, the spec list, the properties
rows, any document, and the project picker.

- `listPage()` sizes a page to the rows the surface is showing, so the keys
  mean the same thing at any terminal height
- `setChangeCursor` remembers the selection, which a page move would otherwise
  lose where a single-row move does not
- the picker's own key block gained the page and half-page keys beside its
  existing `g` and `G`

## Supersedes

`page-from-either-half` (7b6bdcb), shipped an hour earlier, made the paging
keys bypass focus and always reach the document. That was a misreading of the
first report, which turned out to be about lists. Its three requirements are
removed with reasons: once a list is pageable there are two things a key could
move on a split tab, and a key that ignores focus has nothing to choose with.

Coverage floors raised: ui 84.2 to 84.5, total 86.5 to 86.8.
