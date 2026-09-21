## REMOVED Requirements

### Requirement: Rescan is scoped to the open project
**Reason**: The key existed to recover from four things the filesystem watcher
could not see: the store registry, the store's git state, a watcher that failed
to start, and a signal dropped while a scan was in flight. This change closes
the first two and the fourth, and `specgetty-7lc7` reports the third, so the key
has no job left. A manual refresh in a tool that watches the filesystem is a
statement that the watching is not trusted.

**Migration**: Pressing `p` and then `enter` on the open project re-resolves and
re-reads it, which is what `s` did. That path is unchanged and is what the
README now points at.
