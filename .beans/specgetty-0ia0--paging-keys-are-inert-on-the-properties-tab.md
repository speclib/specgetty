---
# specgetty-0ia0
title: paging keys are inert on the properties tab
status: in-progress
type: bug
priority: high
created_at: 2026-09-21T16:54:05Z
updated_at: 2026-09-21T16:54:05Z
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
