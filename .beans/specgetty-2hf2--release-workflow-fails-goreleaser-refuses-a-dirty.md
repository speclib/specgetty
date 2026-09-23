---
# specgetty-2hf2
title: 'Release workflow fails: goreleaser refuses a dirty tree'
status: in-progress
type: bug
priority: critical
created_at: 2026-09-23T09:21:29Z
updated_at: 2026-09-23T09:21:29Z
---

`v0.7.3` failed at the GoReleaser step:

    git is in a dirty state
    ?? release-notes.md

The `Extract the release notes` step writes `release-notes.md` into the
checkout, and goreleaser validates git state before releasing. An untracked
file counts as dirty.

Not a regression: goreleaser was v2.18.2 in both the green `v0.7.2` run and
the red `v0.7.3` one. The step arrived with `publish-the-changelog-entry`
after v0.7.2, so v0.7.3 is the first tag to run it. It has never worked.

Nothing was half-published: no artifacts, no release object. `release-darwin`
never ran, being `needs: release-linux`.

- [ ] write the notes to $RUNNER_TEMP rather than into the checkout, both jobs
- [ ] re-point the v0.7.3 tag at the fix and push it
- [ ] watch the run publish linux and darwin
