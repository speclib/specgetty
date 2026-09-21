## 1. The grammar, transcribed

- [ ] 1.1 Add a code-fence mask to `src/ui/specparse.go` and make every heading
      test consult it, and verify a test asserts a `#### Scenario:` line inside a
      fenced block is content and a `# ` line inside one does not end its node
- [ ] 1.2 Recognise the requirements section as `/^##\s+Requirements\s*$/i` and
      the Purpose as a section titled `Purpose` case-insensitively, and verify a
      test covers `## requirements` and `## PURPOSE`
- [ ] 1.3 Recognise a requirement as `/^###\s+Requirement:\s*(.+)\s*$/i`, and
      verify a test asserts `### requirement: x` is one and `### Requirements`
      is not
- [ ] 1.4 Recognise a scenario as any non-fenced `/^####\s+/` heading, naming it
      by its heading text with an optional ATX closing run and an optional
      `Scenario:` prefix stripped, and verify tests cover `#### Edge case`,
      `#### Scenario: x` and `#### x ####`
- [ ] 1.5 Treat a level-four heading with no content before the next heading as
      not a scenario, and verify a test asserts it is absent from the outline
- [ ] 1.6 Report a requirement that sits outside the requirements section rather
      than reading it, and verify a test asserts the file does not structure and
      the reason names the section
- [ ] 1.7 Report a delta header in a main spec, and verify a test covers all four
      of ADDED, MODIFIED, REMOVED and RENAMED
- [ ] 1.8 Give every reason the line it was found on, and verify a test asserts
      the line number of a delta header on line 3 and of a stray requirement
      further down
- [ ] 1.9 Verify the grammar against openspec itself by a test that runs
      `openspec validate <spec> --type spec` for each fixture and asserts the two
      agree on whether the file is a spec, skipping if the CLI is absent

## 2. A scenario's content

- [ ] 2.1 Replace a scenario's `clauses` and `body` with an ordered list of parts,
      each a clause or a paragraph of prose, and verify a test asserts a scenario
      holding a clause, a paragraph and another clause keeps that order
- [ ] 2.2 Keep every line of a scenario whatever shape it is in, and verify tests
      assert the content of a bare-uppercase scenario, a bold-title-case one and
      a prose-only one all reach the node
- [ ] 2.3 Leave a horizontal rule out of the parts, and verify a test asserts a
      `---` between two clauses is absent while both clauses remain
- [ ] 2.4 Keep the clause test as it is, a bulleted `- **WHEN**` or `- WHEN`,
      case-sensitive, and verify a test asserts `- **When**` becomes prose rather
      than a clause and is still present
- [ ] 2.5 Verify no scenario in the local corpus renders as a bare title, by a
      test that walks every fixture and asserts each scenario node has at least
      one part

## 3. The card

- [ ] 3.1 Render a prose part as wrapped prose in its place among the clauses,
      and verify a test asserts a 229-character paragraph wraps and sits between
      the clauses that surround it
- [ ] 3.2 Keep the clause layout exactly as it is for a recognised clause, and
      verify the existing clause tests pass unchanged
- [ ] 3.3 Verify the card still fits its pane at 70, 80 and 120 columns with a
      scenario of mixed parts

## 4. The report

- [ ] 4.1 Return the reasons a file is not a spec rather than a single error, and
      verify a test asserts a file with three faults yields three reasons
- [ ] 4.2 Word each reason so it names the rule, the consequence and the fix,
      following openspec's own messages, and verify a test asserts the delta
      header reason says a delta header belongs in a change
- [ ] 4.3 Render the report in the spec view with each reason and its line, and
      verify a test asserts all three reasons and all three line numbers are on
      screen
- [ ] 4.4 Offer `E` from the report and verify a test asserts the file opened is
      the spec that would not structure
- [ ] 4.5 Verify the report fits the frame at 60x20 and at 200x50, and that no
      row exceeds the terminal

## 5. The level and the keys

- [ ] 5.1 Make `enter` on the specs tab descend in every case, and verify a test
      asserts the level changes for a file that does not fit
- [ ] 5.2 Give the spec level a second state so the view is either an outline and
      a card or a report, and verify a test asserts `docActive()` and the outline
      keys are inert in the report state
- [ ] 5.3 Make `esc` return to the specs tab from the report with the same spec
      selected, and verify a test round-trips
- [ ] 5.4 Stop putting the reason on the nav bar, and verify a test asserts
      `statusMsg` is empty after descending into a report
- [ ] 5.5 Add the report's arm to `renderNavBar`, listing `esc` and `E` and not
      the navigation keys, and verify a test asserts each
- [ ] 5.6 Re-parse on rescan in the report state too, so a file repaired in the
      editor structures itself without leaving the view, and verify a test
      asserts the view becomes an outline when the content is corrected

## 6. The fixture corpus

- [ ] 6.1 Add `src/ui/testdata/specs/` with one file per shape found in the
      survey: a delta-headed main spec, one with no Purpose, one with
      requirements outside the section, bare uppercase clauses, bold title-case
      clauses, a `####` header without the `Scenario:` prefix, a `####` block
      with no content, a heading inside a fence, and a prose-only scenario
- [ ] 6.2 Name each fixture for what it is and state in a comment which live
      project it was taken from, so the next reader knows it is a specimen rather
      than an invention
- [ ] 6.3 Drive the parser tests from the fixtures rather than from string
      literals in the test file, and verify each fixture's expected outcome is
      asserted
- [ ] 6.4 Keep the walk over this project's own live specs, and verify it still
      asserts every one structures

## 7. Verification

- [ ] 7.1 Revert task 2.2 and confirm the bare-uppercase fixture renders a bare
      title again. Check the build succeeds first
- [ ] 7.2 Revert task 1.7 and confirm the delta-headed fixture structures again
      when it should not
- [ ] 7.3 Verify the count of local specs that structure matches what the survey
      predicted, by a test that walks every spec under a directory given by an
      environment variable and is skipped when it is unset
- [ ] 7.4 `nix flake check` passes, coverage floors included

## 8. Documentation

- [ ] 8.1 Document in README.md what opens as an outline and what opens as a
      report, and that the rules are OpenSpec's own
- [ ] 8.2 Add the change to CHANGELOG.md under Fixed, naming the empty scenarios,
      and under Changed, naming the files that will now refuse to structure
