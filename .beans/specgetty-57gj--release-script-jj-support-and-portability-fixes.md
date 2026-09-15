---
# specgetty-57gj
title: 'release script: jj support and portability fixes'
status: completed
type: task
priority: normal
created_at: 2026-09-15T19:24:09Z
updated_at: 2026-09-15T19:35:56Z
---

OpenSpec change `release-with-jj` (tinychange schema) holds the spec delta and
tasks. Explored 2026-09-15.

## What is wrong

- `scripts/release.sh` assumes git owns the working copy. In a jj workspace the
  clean check, the commit and the push are all wrong.
- The clean check is `git diff --quiet`, which ignores untracked files. A new
  file that belongs in the release is silently left out.
- Line 159 uses a `\n` inside a `sed` replacement to insert the version heading
  into CHANGELOG.md. That is GNU-only; BSD sed writes a literal `n`. This
  project ships darwin builds, so a mac maintainer hits it.

## What jj can and cannot do (jj 0.41)

| Operation | jj                                            |
| --------- | --------------------------------------------- |
| commit    | `jj commit -m ...`                            |
| bookmark  | `jj bookmark set <name> -r @-`                |
| tag       | `jj tag set vX.Y.Z -r @-`                     |
| push      | `jj git push --bookmark <name>` (bookmarks only) |
| push tags | NOT supported; `jj git push` never carries tags |

In a colocated workspace jj exports refs to git automatically, so a tag made
with `jj tag set` lands in `.git/refs/tags` and plain `git push origin vX.Y.Z`
carries it. That is the route the tasks take.

The tag matters more than usual here: the GitHub Actions release workflow is
triggered by a `v*` tag, so a release that commits but does not tag publishes
nothing.

## Already fixed, not part of this change

Commit 7230017 fixed the vendorHash update touching only package.nix while its
comment claimed flake.nix, and pointed RELEASING.md at main instead of a master
branch that does not exist.

## Summary of Changes

Applied as openspec change `release-with-jj` (tinychange), archived to
`openspec/changes/archive/2026-09-15-release-with-jj/`. The delta created
`openspec/specs/release-process/` with 3 requirements.

### What changed in scripts/release.sh

- Detects jj when `jj root` succeeds and uses jj throughout; git otherwise, on
  the path it already had.
- Clean check covers untracked files now. The old `git diff --quiet` passed on
  an untracked file, so a new file belonging in the release was silently left
  out. Demonstrated: with a stray file present the old check reports clean.
- Under jj: `jj commit`, `jj bookmark set <branch> -r @-`, `jj tag set`, then
  `jj git push --bookmark`, then `git push origin <tag>`.
- CHANGELOG heading inserted with awk instead of a `\n` inside a sed
  replacement, which is a GNU extension that writes a literal `n` on BSD.
- All `sed -i` replaced with a helper that works on both userlands.

### Order matters, and testing showed why

`jj git push --bookmark main` fails outright when the bookmark is not tracking
its remote counterpart. In the lab that failure happened *after* the tag push,
leaving the remote with a tag pointing at a commit no branch reached. So the
script now tracks the bookmark first, and pushes the bookmark before the tag.

### A regression I introduced and caught

The first version of the portable editor moved a `mktemp` file over each
target, which carried mktemp's 0600 onto files that were 0644. Git tracks only
the executable bit, so this would never have appeared in a diff. Fixed by
writing back through the original file. File modes are now part of the
verification matrix.

### Verified

A harness in a temp directory builds throwaway repos with a bare remote, stubs
`gum` and `nix`, and runs the real script end to end:

| Case                | Result                                        |
| ------------------- | --------------------------------------------- |
| jj colocated        | branch and tag on the remote at the same commit |
| plain git           | unchanged from today, annotated tag            |
| dirty, modified     | stops, writes nothing (both VCS)               |
| dirty, untracked    | stops, writes nothing (both VCS)               |
| file modes          | 644 before and after                          |
| changelog heading   | on its own line                               |

`nix flake check` passes. `openspec validate --all --strict` is 18/18.

### Known difference between the two paths

The git path makes an annotated tag (`git tag -a`); `jj tag set` has no message
option, so the jj path makes a lightweight one. goreleaser reads the tag name,
so both publish, but the tags differ in kind.

### Not verified

Nothing was released from this repository: `src/VERSION` is still 0.2.0. The BSD
sed behaviour is reasoned from the known GNU/BSD difference, not observed, since
only GNU sed is available here.
