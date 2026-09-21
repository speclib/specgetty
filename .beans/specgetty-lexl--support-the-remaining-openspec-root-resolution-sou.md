---
# specgetty-lexl
title: support the remaining OpenSpec root resolution sources
status: scrapped
type: feature
priority: normal
created_at: 2026-09-18T16:34:14Z
updated_at: 2026-09-21T17:05:52Z
---

`follow-the-store` implements two of OpenSpec's four root sources: the nearest local root and a repo's `store:` declaration. Two are left.

- [ ] `--store <id>` on the command line, selecting a registered store directly
- [ ] the machine-wide `defaultStore` from ~/.config/openspec/config.json, which OpenSpec consults after the nearest root and any declaration have failed

Both are fallbacks that nothing on this machine uses yet, which is why they were left out. Precedence, from OpenSpec's core/root-selection.js: --store, then the nearest qualifying root or its declaration, then defaultStore.
