package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The fixtures in testdata/specs are specimens, not inventions. Each was taken
// from a live spec on the machine this was written on, and its header comment
// says which. The test that missed the bug this change fixes walked only this
// project's own specs, every one written in one shape by one author.
func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "specs", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func parseFixture(t *testing.T, name string) specTree {
	t.Helper()
	tree, problems := parseSpec(name, fixture(t, name))
	if len(problems) > 0 {
		t.Fatalf("%s should be a spec, got %d problems: %v", name, len(problems), problems)
	}
	return tree
}

// --- 1.x the grammar ---

// TestEveryFixtureLandsWhereItShould is the table the rest of the file leans on.
func TestEveryFixtureLandsWhereItShould(t *testing.T) {
	for _, c := range []struct {
		name    string
		isSpec  bool
		because string
	}{
		{"canonical", true, ""},
		{"bare-uppercase-clauses", true, ""},
		{"bold-title-clauses", true, ""},
		{"prose-scenario", true, ""},
		{"mixed-parts", true, ""},
		{"heading-without-prefix", true, ""},
		{"heading-in-fence", true, ""},
		{"empty-scenario", true, ""},
		{"delta-header", false, "delta header"},
		{"no-purpose", false, "`## Purpose`"},
		{"requirements-outside", false, "outside the `## Requirements` section"},
		{"no-requirements", false, "no requirements"},
		{"requirement-without-scenario", false, "no scenario"},
	} {
		t.Run(c.name, func(t *testing.T) {
			tree, problems := parseSpec(c.name, fixture(t, c.name))

			if c.isSpec {
				if len(problems) > 0 {
					t.Fatalf("should be a spec, got: %v", problems)
				}
				if req, scen := tree.counts(); req == 0 || scen == 0 {
					t.Errorf("parsed to %d requirements and %d scenarios", req, scen)
				}
				return
			}

			if len(problems) == 0 {
				t.Fatal("should not be a spec, but parsed")
			}
			if len(tree.nodes) != 0 {
				t.Errorf("a file that is not a spec yields no tree, got %d nodes",
					len(tree.nodes))
			}
			var found bool
			for _, p := range problems {
				if strings.Contains(p.text, c.because) {
					found = true
				}
			}
			if !found {
				t.Errorf("no problem mentions %q; got %v", c.because, problems)
			}
		})
	}
}

func TestAFenceHidesItsHeadings(t *testing.T) {
	tree := parseFixture(t, "heading-in-fence")

	for _, n := range tree.nodes {
		if n.title == "not a requirement" || n.title == "not a scenario" {
			t.Errorf("a line inside a fenced block became a %q node", n.title)
		}
	}
	if req, scen := tree.counts(); req != 1 || scen != 1 {
		t.Errorf("got %d requirements and %d scenarios, want 1 and 1", req, scen)
	}

	// And the fenced lines are still part of the scenario's content, being
	// content rather than structure.
	var joined string
	for _, n := range tree.nodes {
		if n.kind == nodeScenario {
			for _, p := range n.parts {
				joined += p.text + " "
			}
		}
	}
	if !strings.Contains(joined, "Page Title") {
		t.Errorf("the fenced block's content was dropped:\n%s", joined)
	}
}

func TestTheHeadingsAreCaseInsensitive(t *testing.T) {
	src := strings.NewReplacer(
		"## Purpose", "## purpose",
		"## Requirements", "## REQUIREMENTS",
		"### Requirement: The first thing", "### requirement: The first thing",
	).Replace(fixture(t, "canonical"))

	tree, problems := parseSpec("x", src)
	if len(problems) > 0 {
		t.Fatalf("the grammar is case-insensitive on its headings: %v", problems)
	}
	if tree.indexOfPath("req/The first thing") < 0 {
		t.Error("the lower-case requirement heading was not read")
	}
}

func TestARequirementsHeadingIsNotARequirement(t *testing.T) {
	// `### Requirements` is not `### Requirement:`.
	src := strings.Replace(fixture(t, "canonical"),
		"### Requirement: The second thing", "### Requirements of a kind", 1)
	tree, problems := parseSpec("x", src)
	if len(problems) > 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
	for _, n := range tree.nodes {
		if n.kind == nodeRequirement && strings.HasPrefix(n.title, "of a kind") {
			t.Error("a heading that is not `### Requirement:` became one")
		}
	}
}

