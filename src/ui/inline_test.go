package ui

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// The inline renderer draws one line of markdown. Every markdown surface goes
// through it: the specs tab, the detail card, a change's spec deltas, and the
// proposal and design documents.

// styleOpens counts the style-opening sequences in a rendered string, and
// resets counts the sequences that end one. A correctly drawn line closes each
// style before opening the next; a line that nests them does not.
var (
	reOpen  = regexp.MustCompile(`\x1b\[[0-9;]*[1-9][0-9;]*m`)
	reReset = regexp.MustCompile(`\x1b\[0?m`)
)

// TestACodeSpanHoldingMarksKeepsItsStyle is task 1.1, the live bug.
//
// `spec-detail-view/spec.md` quotes the clause format inside a code span, and so
// do nineteen other specs in the local corpus. The bold pass used to run first
// and style the marks inside the span; the code pass then wrapped the result,
// and the inner style's reset ended the cyan part way through the span.
func TestACodeSpanHoldingMarksKeepsItsStyle(t *testing.T) {
	const line = "A clause is one source bullet, `- **WHEN** the user...`."

	got := renderInlineMarkdown(line)

	if strings.Contains(ansi.Strip(got), "`") {
		t.Errorf("the backticks should be gone: %q", ansi.Strip(got))
	}
	// The whole span, drawn as one run. Asserted against the style rendering it
	// rather than by walking escapes, because the defect was precisely that the
	// run ended early.
	want := mdCodeStyle.Render("- **WHEN** the user...")
	if !strings.Contains(got, want) {
		t.Errorf("the span is not drawn as one run:\n  got  %q\n  want it to contain %q",
			got, want)
	}
}

func TestNoStyleIsEverNested(t *testing.T) {
	for _, line := range []string{
		"A clause is one source bullet, `- **WHEN** the user...`.",
		"plain `code` and **bold** and _italic_ together",
		"**bold holding `code` inside it**",
		"`code holding **marks** inside it`",
		"- **WHEN** the value is `a/path.yaml` and it SHALL be written",
	} {
		got := renderInlineMarkdown(line)

		// Walk the escape sequences: an opening may never follow an opening
		// without a reset between them.
		depth := 0
		for _, m := range regexp.MustCompile(`\x1b\[[0-9;]*m`).FindAllString(got, -1) {
			if reReset.MatchString(m) {
				depth--
				if depth < 0 {
					depth = 0
				}
				continue
			}
			depth++
			if depth > 1 {
				t.Errorf("a style opens inside another one:\n  in:  %s\n  raw: %q",
					line, got)
				break
			}
		}
	}
}

