---
# specgetty-57gj
title: 'release script: jj support and portability fixes'
status: todo
type: task
priority: normal
created_at: 2026-09-15T19:24:09Z
updated_at: 2026-09-15T19:24:09Z
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
