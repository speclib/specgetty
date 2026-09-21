---
# specgetty-0ia0
title: paging keys are inert on the properties tab
status: completed
type: bug
priority: high
created_at: 2026-09-21T16:54:05Z
updated_at: 2026-09-21T16:57:48Z
---

Regression from `properties-tab` (8f7aab5).

The config tab was a single document and `docActive()` returned true for it
regardless of focus, so pgup/pgdown, ctrl+f/b, ctrl+d/u and gg/G scrolled it
immediately. Making it a split gave it a list that takes the keyboard first, and
those keys now do nothing until `tab` is pressed.

`document-viewer` already says what should happen: the paging keys apply when a
document `is displayed`, not when it is focused. The specs tab under-delivered
against the same requirement from the day it shipped; nobody noticed because its
list is the thing you interact with.

## Summary of Changes

Shipped as `page-from-either-half` (commit 7b6bdcb).

`docDisplayed()` sits beside `docActive()`, differing only by the focus check
on a split tab. The paging and jump keys ask the new one; `j` and `k` keep
asking the old one, because they have to choose between moving a list and
scrolling a document and focus is what decides that.

The reading position in the panel title was left alone: `specs-tab` reports it
only while the content holds the keyboard, which is deliberate and was not what
broke.

The specs tab had the same defect since it shipped, against the same
`document-viewer` requirement. Nobody noticed because its list is the half you
came to use, where the properties tab's list is three fixed rows.

Coverage floor raised: ui 83.9 to 84.2.