func TestTheExistingInlineBehaviourIsUnchanged(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"a `code span` in prose", "a code span in prose"},
		{"an `unclosed span in prose", "an `unclosed span in prose"},
		{"**bold** text", "bold text"},
		{"an **unclosed bold in prose", "an **unclosed bold in prose"},
		{"_italic_ text", "italic text"},
		{"no marks at all", "no marks at all"},
		{"", ""},
	} {
		if got := ansi.Strip(renderInlineMarkdown(c.in)); got != c.want {
			t.Errorf("%q rendered as %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEachMarkIsStyled(t *testing.T) {
	if got := renderInlineMarkdown("a `code span` in prose"); !strings.Contains(got, mdCodeStyle.Render("code span")) {
		t.Errorf("the code span is not styled: %q", got)
	}
	if got := renderInlineMarkdown("**bold** text"); !strings.Contains(got, mdBoldStyle.Render("bold")) {
		t.Errorf("the bold span is not styled: %q", got)
	}
	if got := renderInlineMarkdown("_italic_ text"); !strings.Contains(got, mdItalicStyle.Render("italic")) {
		t.Errorf("the italic span is not styled: %q", got)
	}
}

func TestACodeSpanBindsTighterThanBold(t *testing.T) {
	// Markdown reads a code span first, and its contents are never marks.
	got := renderInlineMarkdown("`code holding **marks** inside it`")

	if !strings.Contains(ansi.Strip(got), "**marks**") {
		t.Errorf("the marks inside a code span are content: %q", ansi.Strip(got))
	}
	if strings.Contains(got, mdBoldStyle.Render("marks")) {
		t.Errorf("the marks inside a code span were read as bold: %q", got)
	}
}

// --- 2.x the vocabulary ---

func TestAClauseKeywordAtTheStartOfALine(t *testing.T) {
	for _, c := range []struct {
		line  string
		word  string
		style lipgloss.Style
	}{
		{"- **WHEN** the user presses enter", "WHEN", kwConditionStyle},
		{"- WHEN the user presses enter", "WHEN", kwConditionStyle},
		{"WHEN the user presses enter", "WHEN", kwConditionStyle},
		{"  - **GIVEN** a project", "GIVEN", kwConditionStyle},
		{"- **THEN** it SHALL open", "THEN", kwAssertionStyle},
		{"THEN it opens", "THEN", kwAssertionStyle},
		{"- **AND** the other one too", "AND", kwContinuationStyle},
		{"* **WHEN** a star marks the item", "WHEN", kwConditionStyle},
		{"+ WHEN a plus marks it", "WHEN", kwConditionStyle},
	} {
		got := renderInlineMarkdown(c.line)
		if !strings.Contains(got, c.style.Render(c.word)) {
			t.Errorf("%q: %s is not drawn in its role's style:\n  %q",
				c.line, c.word, got)
		}
		if strings.Contains(ansi.Strip(got), "**") {
			t.Errorf("%q: the bold marks should be gone: %q", c.line, ansi.Strip(got))
		}
	}
}

func TestTheThreeClauseRolesAreToldApart(t *testing.T) {
	seen := map[string]string{}
	for _, word := range []string{"GIVEN", "WHEN", "THEN", "AND"} {
		seen[word] = clauseStyleFor(word).Render(word)
	}
	if seen["GIVEN"] == seen["THEN"] {
		t.Error("a condition and an assertion must read apart")
	}
	if seen["THEN"] == seen["AND"] {
		t.Error("an assertion and a continuation must read apart")
	}
	if seen["GIVEN"] == seen["AND"] {
		t.Error("a condition and a continuation must read apart")
	}
	// GIVEN and WHEN share a role, which is the one pair that may match.
	if clauseStyleFor("GIVEN").Render("x") != clauseStyleFor("WHEN").Render("x") {
		t.Error("GIVEN and WHEN both open a condition and share a style")
	}
}

func TestAClauseKeywordInTheMiddleOfASentenceIsProse(t *testing.T) {
	for _, line := range []string{
		"the value is written AND the file is closed",
		"a scenario may open with GIVEN or with WHEN",
		"reported THEN and only then",
	} {
		got := renderInlineMarkdown(line)
		for _, style := range []lipgloss.Style{kwConditionStyle, kwAssertionStyle, kwContinuationStyle} {
			for _, word := range []string{"GIVEN", "WHEN", "THEN", "AND"} {
				if strings.Contains(got, style.Render(word)) {
					t.Errorf("%q: %s was drawn as a keyword mid-sentence", line, word)
				}
			}
		}
	}
}

func TestANormativeKeywordAnywhereInProse(t *testing.T) {
	for _, c := range []struct{ line, word string }{
		{"the bar SHALL display the battery", "SHALL"},
		{"the project MUST include a package", "MUST"},
		{"SHALL be the first word too", "SHALL"},
		{"- **THEN** the change list SHALL be displayed", "SHALL"},
	} {
		got := renderInlineMarkdown(c.line)
		if !strings.Contains(got, kwBindingStyle.Render(c.word)) {
			t.Errorf("%q: %s is not drawn:\n  %q", c.line, c.word, got)
		}
	}
}

func TestATwoWordKeywordIsOneKeyword(t *testing.T) {
	for _, c := range []struct{ line, word string }{
		{"it SHALL NOT be displayed", "SHALL NOT"},
		{"the file MUST NOT be created", "MUST NOT"},
	} {
		got := renderInlineMarkdown(c.line)
		if !strings.Contains(got, kwBindingStyle.Render(c.word)) {
			t.Errorf("%q: %q is not drawn as one keyword:\n  %q", c.line, c.word, got)
		}
		// And never as a drawn first word beside a plain second one.
		first := strings.Fields(c.word)[0]
		if strings.Contains(got, kwBindingStyle.Render(first)+" NOT") {
			t.Errorf("%q: drawn as %s followed by a plain NOT", c.line, first)
		}
	}
}

func TestLowerCaseIsProse(t *testing.T) {
	for _, line := range []string{
		"the rule shall apply to everyone",
		"when the user presses enter",
		"then it opens",
	} {
		if got := renderInlineMarkdown(line); got != line {
			t.Errorf("%q was drawn: %q", line, got)
		}
	}
}

func TestAKeywordIsAWholeWord(t *testing.T) {
	for _, line := range []string{
		"the SHALLOW copy is enough",
		"ANDROID is not a keyword",
		"a MUSTER of keywords",
		"WHENEVER it happens",
	} {
		if got := renderInlineMarkdown(line); got != line {
			t.Errorf("%q had a word clipped into a keyword: %q", line, got)
		}
	}
}

func TestAKeywordInsideACodeSpanIsNotDrawn(t *testing.T) {
	for _, line := range []string{
		"a sentence about `SHALL` as a word",
		"the grammar reads `WHEN` at the start of a line",
		"`- **THEN** ...` is the shape of a clause",
	} {
		got := renderInlineMarkdown(line)
		for _, style := range []lipgloss.Style{
			kwBindingStyle, kwConditionStyle, kwAssertionStyle, kwContinuationStyle,
		} {
			for _, word := range []string{"SHALL", "WHEN", "THEN", "AND"} {
				if strings.Contains(got, style.Render(word)) {
					t.Errorf("%q: %s inside a code span was drawn as a keyword:\n  %q",
						line, word, got)
				}
			}
		}
	}
}

// --- 3.x the styles ---

func TestAKeywordIsDrawnAtItsOwnWidth(t *testing.T) {
	// A style must not change how wide a word is, or every column count in the
	// application is wrong by however many keywords a line holds.
	for _, line := range []string{
		"- **WHEN** the value SHALL NOT be written",
		"THEN it opens AND the cursor moves",
		"plain prose with no keywords",
	} {
		got := renderInlineMarkdown(line)
		want := strings.ReplaceAll(strings.ReplaceAll(line, "**", ""), "`", "")
		if ansi.StringWidth(got) != ansi.StringWidth(want) {
			t.Errorf("%q: drawn width %d, want %d\n  %q",
				line, ansi.StringWidth(got), ansi.StringWidth(want), got)
		}
	}
}

func TestPlainProseCarriesNoEscapes(t *testing.T) {
	const line = "a line with nothing in it to draw"
	if got := renderInlineMarkdown(line); got != line {
		t.Errorf("a plain line gained styling: %q", got)
	}
}

// --- 4.x every surface ---

func TestTheSpecsTabDrawsKeywordsInMarkdown(t *testing.T) {
	const src = "## Requirements\n\n### Requirement: A thing\nThe bar SHALL display it.\n\n" +
		"#### Scenario: Doing it\n- **WHEN** asked\n- **THEN** done\n"

	out := renderMarkdown(src, 60)

	for _, want := range []struct {
		word  string
		style lipgloss.Style
	}{
		{"SHALL", kwBindingStyle},
		{"WHEN", kwConditionStyle},
		{"THEN", kwAssertionStyle},
	} {
		if !strings.Contains(out, want.style.Render(want.word)) {
			t.Errorf("%s is not drawn in the markdown view", want.word)
		}
	}
}

func TestTheDeltaColoursAndTheKeywordColoursAreToldApart(t *testing.T) {
	// A spec delta shows what a change does to a requirement beside the
	// requirement's own clauses. Green and red mean added and removed there, so
	// no keyword may claim either.
	for _, kw := range []lipgloss.Style{
		kwConditionStyle, kwAssertionStyle, kwContinuationStyle, kwBindingStyle,
	} {
		for _, op := range []lipgloss.Style{opAddedStyle, opRemovedStyle, addedWordStyle} {
			if kw.Render("x") == op.Render("x") {
				t.Error("a keyword style is indistinguishable from a delta operation style")
			}
		}
	}
}

func TestACardDrawsABindingKeywordInsideAClause(t *testing.T) {
	n := specNode{kind: nodeScenario, title: "s", parts: []specPart{
		{kind: partClause, keyword: "THEN", text: "the file SHALL NOT be written"},
	}}

	card := renderSpecCard(n, 56)
	if !strings.Contains(card, kwBindingStyle.Render("SHALL NOT")) {
		t.Errorf("the binding keyword inside the clause is not drawn:\n%q", card)
	}
	if !strings.Contains(card, kwAssertionStyle.Render("THEN")) {
		t.Errorf("the clause keyword is not drawn in its role's style:\n%q", card)
	}
}

func TestACardTellsTheClauseRolesApart(t *testing.T) {
	n := specNode{kind: nodeScenario, title: "s", parts: []specPart{
		{kind: partClause, keyword: "WHEN", text: "a"},
		{kind: partClause, keyword: "THEN", text: "b"},
		{kind: partClause, keyword: "AND", text: "c"},
	}}

	card := renderSpecCard(n, 56)
	for _, c := range []struct {
		word  string
		style lipgloss.Style
	}{
		{"WHEN", kwConditionStyle}, {"THEN", kwAssertionStyle}, {"AND", kwContinuationStyle},
	} {
		if !strings.Contains(card, c.style.Render(c.word)) {
			t.Errorf("%s is not drawn in its role's style", c.word)
		}
	}
}

// TestProseTheCardDidNotLayOutKeepsItsKeywords is task 4.6, and the case the
// narrow clause parser leaves behind: a scenario written without bullets is
// drawn as prose, and its keywords are what give that prose its shape back.
func TestProseTheCardDidNotLayOutKeepsItsKeywords(t *testing.T) {
	tree, _ := mustParse(t, "bare", fixture(t, "bare-uppercase-clauses"))

	var n specNode
	for _, node := range tree.nodes {
		if node.title == "enable airplane mode" {
			n = node
		}
	}
	if len(n.parts) < 3 {
		t.Fatalf("the scenario has %d part(s); each keyword line is its own "+
			"paragraph, or the lines run together into one sentence", len(n.parts))
	}
	for _, p := range n.parts {
		if p.kind != partProse {
			t.Errorf("the clause parser stays narrow: this shape is prose, got kind %d", p.kind)
		}
	}

	card := renderSpecCard(n, 52)
	for _, c := range []struct {
		word  string
		style lipgloss.Style
	}{
		{"GIVEN", kwConditionStyle}, {"WHEN", kwConditionStyle}, {"THEN", kwAssertionStyle},
	} {
		if !strings.Contains(card, c.style.Render(c.word)) {
			t.Errorf("%s is not drawn in prose the card did not lay out:\n%s",
				c.word, ansi.Strip(card))
		}
	}
	if !strings.Contains(card, kwBindingStyle.Render("SHALL")) {
		t.Errorf("SHALL is not drawn in that prose:\n%s", ansi.Strip(card))
	}
}

// TestKeywordsChangeNoRowsWidth is the invariant this change owns: drawing a
// keyword must not move a single column, or every width the application
// computes is wrong by however many keywords a line holds.
//
// It is not an assertion that nothing overflows. `ansi.Wrap` lets a breakpoint
// character sit one column past the limit, which it did before this change and
// does identically after it.
func TestKeywordsChangeNoRowsWidth(t *testing.T) {
	dir := filepath.Join("..", "..", "openspec", "specs")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Skip(err)
	}

	var checked int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name(), "spec.md"))
		if err != nil {
			continue
		}
		checked++
		const width = 52
		for _, line := range strings.Split(renderMarkdown(string(b), width), "\n") {
			plain := ansi.Strip(line)
			if ansi.StringWidth(line) != ansi.StringWidth(plain) {
				t.Errorf("%s: a styled row measures %d and its text %d: %q",
					e.Name(), ansi.StringWidth(line), ansi.StringWidth(plain), plain)
				break
			}
		}
	}
	if checked == 0 {
		t.Fatal("no specs were read")
	}
}

