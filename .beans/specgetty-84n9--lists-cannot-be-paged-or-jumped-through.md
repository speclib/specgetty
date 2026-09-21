---
# specgetty-84n9
title: lists cannot be paged or jumped through
status: in-progress
type: bug
priority: high
created_at: 2026-09-21T17:01:56Z
updated_at: 2026-09-21T17:01:56Z
---

The change list has 38 rows and moves one at a time: no pgdown, no gg, no G. The specs list has 24 and the same. Both side lists send those keys to the document beside them instead, which `page-from-either-half` (7b6bdcb) did on a misreading of the first report.

Decided: the keys act on whatever holds the keyboard. Lists become pageable, and paging a document again needs the keyboard on it. The picker gets the same keys; it has gg and G already but no paging.
