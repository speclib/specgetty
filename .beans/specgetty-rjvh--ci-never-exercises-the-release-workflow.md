---
# specgetty-rjvh
title: CI never exercises the release workflow
status: todo
type: task
priority: high
created_at: 2026-09-23T09:34:44Z
updated_at: 2026-09-23T09:34:44Z
---

Both v0.7.3 release bugs reached a real release because nothing runs the tag-triggered workflow until a tag is pushed. `nix flake check` cannot reach it.

`goreleaser release --snapshot --clean` builds and packages without a tag and without publishing. Run in the Check workflow it would have caught the dirty-tree failure and any future error in `.goreleaser-*.yaml` in seconds.

Note: a local `goreleaser release --clean` deleted the untracked-but-tracked contents of `design/` during investigation on 2026-09-23. Cause not established. Worth understanding before wiring goreleaser into a workflow that runs on every push.

Related: specgetty-2hf2.
