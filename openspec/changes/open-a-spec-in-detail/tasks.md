## 1. The parser

- [ ] 1.1 Add `src/ui/specparse.go` with a node tree of Purpose, requirements and
      scenarios, each carrying its title and its body, and verify a unit test
      parses one hand-written spec into the expected tree
- [ ] 1.2 Parse a scenario's clauses into keyword and text, joining a bullet with
      its continuation lines, and verify a test covers `GIVEN`, `WHEN`, `THEN`,
      `AND` and a clause that spans three source lines
- [ ] 1.3 Report a parse failure when the file yields no requirements or a
      requirement has no scenarios, and verify tests cover both and that neither
      panics on an empty file
- [ ] 1.4 Verify the parser accepts all 24 live specs in this repository, by a
      test that walks `openspec/specs/*/spec.md` and asserts 140 requirements and
      428 scenarios in total
- [ ] 1.5 Give each node a path that identifies it across a re-parse, and verify a
      test finds the same node in a tree re-parsed from edited content

## 2. Inline markdown

- [ ] 2.1 Highlight backticked spans in `renderInlineMarkdown` and drop the
      backticks, and verify a test shows the span styled and the marks gone
- [ ] 2.2 Leave a line with an unclosed backtick unchanged, and verify a test
      asserts the rest of the line survives

## 3. The level

- [ ] 3.1 Add `levelSpec` and make `enter` on the specs tab descend into it when
      the selected spec parses, from either half, and verify tests for both halves
- [ ] 3.2 Make `enter` on a spec that does not parse set `statusMsg` and leave the
      level and the cursor alone, and verify a test asserts all three
- [ ] 3.3 Make `esc` return to the specs tab with the same spec selected and the
      list holding the keyboard, and verify a test round-trips
- [ ] 3.4 Make `enter` at `levelSpec` do nothing, and verify a test asserts the
      level is unchanged
- [ ] 3.5 Hold the parsed tree for the open spec and drop it on `esc`, and verify
      a test asserts the parser runs once across several renders

## 4. The outline

- [ ] 4.1 Add `specDetailSplit` at four tenths with a minimum, used by both the
      renderer and `docRegion`, and verify a test asserts the two agree at several
      widths
- [ ] 4.2 Lay the outline out with wrapped labels, requirements in
      `sectionHeaderStyle` and scenarios in the normal style, and verify a test
      asserts a 63-column title wraps at width 30 and is not cut
- [ ] 4.3 Map each node to its row range and highlight every row of the selected
      node, and verify a test asserts a two-row node is highlighted on both rows
- [ ] 4.4 Scroll the outline so the whole selected node is visible, and verify a
      test asserts a node whose last row is past the bottom brings both rows in
- [ ] 4.5 Move the cursor one node per `j` or `k` regardless of node height, and
      verify a test crosses a wrapped node in one step
- [ ] 4.6 Keep the cursor by node path across a re-parse, and verify tests for an
      edit above the cursor and for the node disappearing

## 5. The card

- [ ] 5.1 Render Purpose, a requirement's text and a scenario as the three card
      kinds, and verify a test asserts a requirement's card excludes its scenarios
- [ ] 5.2 Lay a clause out with its keyword on its own row and the text indented,
      wrapping with a hanging indent, and verify a test asserts a 229-character
      clause wraps with every row after the first still indented
- [ ] 5.3 Colour the keywords and reuse the inline styling for bold and code
      spans, and verify a test asserts a keyword and a code span are both styled
- [ ] 5.4 Shrink the padding as the pane narrows by one formula, and verify a test
      asserts no row exceeds the pane at 70, 80 and 120 columns
- [ ] 5.5 Feed the card through `currentDoc` keyed by project, spec and node path,
      and verify a test asserts moving to another node starts at the first row

## 6. The keys

- [ ] 6.1 Add the `levelSpec` arm to `docRegion`, `docActive` and `splitTab`, and
      verify a test asserts the card reports its position only while it has the
      keyboard
- [ ] 6.2 Make `tab` cycle outline, card and the log panel when it is open, and
      verify a test asserts the cycle in both directions of the open log
- [ ] 6.3 Teach `listPage`, `moveListCursor` and `gotoListEnd` the outline, with
      the page walking node heights until one pane of rows is covered, and verify
      tests for unequal heights and for all-one-row reducing to today's behaviour
- [ ] 6.4 Add the `levelSpec` arm to `renderNavBar`, and verify a test asserts it
      lists `esc`, `tab` and the scroll keys with the card focused

## 7. Documentation

- [ ] 7.1 Document the new level and its keys in README.md, and verify the keys
      listed match the ones bound in `renderNavBar`
- [ ] 7.2 Add the change to CHANGELOG.md