func TestAnyLevelFourHeadingIsAScenario(t *testing.T) {
	tree := parseFixture(t, "heading-without-prefix")

	want := []string{"Parsing", "Rendering", "Edge case"}
	var got []string
	for _, n := range tree.nodes {
		if n.kind == nodeScenario {
			got = append(got, n.title)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("scenarios %v, want %v: openspec counts any `#### ` heading, "+
			"and strips a `Scenario:` prefix and an ATX closing run from the name",
			got, want)
	}
}

func TestAScenarioHeadingWithNoContentIsNotAScenario(t *testing.T) {
	tree := parseFixture(t, "empty-scenario")

	for _, n := range tree.nodes {
		if n.title == "Never filled in" {
			t.Error("a heading with nothing under it became a scenario; openspec " +
				"does not count it")
		}
	}
	if _, scen := tree.counts(); scen != 1 {
		t.Errorf("got %d scenarios, want the one that has content", scen)
	}
}

func TestADeltaHeaderIsReported(t *testing.T) {
	for _, kind := range []string{"ADDED", "MODIFIED", "REMOVED", "RENAMED"} {
		src := strings.Replace(fixture(t, "delta-header"),
			"## ADDED Requirements", "## "+kind+" Requirements", 1)
		_, problems := parseSpec("x", src)

		var found bool
		for _, p := range problems {
			if strings.Contains(p.text, "delta header") && strings.Contains(p.text, kind) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: no problem names it as a delta header; got %v", kind, problems)
		}
	}
}

func TestEveryProblemNamesItsLine(t *testing.T) {
	_, problems := parseSpec("x", fixture(t, "delta-header"))

	var deltaLine int
	for _, p := range problems {
		if strings.Contains(p.text, "delta header") {
			deltaLine = p.line
		}
	}
	// The comment block, the title, a blank, then the delta header.
	want := 0
	for i, l := range strings.Split(fixture(t, "delta-header"), "\n") {
		if strings.HasPrefix(l, "## ADDED") {
			want = i + 1
		}
	}
	if deltaLine != want {
		t.Errorf("the delta header was reported on line %d, want %d", deltaLine, want)
	}

	_, outside := parseSpec("x", fixture(t, "requirements-outside"))
	var lines string
	for _, p := range outside {
		if strings.Contains(p.text, "outside") {
			lines = p.text
			if p.line == 0 {
				t.Error("a problem about particular lines must name one")
			}
		}
	}
	if !strings.Contains(lines, "line ") {
		t.Errorf("the reason should list the lines: %q", lines)
	}
}

func TestAProblemSaysWhatToDo(t *testing.T) {
	_, problems := parseSpec("x", fixture(t, "delta-header"))

	var text string
	for _, p := range problems {
		if strings.Contains(p.text, "delta header") {
			text = p.text
		}
	}
	for _, want := range []string{"change", "## Requirements", "invisible"} {
		if !strings.Contains(text, want) {
			t.Errorf("the reason should mention %q, so the reader can fix it: %q",
				want, text)
		}
	}
}

// TestTheGrammarAgreesWithOpenspec is task 1.9. The tool of record decides
// whether a file is a spec; this asserts specgetty has not invented a rule of
// its own, nor missed one.
func TestTheGrammarAgreesWithOpenspec(t *testing.T) {
	if _, err := exec.LookPath("openspec"); err != nil {
		t.Skip("openspec is not installed")
	}

	entries, err := os.ReadDir(filepath.Join("testdata", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".md")
		t.Run(name, func(t *testing.T) {
			// A project holding just this spec, so openspec validates it alone.
			root := t.TempDir()
			dir := filepath.Join(root, "openspec", "specs", name)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			body := fixture(t, name)
			if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command("openspec", "validate", name, "--type", "spec")
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "NO_COLOR=1", "GIT_TERMINAL_PROMPT=0")
			out, _ := cmd.CombinedOutput()
			theirs := !strings.Contains(string(out), "[ERROR]")

			_, problems := parseSpec(name, body)
			ours := len(problems) == 0

			if ours != theirs {
				t.Errorf("specgetty says spec=%v, openspec says spec=%v\nopenspec:\n%s\nspecgetty: %v",
					ours, theirs, out, problems)
			}
		})
	}
}

