## Why

The release notes on GitHub are a list of commit hashes. This is the whole of
what v0.7.1 published:

```
  ## Changelog
  * 338f58ef… chore: release v0.7.1
  * 42fe7a04… close the editor bean
  * f24aed3d… feat: E opens the file a pane is showing in your editor
```

and v0.7.0 published eleven such lines, five of them `chore: close the … bean`.

Meanwhile `CHANGELOG.md` carried, for that same release, the paragraph that
actually tells someone what changed:

> `E` opens the file a pane is showing in your own editor, from a change's
> artifacts, a spec, or the project configuration. The editor is `$VISUAL`, then
> `$EDITOR` [...] With neither set the key says so and does nothing rather than
> guessing.

Nobody installing `spg` reads `CHANGELOG.md`. They read the release page.

### The changelog itself is fine

Worth stating, because the symptom points the wrong way. Each release commit
changes `CHANGELOG.md` by exactly two lines, which looks like a script that
gave up:

```
  338f58e chore: release v0.7.1   CHANGELOG.md | 2 ++
  e36899a chore: release v0.7.0   CHANGELOG.md | 2 ++
  2150ec8 chore: release v0.6.0   CHANGELOG.md | 2 ++
```

Those two lines are the version heading being inserted above content that was
already written under `## [Unreleased]`, which is exactly what
`insert_version_heading` is for and exactly what Keep a Changelog asks. The
entries are right, dated and in place. Nothing downstream reads them.

Neither `.goreleaser-linux.yaml` nor `.goreleaser-darwin.yaml` mentions
`CHANGELOG.md`. Both carry a `changelog:` block that builds notes from git
commit subjects, and its filters exclude `^docs:`, `^test:`, `^ci:` and merges
but not `chore:` or a bare subject, which is why half of v0.7.0's notes are bean
bookkeeping.

## What Changes

- The release publishes the changelog entry it just wrote. `release.sh` extracts
  the section for the version being released, from its `## [X.Y.Z]` heading to
  the next `## [`, and the workflow hands that file to goreleaser as the release
  notes.

- The commit-derived changelog is turned off in both goreleaser configurations.
  A release publishes one account of itself, and two would disagree the first
  time a commit subject and a changelog entry described the same change
  differently.

- A release refuses to publish an empty entry. A version whose section holds
  nothing means the changelog was not written before the release was cut, which
  is worth stopping for rather than publishing a heading with nothing under it.

- The notes are passed to both goreleaser jobs rather than to one. The darwin
  job already waits on the linux job, so the release is created once and the
  second run adds its artifacts to it; passing the notes to both means neither
  run can publish an account the other would contradict, whatever either does
  with a release that already exists.

Not in scope:

- **Tidying the commit subjects.** Once the notes come from the changelog, what
  a `chore:` commit is called stops being published, so a naming convention for
  them is a separate question and not an urgent one.

## Capabilities

### Modified Capabilities
- `release-process`: a release publishes the changelog entry it wrote, and
  refuses to publish an empty one

## Impact

- `scripts/release.sh`: extracts the entry, and stops on an empty one
- `.github/workflows/release.yml`: passes the notes file, and sequences the two
  jobs
- `.goreleaser-linux.yaml`, `.goreleaser-darwin.yaml`: the commit-derived
  changelog is disabled

## Rollback

Its own commit. Reverting returns the release notes to a list of commit hashes,
which is what v0.7.0 and v0.7.1 published.
