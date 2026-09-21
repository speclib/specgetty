## Context

See proposal.md for why. What matters here is where the notes come from today
and what the changelog's shape already gives us.

```
  release.sh            inserts "## [X.Y.Z] - date" under [Unreleased]
                        bumps src/VERSION, commits, tags, pushes
         │  tag push
         ▼
  release.yml
         ├─ goreleaser (linux)    release --clean
         └─ goreleaser (darwin)   release --clean
                  └─ changelog:   sort: asc, filters excluding ^docs: ^test:
                                  ^ci: and merges
```

The `changelog:` block in both configurations is goreleaser's commit-derived
notes. Nothing in either file names `CHANGELOG.md`.

## Goals / Non-Goals

Goals: the release page says what the changelog says; a release that has nothing
written stops rather than publishing a heading.

Non-Goals: changing how `CHANGELOG.md` is written or dated, which is correct.
A commit message convention. Generating a changelog from commits, which is the
thing being removed.

## Decisions

### The entry is extracted by heading, not by version number

The changelog's own shape is the interface. A section runs from `## [X.Y.Z]` to
the next line beginning `## [`, which is true of every entry in the file and is
what Keep a Changelog defines. Matching on the heading rather than searching for
the version string anywhere means a version mentioned inside prose, such as a
note that a key was removed in 0.8.0, cannot be mistaken for the start of its
section.

The extraction happens in `release.sh`, at the moment the heading is inserted
and the content under it is known to be what was just released. Doing it in the
workflow instead would mean parsing the file again from a checkout, at a point
where an empty entry can no longer stop anything: the tag is already pushed.

### The commit-derived changelog is turned off rather than filtered

Filtering `chore:` out of goreleaser's list would make the published notes less
noisy and still leave two descriptions of one release, written in two places, in
two registers. The first time a commit subject and a changelog entry disagreed
about the same change, the release page would carry both.

So `changelog: disable: true` in both configurations, and the notes come from
one place.

### An empty entry stops the release

An empty section is the observable form of "the changelog was not written before
the release was cut". The check belongs in `release.sh` before it tags, because
after the tag is pushed the workflow has already started and a release with
empty notes exists.

This sits beside the check that already refuses a dirty working copy, and for
the same reason: both are things that would make the release wrong in a way that
is annoying to undo.

### The release is created once, and both runs carry the notes

The two goreleaser runs are already sequenced: `release-darwin` declares
`needs: release-linux`, and has since before v0.7.0. So the linux run creates
the release and the darwin run adds its artifacts to one that exists. This was
checked rather than assumed, after an earlier draft of this proposal asserted a
race that is not there.

What changes is that the notes become something a run could disagree about.
They are passed to both rather than only to the run that creates the release,
so that whatever either does with an existing release, neither can publish a
different account of it.

The sequencing is load-bearing for this change without being introduced by it,
which is why it becomes a requirement here: an undocumented guarantee that
something now depends on is worth writing down before it is changed by someone
who does not know it is relied upon.

## Risks / Trade-offs

- The extraction is text handling in bash, and this project ships darwin builds,
  where BSD tools differ. → The changelog editing already in `release.sh` uses
  awk for exactly this reason, and says so in a comment. The extraction follows
  it rather than reaching for `sed -i`.
- A changelog entry containing a line that begins `## [` would truncate its own
  section. → That is a heading in Keep a Changelog, so a line beginning that way
  inside an entry is already malformed. Worth nothing more than knowing it.
- The release notes reach the workflow through `CHANGELOG.md` in the checkout
  rather than through a committed file, so the tag must carry the changelog
  entry it publishes. → It does by construction: the release commit that
  inserts the heading is the commit the tag points at.

## Migration Plan

None. The first release after this publishes its own entry; the ones already
published keep the notes they have.