// --- 2.x a scenario's content ---

func TestNoScenarioIsEmptyWhateverItsShape(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("testdata", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	var checked int
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".md")
		tree, problems := parseSpec(name, fixture(t, name))
		if len(problems) > 0 {
			continue
		}
		for _, n := range tree.nodes {
			if n.kind != nodeScenario {
				continue
			}
			checked++
			if len(n.parts) == 0 {
				t.Errorf("%s: scenario %q has no content; it would render as a "+
					"bare title, which is the bug this change removes", name, n.title)
			}
		}
	}
	if checked < 10 {
		t.Fatalf("checked %d scenarios, expected the fixture corpus to hold more", checked)
	}
}

func TestBareUppercaseClausesAreKept(t *testing.T) {
	tree := parseFixture(t, "bare-uppercase-clauses")

	var n specNode
	for _, node := range tree.nodes {
		if node.title == "enable airplane mode" {
			n = node
		}
	}
	if len(n.parts) == 0 {
		t.Fatal("the scenario is empty, which is bean specgetty-vfnn exactly")
	}
	var joined string
	for _, p := range n.parts {
		joined += p.text + " "
	}
	for _, want := range []string{"airplane mode is off", "clicks the airplane mode toggle", "rfkill block all"} {
		if !strings.Contains(joined, want) {
			t.Errorf("%q was dropped:\n%s", want, joined)
		}
	}
}

func TestBoldTitleCaseClausesAreKept(t *testing.T) {
	tree := parseFixture(t, "bold-title-clauses")

	for _, n := range tree.nodes {
		if n.kind != nodeScenario {
			continue
		}
		if len(n.parts) == 0 {
			t.Errorf("scenario %q is empty", n.title)
		}
	}

	// They are kept as prose, the clause layout being deliberately narrow.
	var n specNode
	for _, node := range tree.nodes {
		if node.title == "Initialize Docusaurus Package" {
			n = node
		}
	}
	var joined string
	for _, p := range n.parts {
		if p.kind == partClause {
			t.Errorf("`**Given**` is not laid out as a clause by design; got keyword %q",
				p.keyword)
		}
		joined += p.text + " "
	}
	if !strings.Contains(joined, "the monorepo root") {
		t.Errorf("the content was dropped:\n%s", joined)
	}
}

func TestAProseOnlyScenarioIsKept(t *testing.T) {
	tree := parseFixture(t, "prose-scenario")

	var n specNode
	for _, node := range tree.nodes {
		if node.title == "Site settings (Future)" {
			n = node
		}
	}
	if len(n.parts) != 1 || n.parts[0].kind != partProse {
		t.Fatalf("want one prose part, got %d parts: %+v", len(n.parts), n.parts)
	}
	if !strings.Contains(n.parts[0].text, "deferred to a future change") {
		t.Errorf("the paragraph was truncated: %q", n.parts[0].text)
	}
}

func TestClausesAndProseKeepTheirOrder(t *testing.T) {
	tree := parseFixture(t, "mixed-parts")

	var n specNode
	for _, node := range tree.nodes {
		if node.title == "Clauses around a paragraph" {
			n = node
		}
	}
	if len(n.parts) != 3 {
		t.Fatalf("want three parts, got %d: %+v", len(n.parts), n.parts)
	}
	if n.parts[0].kind != partClause || n.parts[0].keyword != "WHEN" {
		t.Errorf("part 0 = %+v, want the WHEN clause", n.parts[0])
	}
	if n.parts[1].kind != partProse || !strings.Contains(n.parts[1].text, "Rationale") {
		t.Errorf("part 1 = %+v, want the paragraph between them", n.parts[1])
	}
	if n.parts[2].kind != partClause || n.parts[2].keyword != "THEN" {
		t.Errorf("part 2 = %+v, want the THEN clause", n.parts[2])
	}
}

func TestAHorizontalRuleIsNotContent(t *testing.T) {
	tree := parseFixture(t, "mixed-parts")

	for _, n := range tree.nodes {
		for _, p := range n.parts {
			if strings.Contains(p.text, "---") {
				t.Errorf("scenario %q kept a separator as content: %q", n.title, p.text)
			}
		}
	}
	// And the requirement after the rule is still there.
	if tree.indexOfPath("req/The one after the rule") < 0 {
		t.Error("the requirement after the rule was lost with it")
	}
}

