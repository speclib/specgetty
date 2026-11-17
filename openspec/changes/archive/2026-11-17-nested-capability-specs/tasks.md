# Tasks

## 1. Capability discovery

- [x] 1.1 Add one helper that walks a `specs/` directory and returns every
  directory that directly holds a `spec.md`, named by its path relative to that
  directory with `/` between segments, and continues past a subtree it cannot
  read. Verify with a scanner test covering: two nested capabilities, a container
  holding no `spec.md`, a directory holding both a `spec.md` and a nested
  capability, and a flat single-level tree returning the same names as before.

- [x] 1.2 Use the helper in `parseChangeDir` and in `ParseProjectInfo`, and make
  `SpecCount` the number of capabilities it returns. Rewrite
  `TestParseProjectInfo`, whose `cap-c` case asserts the old rule, into a
  container that is no longer counted, and add a nested fixture asserting
  `SpecCount` and `SpecNames`. Verify with `go test ./src/scanner/`.

- [x] 1.3 Add a test asserting that a change's nested delta names are the same
  strings as the project's nested spec names, which is what makes
  `buildChangeSpecTree` find an original. Verify it fails when either walk is
  reverted to one level.

## 2. Specs tab rendering

- [x] 2.1 Derive groups from the first path segment at render time and draw the
  spec list as a header per group carrying its count, the remaining path as the
  row, a blank line between groups, and capabilities without a parent segment as
  bare rows at their sort position. Truncate a header or row wider than the pane
  with `ansi.Truncate`. Verify with a render test: a nested project shows one
  header and its leaves, a flat project renders as it does today, and a narrow
  width does not widen the pane.

- [x] 2.2 Replace the spec list's own offset and page arithmetic with
  `itemLines` (`offsetFor`, `fit`), so the pane scrolls to the selected
  capability and a page moves by the rows it shows, headers and blank lines
  included. Verify with a paging test that the cursor never lands on a header or
  a blank line under `j`, `k`, `pgdown`, `pgup`, `gg` and `G`, and that the
  selected row is on screen after a jump to the last capability of a group.

- [x] 2.3 Add a `### Fixed` entry under `[Unreleased]` in `CHANGELOG.md` saying
  that nested capabilities are listed, grouped, and that a change's nested
  deltas no longer collapse. Verify the entry reads as one paragraph in the
  file's existing voice.

## 3. Integration

- [x] 3.1 Run `go test ./...` and `make build`; both succeed.

- [x] 3.2 Open `chatcms-app` in the built binary and confirm: the picker's
  specs column reads 52, the specs tab lists 14 groups whose rows open content,
  and `hosted-until-connected` shows three delta capability nodes with their
  MODIFIED marks drawn.

## Workflow follow-up

- Archive the change after review.
- Keep bean `specgetty-9xiu` linked; mark it completed once the archived result
  is verified.
