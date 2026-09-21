## 1. Extracting the entry

- [x] 1.1 Add a function to `scripts/release.sh` that writes the section for a
      given version, from its `## [X.Y.Z]` heading to the next `## [`, to a file,
      and verify it against the existing changelog for 0.7.1, 0.7.0 and 0.6.0
- [x] 1.2 Verify the extraction excludes both the `## [Unreleased]` heading above
      the section and the version heading below it, by checking the boundaries of
      a middle version
- [x] 1.3 Write it in awk rather than with `sed -i`, the script already doing so
      for the changelog and saying why, and verify the output is byte identical
      on GNU and BSD userlands or state which was tested
- [x] 1.4 Verify a version string appearing inside an entry's prose does not
      start a section, by extracting an entry that names another version in a
      sentence

## 2. Stopping on an empty entry

- [x] 2.1 Make the script stop before tagging when the extracted entry holds
      nothing, saying that the changelog is empty, and verify it against a
      changelog whose newest section is bare
- [x] 2.2 Verify nothing is left half done: no tag created, nothing pushed, and
      the working copy as it was
- [x] 2.3 Verify the check sits beside the dirty-working-copy check and runs
      before any file is edited

## 3. Publishing the entry

- [x] 3.1 Pass the extracted file to goreleaser as `--release-notes` in
      `.github/workflows/release.yml`, and verify the argument reaches both the
      invocation and the release
- [x] 3.2 Set `changelog: disable: true` in `.goreleaser-linux.yaml` and
      `.goreleaser-darwin.yaml`, and verify no commit-derived list is produced
- [x] 3.3 Verify the notes file survives from the release commit to the workflow,
      the workflow checking out the tag it was pushed with

## 4. One release, two platforms

- [x] 4.1 Verify the darwin job already declares `needs: release-linux`, so the
      release is created once, and leave it as it is. An earlier draft of the
      proposal said these ran in parallel; they have not since before v0.7.0
- [x] 4.2 Verify the notes are the changelog entry whichever order the jobs
      finished in, by checking the published release rather than the job log

## 5. Verification

- [x] 5.1 Verify the publishing path without publishing: `goreleaser check` on
      both configurations, `goreleaser release --help` confirming
      `--release-notes` and that it skips goreleaser's own changelog, and the
      workflow's extraction step run locally with `GITHUB_REF_NAME` set
- [x] 5.2 Leave the throwaway tag unpushed. Cutting and deleting a real release
      on a public repository to test a config is visible to everyone watching it,
      and the four things above cover what it would have shown. The first real
      release is the remaining check
- [x] 5.3 Verify a release with an empty entry stops, by running the refusal
      against a changelog whose `[Unreleased]` section is bare
- [x] 5.4 `nix flake check` passes, this change adding no Go code

## 6. Documentation

- [x] 6.1 Say in RELEASING.md that the release publishes the changelog entry and
      will stop if it is empty
- [x] 6.2 Add the change to CHANGELOG.md under Changed