func TestTheClauseTestStaysNarrow(t *testing.T) {
	// Deliberately unchanged: a bulleted upper-case keyword. Since nothing is
	// dropped any more, which shapes are laid out is presentation.
	for _, c := range []struct {
		line    string
		clause  bool
		keyword string
	}{
		{"- **WHEN** a thing", true, "WHEN"},
		{"- WHEN a thing", true, "WHEN"},
		{"- **GIVEN** a thing", true, "GIVEN"},
		{"- **AND** a thing", true, "AND"},
		{"- **When** a thing", false, ""},
		{"**WHEN** a thing", false, ""},
		{"WHEN a thing", false, ""},
		{"- something else", false, ""},
	} {
		kw, _, ok := clauseOf(c.line)
		if ok != c.clause {
			t.Errorf("%q: clause=%v, want %v", c.line, ok, c.clause)
		}
		if ok && kw != c.keyword {
			t.Errorf("%q: keyword %q, want %q", c.line, kw, c.keyword)
		}
	}
}

func TestAClauseSpanningSeveralSourceLines(t *testing.T) {
	tree := parseFixture(t, "canonical")

	var n specNode
	for _, node := range tree.nodes {
		if node.title == "A clause over three lines" {
			n = node
		}
	}
	if len(n.parts) != 2 {
		t.Fatalf("want two clauses, got %d: %+v", len(n.parts), n.parts)
	}
	want := "a clause whose text runs on across a second source line and a third"
	if n.parts[0].text != want {
		t.Errorf("got %q, want %q", n.parts[0].text, want)
	}
	if n.parts[1].keyword != "AND" {
		t.Errorf("the second clause is %+v", n.parts[1])
	}
}

func TestEveryClauseKeywordIsRead(t *testing.T) {
	src := `## Purpose
Long enough to be a purpose of a capability.

## Requirements

### Requirement: k
The system SHALL do it.

#### Scenario: all four
- **GIVEN** a
- **WHEN** b
- **THEN** c
- **AND** d
`
	tree, problems := parseSpec("k", src)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	var got []string
	for _, n := range tree.nodes {
		for _, p := range n.parts {
			got = append(got, p.keyword)
		}
	}
	if strings.Join(got, ",") != "GIVEN,WHEN,THEN,AND" {
		t.Errorf("keywords %v", got)
	}
}

// --- node paths ---

func TestNodePathsSurviveAReParse(t *testing.T) {
	tree := parseFixture(t, "canonical")
	path := tree.nodes[len(tree.nodes)-1].path

	edited := strings.Replace(fixture(t, "canonical"),
		"### Requirement: The first thing",
		"### Requirement: A new one\nIts prose.\n\n#### Scenario: Inserted\n- **WHEN** x\n- **THEN** y\n\n### Requirement: The first thing", 1)
	again, problems := parseSpec("thing", edited)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	if at := again.indexOfPath(path); at < 0 {
		t.Errorf("the node at %q is gone after an edit above it", path)
	}

	renamed := strings.Replace(fixture(t, "canonical"),
		"### Requirement: The second thing", "### Requirement: Renamed", 1)
	third, problems := parseSpec("thing", renamed)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	if third.indexOfPath("req/The second thing") >= 0 {
		t.Error("a renamed requirement keeps its old path")
	}
}

func TestNodePathsAreUnique(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("testdata", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".md")
		tree, problems := parseSpec(name, fixture(t, name))
		if len(problems) > 0 {
			continue
		}
		seen := map[string]bool{}
		for _, n := range tree.nodes {
			if seen[n.path] {
				t.Errorf("%s: duplicate path %q", name, n.path)
			}
			seen[n.path] = true
		}
	}
}

// --- the live corpus ---

