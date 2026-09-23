---
# specgetty-2hf2
title: 'Release workflow fails: goreleaser refuses a dirty tree'
status: completed
type: bug
priority: critical
created_at: 2026-09-23T09:21:29Z
updated_at: 2026-09-23T09:34:44Z
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

- [x] write the notes to $RUNNER_TEMP rather than into the checkout, both jobs
- [x] re-point the v0.7.3 tag at the fix and push it
- [x] watch the run publish linux and darwin

## Summary of Changes

Two bugs, both introduced by `publish-the-changelog-entry` and both first
exercised by v0.7.3, which was the first tag cut after it landed.

**1. `git is in a dirty state`.** The `Extract the release notes` step wrote
`release-notes.md` into the checkout, and goreleaser validates git state before
releasing and counts an untracked file as dirty. Fixed in `bafc146`: the notes
go to `$RUNNER_TEMP` in both jobs.

Not a regression from outside. goreleaser was v2.18.2 in both the green v0.7.2
run and the red v0.7.3 one; the step simply did not exist at v0.7.2.

**2. An empty release body.** With the first bug fixed the run went green and
uploaded all six artifacts, but the body was one newline. `changelog: disable:
true` turns off the pipe that reads `--release-notes`, so the notes file was
never read. Verified offline rather than guessed: with the key present
goreleaser prints no `generating changelog` line and writes no body; without it
the pipe runs. `--release-notes` alone already replaces the commit-derived list,
which is all the disable was ever for. Fixed in `b9ed269`, both configs.

`.gitignore` gained `dist/` in the same commit, a local goreleaser run leaving
it behind being another way to dirty the tree bug 1 tripped over.

## What shipped

v0.7.3 is published with all six artifacts, linux and darwin, and its notes were
set by hand from the changelog. The config fix is verified offline; it is proven
end to end by the next release.

## Not done

Nothing exercises a tag-triggered workflow, which is why both bugs reached a
real release. A `goreleaser release --snapshot` in the Check workflow would
catch this class in seconds. Left for a follow-up bean rather than folded in
here.
