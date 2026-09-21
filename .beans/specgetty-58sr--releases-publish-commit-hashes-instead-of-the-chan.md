---
# specgetty-58sr
title: releases publish commit hashes instead of the changelog
status: completed
type: bug
priority: normal
created_at: 2026-09-21T22:33:36Z
updated_at: 2026-09-21T22:46:29Z
---

The GitHub release notes are a list of commit hashes, not the changelog.

v0.7.1 published three lines, one of them `chore: release v0.7.1`. v0.7.0
published eleven, five of them `chore: close the ... bean`.

`CHANGELOG.md` itself is correct. Each release commit changes it by exactly two
lines, which is the version heading inserted above content already written
under `## [Unreleased]`, which is what `insert_version_heading` is for. Nothing
downstream reads it: neither goreleaser config mentions `CHANGELOG.md`, and both
carry a `changelog:` block that builds notes from commit subjects. Its filters
exclude `^docs:`, `^test:`, `^ci:` and merges, but not `chore:` or a bare
subject.

## OpenSpec change

`publish-the-changelog-entry`, validated strict, 21 tasks.

- `release.sh` extracts the section for the version being released and the
  workflow hands it to goreleaser as the notes.
- The commit-derived changelog is disabled in both configs rather than filtered:
  two accounts of one release would disagree eventually.
- A release with an empty entry stops before tagging, beside the existing
  dirty-working-copy check.
- The two goreleaser jobs both run `release --clean` in parallel against the
  same tag. It survives today because neither writes anything the other would
  contradict; both writing notes is where that stops being true. The darwin job
  is sequenced after linux.

## Summary of Changes

Shipped as `publish-the-changelog-entry` (commit `8c78c21`).

- `scripts/changelog-entry.sh` reads one section of `CHANGELOG.md`, from its
  `## [<label>]` heading to the next `## [`. One reader, two callers: the
  release script checks something was written, the workflow publishes it.
- `release.sh` refuses to release an empty `[Unreleased]` section, before any
  file is edited. After the tag is pushed the workflow has started and a
  release with empty notes already exists.
- Both goreleaser runs are handed the entry with `--release-notes`, and
  `changelog: disable: true` replaces the commit-derived list in both configs.

## A claim in my own proposal was wrong

The proposal said the two goreleaser jobs race, both running `release --clean`
in parallel. They do not: `release-darwin` has declared `needs: release-linux`
since before v0.7.0. I found it while implementing, and corrected the proposal,
the design, the spec delta and task 4.1 before archiving rather than shipping a
record that says something untrue.

The requirement about one release being created stayed, reworded: it is an
existing guarantee that this change makes load-bearing, and that is worth
writing down before someone who does not know it is relied upon changes it.

## What I did not do

Task 5.1 called for cutting a throwaway release on a real tag and deleting it.
I did not: that is visible to everyone watching a public repository. Instead:
`goreleaser check` passes on both configs, `goreleaser release --help` confirms
`--release-notes` exists and skips its own changelog generation, and the
workflow's extraction step was run locally with `GITHUB_REF_NAME=v0.7.2`, which
produced the 76-line entry. Tasks 5.1 and 5.2 were rewritten to say so. The
first real release is the remaining check.

## Notes

- One test passed for the wrong reason and was tightened: the ordering assertion
  used the substring `"tag"`, which matched a comment earlier in the file. It
  now names `git tag -a` and `jj tag set`.
- Two pre-existing em dashes in RELEASING.md were fixed while editing it.
- Coverage: src 87.7%, total 90.5%.