// TestTheParserAcceptsEveryLiveSpec keeps the walk over this project's own
// specs. It is not enough on its own, which is what the fixtures are for, but
// it is what will notice the next shape this project starts writing.
func TestTheParserAcceptsEveryLiveSpec(t *testing.T) {
	dir := filepath.Join("..", "..", "openspec", "specs")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Skip(err)
	}

	var specs, reqs, scens int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name(), "spec.md")
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		specs++
		tree, problems := parseSpec(e.Name(), string(b))
		if len(problems) > 0 {
			t.Errorf("%s is not read as a spec: %v", e.Name(), problems)
			continue
		}
		r, s := tree.counts()
		reqs += r
		scens += s

		for _, n := range tree.nodes {
			if n.kind == nodeScenario && len(n.parts) == 0 {
				t.Errorf("%s: scenario %q is empty", e.Name(), n.title)
			}
		}
	}

	if specs == 0 {
		t.Fatal("no specs found")
	}
	if reqs < specs || scens < reqs {
		t.Errorf("%d specs, %d requirements, %d scenarios: every spec has at "+
			"least one requirement and every requirement at least one scenario",
			specs, reqs, scens)
	}
	t.Logf("%d specs, %d requirements, %d scenarios", specs, reqs, scens)
}

// TestTheLocalCorpusIsRead is task 7.3: a sweep over a directory of real specs
// outside this repository, skipped unless one is named.
func TestTheLocalCorpusIsRead(t *testing.T) {
	root := os.Getenv("SPECGETTY_CORPUS")
	if root == "" {
		t.Skip("set SPECGETTY_CORPUS to a directory of OpenSpec projects to run this")
	}

	var total, structured, blank int
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
		total++
		tree, problems := parseSpec(filepath.Base(filepath.Dir(path)), string(b))
		if len(problems) > 0 {
			return nil
		}
		structured++
		for _, n := range tree.nodes {
			if n.kind == nodeScenario && len(n.parts) == 0 {
				blank++
				t.Errorf("%s: scenario %q is empty", path, n.title)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d specs, %d structured, %d blank scenarios", total, structured, blank)
	if blank > 0 {
		t.Errorf("%d scenarios would render as a bare title", blank)
	}
}

func TestManyLinesAreSummarised(t *testing.T) {
	// A file with a dozen requirements outside the section should not produce a
	// dozen line numbers on one row of a report.
	var b strings.Builder
	b.WriteString("## Purpose\nLong enough a sentence to count as a purpose.\n\n")
	for i := 0; i < 9; i++ {
		b.WriteString("### Requirement: R\nIt SHALL.\n\n#### Scenario: S\n- **WHEN** a\n- **THEN** b\n\n")
	}
	b.WriteString("## Requirements\n\n### Requirement: Inside\nIt SHALL.\n\n" +
		"#### Scenario: S\n- **WHEN** a\n- **THEN** b\n")

	_, problems := parseSpec("x", b.String())
	var text string
	for _, p := range problems {
		if strings.Contains(p.text, "outside") {
			text = p.text
		}
	}
	if text == "" {
		t.Fatalf("no problem about requirements outside the section: %v", problems)
	}
	if !strings.Contains(text, "and 5 more") {
		t.Errorf("nine lines should be summarised after the fourth: %q", text)
	}
	if strings.Count(text, "line ") != 4 {
		t.Errorf("want four lines named, got %d: %q", strings.Count(text, "line "), text)
	}
}

func TestAnEmptyPurposeIsNoPurpose(t *testing.T) {
	// openspec reads an empty Purpose section as no Purpose: `purpose || ''`.
	src := "## Purpose\n\n## Requirements\n\n### Requirement: A\nIt SHALL.\n\n" +
		"#### Scenario: S\n- **WHEN** a\n- **THEN** b\n"
	_, problems := parseSpec("x", src)

	var found bool
	for _, p := range problems {
		if strings.Contains(p.text, "Purpose") {
			found = true
			if p.line == 0 {
				t.Error("an empty section has a line to point at")
			}
		}
	}
	if !found {
		t.Errorf("an empty Purpose was accepted: %v", problems)
	}
}

func TestCarriageReturnsAndABOMDoNotHideTheHeadings(t *testing.T) {
	src := "\ufeff" + strings.ReplaceAll(
		"## Purpose\nLong enough a sentence to count as a purpose here.\n\n"+
			"## Requirements\n\n### Requirement: A\nIt SHALL.\n\n"+
			"#### Scenario: S\n- **WHEN** a\n- **THEN** b\n", "\n", "\r\n")

	_, problems := parseSpec("x", src)
	if len(problems) > 0 {
		t.Errorf("a CRLF file with a BOM should read the same: %v", problems)
	}
}