func TestStylingDoesNotChangeWhereALineWraps(t *testing.T) {
	// The same content, drawn and undrawn, breaks in the same places.
	for _, src := range []string{
		"- **THEN** the content SHALL say that it is being read, and the interface SHALL stay responsive",
		"- **WHEN** a change is open and its `.openspec.yaml` records a schema",
		"The bar SHALL display the battery AND the clock, and SHALL NOT hide either",
	} {
		styled := wrapStyled(styleMarkdownLine(src), 52)
		plain := wrapStyled(strings.ReplaceAll(strings.ReplaceAll(src, "**", ""), "`", ""), 52)

		if len(styled) != len(plain) {
			t.Fatalf("%q wrapped to %d rows styled and %d plain", src, len(styled), len(plain))
		}
		for i := range styled {
			if ansi.Strip(styled[i]) != strings.TrimRight(plain[i], " ") &&
				strings.TrimSpace(ansi.Strip(styled[i])) != strings.TrimSpace(plain[i]) {
				t.Errorf("row %d differs:\n  styled %q\n  plain  %q",
					i, ansi.Strip(styled[i]), plain[i])
			}
		}
	}
}

// TestTheLocalCorpusRendersWithKeywords is task 5.3: every spec on this machine
// drawn through the renderer, asserting only what this change owns.
func TestTheLocalCorpusRendersWithKeywords(t *testing.T) {
	root := os.Getenv("SPECGETTY_CORPUS")
	if root == "" {
		t.Skip("set SPECGETTY_CORPUS to a directory of OpenSpec projects to run this")
	}

	var specs, drawn int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Base(path) != "spec.md" {
			return nil
		}
		if !strings.Contains(path, "/openspec/specs/") || strings.Contains(path, "/archive/") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		specs++
		const width = 52
		for _, line := range strings.Split(renderMarkdown(string(b), width), "\n") {
			if ansi.StringWidth(line) != ansi.StringWidth(ansi.Strip(line)) {
				t.Errorf("%s: a styled row measures differently from its text", path)
				return nil
			}
			if strings.Contains(line, "\x1b[") {
				drawn++
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d specs, %d drawn rows", specs, drawn)
	if specs == 0 {
		t.Fatal("no specs found")
	}
}

func TestClauseStyleForAWordThatIsNotAKeyword(t *testing.T) {
	// The card asks for a keyword's style by the word the parser found. A word
	// outside the vocabulary cannot reach it today, and if one ever does it
	// gets a style rather than an empty one.
	if got := clauseStyleFor("BECAUSE").Render("x"); got == "x" {
		t.Error("an unknown keyword should still be drawn")
	}
	for _, word := range []string{"GIVEN", "WHEN", "THEN", "AND"} {
		if clauseStyleFor(word).Render(word) == word {
			t.Errorf("%s is not drawn at all", word)
		}
	}
}

func TestAnEmptyRunIsNotDrawn(t *testing.T) {
	// withBinding emits the text on either side of a keyword, and either side
	// can be empty: a line that is nothing but a keyword.
	got := renderInlineMarkdown("SHALL")
	if got != kwBindingStyle.Render("SHALL") {
		t.Errorf("a line of one keyword drew as %q", got)
	}
	if renderInlineMarkdown("") != "" {
		t.Error("an empty line drew something")
	}
}
