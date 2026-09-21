## 1. Reading a delta

- [x] 1.1 Add a delta entry point beside `parseSpec`, sharing `fenceMask`,
      `sectionBody`, `headingLevel`, `scenarioName` and `partsOf`, and verify a
      test parses a hand-written delta with all three operations into the
      expected tree
- [x] 1.2 Read a delta header as the operation for the requirements beneath it,
      carrying it on each requirement node, and verify a test asserts the
      operation of every requirement in a three-operation file
- [x] 1.3 Accept a requirement with no scenarios where the operation is a
      removal, and report one where it is not, and verify tests for both
- [x] 1.4 Make `## Purpose` optional and report a `## Requirements` section as a
      fault, that section being unread by OpenSpec in a change, and verify tests
      for both
- [x] 1.5 Keep an unknown operation heading rather than hiding what is under it,
      and verify a test with `## NEW Requirements`, which one archived change in
      this repository carries
- [x] 1.6 Verify the delta parser accepts every spec file of every change in this
      repository, by a test walking `openspec/changes/**/specs/**/*.md` and
      asserting the counts per operation

## 2. Comparing against the live spec

- [x] 2.1 Match a `MODIFIED` requirement to the live spec's requirement of the
      same heading, and verify a test covers a match and a miss
- [x] 2.2 Mark each scenario of a matched requirement unchanged, edited or added,
      and verify a test asserts all three against a hand-written pair
- [x] 2.3 Offer no marks and no comparison when there is no match, and verify a
      test asserts the absence rather than a default
- [x] 2.4 Offer no comparison for an archived change, and verify a test asserts
      that a delta equal to the live spec is not reported as altering nothing
- [x] 2.5 Verify a test asserts both sides come from memory, no file being read
      and no subprocess started

## 3. The outline

- [x] 3.1 List every delta file of the change as a root node in a stable order,
      and verify a test asserts the order is the same across two parses
- [x] 3.2 Draw the operation mark on a requirement row, and verify a test asserts
      the mark for added, modified, removed and unknown
- [x] 3.3 Draw the scenario mark inside a modified requirement and dim an
      unchanged scenario, and verify a test asserts the three marks and the
      dimming
- [x] 3.4 Keep the mark on the first row of a wrapped label and off its
      continuations, and verify a test asserts it at a width that wraps a
      63-column title

## 4. The card

- [x] 4.1 Render a capability node's card naming what the change does to it, and
      verify a test asserts the counts per operation
- [x] 4.2 Render a removal's card from the parts it carries, its reason and
      migration arriving as prose, and verify a test asserts both are shown in
      full
- [x] 4.3 Add the three-way row above the card for a node with an original,
      opening on the difference, and verify a test asserts the row and the
      opening view
- [x] 4.4 Omit the row for a node with no original, and verify a test asserts its
      absence for an added requirement, a removed one and a node of an archived
      change
- [x] 4.5 Render the difference at word level, the prose reflowing when it is
      edited, and verify a test asserts an added clause is marked and the
      untouched sentences around it are not
- [x] 4.6 Open each node on the difference rather than keeping the last choice,
      and verify a test moves between two comparable nodes and asserts both open
      on the difference

## 5. The level

- [x] 5.1 Add the level, entered by `enter` on a change's specs sub-tab and left
      by `esc` back to that sub-tab, and verify a test round-trips
- [x] 5.2 Leave `enter` inert on the artifact sub-tabs, and verify a test asserts
      the level is unchanged on proposal, design and tasks
- [x] 5.3 Give `left` and `right` the card's three-way row at this level, and
      verify a test asserts they move the row and stop at both ends
- [x] 5.4 Verify a test asserts `left` and `right` change neither the project tab
      bar nor the change's artifact sub-tab at this level, which is the rule
      `keep-the-arrows-in-the-spec-view` established one level across
- [x] 5.5 Add the level's arm to `docRegion`, `docActive`, `splitTab`,
      `listPage`, `moveListCursor`, `gotoListEnd` and `renderNavBar`, and verify a
      test walks each
- [x] 5.6 Keep the cursor by node path across a re-parse, and verify a test edits
      a delta above the cursor and asserts the cursor did not move

## 6. Opening the file

- [x] 6.1 Set `document.path` to the delta file the selected node belongs to, and
      verify a test asserts `E` opens the right file for a node in each of a
      multi-capability change's deltas
- [x] 6.2 Verify a test asserts the specs sub-tab one level up still offers no
      `E`, it still showing several files

## 7. Verification

- [x] 7.1 Assert the outline, the card and the paging behave as
      `spec-detail-view` requires, by running the same assertions against this
      view
- [x] 7.2 Revert 2.2 and confirm the scenario-mark tests fail. Check the build
      succeeds first
- [x] 7.3 `nix flake check` passes, coverage floors included

## 8. Documentation

- [x] 8.1 Document the level, its marks and its three-way row in README.md
- [x] 8.2 Add the change to CHANGELOG.md
