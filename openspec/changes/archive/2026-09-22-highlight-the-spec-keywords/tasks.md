## 1. The renderer, rebuilt before anything is added to it

- [x] 1.1 Write a failing test for the live bug: a line containing a code span
      that itself contains bold marks, as `spec-detail-view/spec.md` line 119
      does, and assert the whole span is one style. Confirm it fails
- [x] 1.2 Replace the wrapping passes in `renderInlineMarkdown` with one pass
      that scans the line and emits each segment with exactly one style, and
      verify every existing inline markdown test passes unchanged
- [x] 1.3 Read a code span before bold and italics, and verify a test asserts
      `` `- **WHEN** x` `` is drawn wholly as a code span with no second style
      inside it
- [x] 1.4 Keep an unclosed mark literal, and verify the existing unclosed
      backtick test passes and a new one covers an unclosed `**`
- [x] 1.5 Verify no style is ever nested, by a test that renders a line with
      every kind of mark and asserts the output contains no style opening before
      a previous one has closed

## 2. The vocabulary

- [x] 2.1 Add the keyword set: `GIVEN` and `WHEN` as a condition, `THEN` as an
      assertion, `AND` as a continuation, and `SHALL`, `SHALL NOT`, `MUST` and
      `MUST NOT` as binding, and verify a table test covers each word
- [x] 2.2 Match the two-word keywords before the one-word ones, and verify a test
      asserts `SHALL NOT` is one keyword and never a drawn `SHALL` beside a
      plain `NOT`
- [x] 2.3 Recognise clause keywords at the start of a line only, bare, as a list
      item, or wrapped in bold marks, and verify tests cover all three forms and
      assert a mid-sentence `AND` is left as prose
- [x] 2.4 Recognise normative keywords anywhere in prose, and verify a test
      asserts one mid-sentence
- [x] 2.5 Match upper case only, and verify a test asserts lower-case `shall` and
      `when` are left as prose
- [x] 2.6 Match on word boundaries, and verify a test asserts `SHALLOW` and
      `ANDROID` are untouched
- [x] 2.7 Leave a keyword inside a code span alone, and verify a test asserts a
      sentence about `` `SHALL` `` draws one style and not two

## 3. The styles

- [x] 3.1 Add the four styles, condition bold yellow, assertion bold blue,
      continuation yellow unbolded, binding bold magenta, and verify a test
      asserts each keyword renders with its own
- [x] 3.2 Verify the four are distinguishable from one another, by a test
      asserting no two of the four render identically
- [x] 3.3 Verify a keyword's drawn width equals the word's width, so that no
      column arithmetic anywhere is thrown by a style

## 4. Every surface

- [x] 4.1 Verify the specs tab draws keywords in a spec's markdown, by a test
      over a rendered document
- [x] 4.2 Verify a change's spec deltas draw them, and that the delta operation
      colours and the keyword colours are still told apart on the same screen
- [x] 4.3 Verify a proposal and a design draw them, since one renderer draws
      every document
- [x] 4.4 Draw the card's clause keyword by its role rather than one style for
      all four, and verify a test asserts a `THEN` row differs from a `WHEN` row
- [x] 4.5 Verify a `SHALL` inside a clause's text is drawn, by a test over a card
- [x] 4.6 Verify a scenario the card draws as prose still has its line-initial
      keywords drawn, using the bare-uppercase fixture, which is the case the
      narrow clause parser leaves behind. This needed a change the task did not
      anticipate: the card joined a scenario's consecutive source lines into one
      paragraph, as markdown says to, so only the first keyword was at a line
      start and the other two were mid-sentence. A prose line that opens with a
      clause keyword now starts a paragraph of its own. The lines stay prose and
      the clause parser stays narrow; they just keep the boundaries the author
      wrote

## 5. Verification

- [x] 5.1 Revert 1.3 and confirm 1.1 fails again. Check the build succeeds first.
      The first attempt at this revert was a no-op: swapping the order of the
      code-span and bold branches changes nothing, because a position is either
      a backtick or an asterisk and never both. What carries the fix is
      consuming a span whole, so the revert makes the span's contents pass
      through bold handling, which reproduces the original bug exactly and fails
      all three tests
- [x] 5.2 Revert 2.2 and confirm the `SHALL NOT` test fails
- [x] 5.3 Verify every spec on this machine renders without a panic and with no
      row wider than its pane, by a test over a directory given by an environment
      variable and skipped when it is unset
- [x] 5.4 Verify the frame still fits at 60x20, 100x30 and 200x50 on the specs
      tab and in the detail card, covered by the frame tests those views already
      have. The stronger assertion this change owns is that drawing a keyword
      moves no column: a styled row measures exactly as its text does, asserted
      over every live spec. An assertion that no row ever exceeds its pane would
      fail for a reason this change did not cause: `ansi.Wrap` lets a breakpoint
      character sit one column past the limit, identically before and after
- [x] 5.5 `nix flake check` passes, coverage floors included

## 6. Documentation

- [x] 6.1 Document the keywords and what each colour means in README.md, naming
      `openspec.nvim` as where the vocabulary comes from
- [x] 6.2 Add the change to CHANGELOG.md, the renderer fix under Fixed and the
      keywords under Added
