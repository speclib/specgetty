---
# specgetty-58sr
title: releases publish commit hashes instead of the changelog
status: in-progress
type: bug
priority: normal
created_at: 2026-09-21T22:33:36Z
updated_at: 2026-09-21T22:40:47Z
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
