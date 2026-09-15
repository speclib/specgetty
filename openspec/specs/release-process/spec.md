# release-process Specification

## Purpose
Cutting a release: bumping the version, dating the changelog, keeping the two
vendor hashes in step, and committing, tagging and pushing, whether the working
copy is managed by git or by jj.

## Requirements

### Requirement: The release works in a jj workspace
The release script SHALL detect a jj workspace and use jj to commit, tag, move
the bookmark and push, rather than assuming git owns the working copy.

#### Scenario: Releasing from a jj colocated workspace
- **WHEN** the script runs in a directory where `jj root` succeeds
- **THEN** it SHALL create the release commit, tag and bookmark through jj, and
  push them to the remote

#### Scenario: Releasing from a plain git repository
- **WHEN** the script runs in a directory with no jj workspace
- **THEN** it SHALL behave exactly as it does today, using git throughout

#### Scenario: The tag reaches the remote
- **WHEN** a release is made from a jj workspace
- **THEN** the `vX.Y.Z` tag SHALL exist on the remote, because the release is
  published by a tag and a release that does not tag does not publish

### Requirement: A dirty working copy stops the release
The script SHALL refuse to release when the working copy has changes that would
not be part of the release commit, whichever tool manages it.

#### Scenario: Uncommitted changes under jj
- **WHEN** the jj working copy commit contains changes
- **THEN** the script SHALL stop and say so, without editing any file

#### Scenario: Untracked files
- **WHEN** the working copy contains a file that is not tracked
- **THEN** the script SHALL stop and say so, because an untracked file that
  belongs in the release would otherwise be silently left out

### Requirement: File edits are portable
The in-place edits to the version file, the changelog and the vendor hashes
SHALL produce the same result on GNU and BSD userlands.

#### Scenario: Dating the changelog on macOS
- **WHEN** the script inserts the new version heading into `CHANGELOG.md` on a
  system with BSD sed
- **THEN** the heading SHALL be inserted on its own line, not with a literal `n`
  where the line break belongs

#### Scenario: The two vendor hashes stay in step
- **WHEN** the vendor hash is recomputed during a release
- **THEN** both `package.nix` and `flake.nix` SHALL receive the new value, since
  a stale hash in either one fails `nix flake check`
