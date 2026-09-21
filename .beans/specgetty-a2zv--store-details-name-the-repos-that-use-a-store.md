---
# specgetty-a2zv
title: store details name the repos that use a store
status: scrapped
type: feature
priority: normal
created_at: 2026-09-18T16:34:14Z
updated_at: 2026-09-21T17:05:43Z
---

The store details sub-tab on the config tab reports id, root, origin, remote and git state. It does not say which other repos read from the same store.

The reverse index is nearly free: the walk already descends into every pointing repo's openspec/ directory and returns nil there, so collecting the `store:` pointers on the way past costs one file read per skipped candidate.

Deferred out of the `follow-the-store` change (specgetty-8quc) so that change stayed about following one pointer rather than building an index.

- [ ] Collect store pointers during the walk
- [ ] Carry the reverse index to the open project
- [ ] Show `used by` on the store details sub-tab
