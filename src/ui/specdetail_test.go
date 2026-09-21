package ui

import (
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

const twoReqSpec = `# demo Specification

## Purpose
Why this capability exists, in a sentence that names ` + "`a-thing`" + ` in code.

## Requirements

### Requirement: A requirement with a title long enough to wrap in a narrow pane
Its own prose, which the card shows and the scenarios do not.

#### Scenario: The first one
- **WHEN** something happens
- **THEN** something else does

#### Scenario: A scenario whose title is also long enough that it has to wrap somewhere
- **GIVEN** a clause whose text is long enough to need wrapping across more than one row of any reasonable pane, so that the hanging indent has something to hang
- **THEN** it still reads as one clause

### Requirement: The second requirement
More prose.

#### Scenario: Under the second
- **WHEN** a
- **THEN** b
`

// specDetailModel opens the demo spec at levelSpec.
func specDetailModel(t *testing.T) model {
	t.Helper()
	m := newModel(&scanner.Config{}, true, "0.0.0")
	m.width, m.height = 100, 24
	m.repoPaths = []string{"/p"}
	m.displayNames = []string{"p"}
	m.detailTab = tabSpecs
	m.focus = focusListPane
	m.fields = append([]string(nil), defaultFields...)
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Root: "/p", Origin: "/p",
		SpecCount:    2,
		SpecNames:    []string{"demo", "unparsable"},
		SpecContents: map[string]string{"demo": twoReqSpec, "unparsable": "just prose\n"},
	}}}
	m.recalcLayout()
	m.syncDocument()
	return m
}

// --- 3.1 to 3.5 the level ---

func TestEnterOpensTheSelectedSpec(t *testing.T) {
	for _, focus := range []int{focusListPane, focusContentPane} {
		m := specDetailModel(t)
		m.focus = focus
		opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})

		if opened.level != levelSpec {
			t.Fatalf("focus %d: level = %d, want levelSpec (%q)", focus, opened.level, opened.statusMsg)
		}
		if opened.specTree.name != "demo" {
			t.Errorf("opened %q", opened.specTree.name)
		}
		if opened.focus != focusListPane {
			t.Errorf("the outline holds the keyboard on open, got focus %d", opened.focus)
		}
		if opened.specNode != 0 {
			t.Errorf("the cursor starts on the first node, got %d", opened.specNode)
		}
	}
}

// TestEnterDescendsIntoAReport replaces the test that asserted the cursor
// stayed put and the nav bar carried the reason. A file usually fails in more
// than one way, and one transient line can hold neither the reasons nor the
// lines they sit on nor the key that opens an editor on them.
func TestEnterDescendsIntoAReport(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 1 // "unparsable"
	m.syncDocument()

	after := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})

	if after.level != levelSpec {
		t.Fatalf("level = %d, want the spec level", after.level)
	}
	if after.specStructured() {
		t.Error("the file does not fit the grammar, so there is no outline")
	}
	if len(after.specProblems) == 0 {
		t.Fatal("no reasons were given")
	}
	if after.specName != "unparsable" {
		t.Errorf("the report is about %q", after.specName)
	}
	if after.statusMsg != "" {
		t.Errorf("status = %q, want the nav bar left alone", after.statusMsg)
	}
}

func TestTheReportNamesEveryReasonAndItsLine(t *testing.T) {
	m := specDetailModel(t)
	// A file that fails in three ways at once, which is the usual case: 102 of
	// the local corpus carry a delta header, and most of those also lack a
	// Purpose and a requirements section.
	broken := "# thing\n\n## ADDED Requirements\n\n### Requirement: A\nIt SHALL.\n\n" +
		"#### Scenario: S\n- **WHEN** a\n- **THEN** b\n"
	info := m.projects["/p"].Info
	info.SpecContents = map[string]string{"demo": twoReqSpec, "unparsable": broken}
	m.projects["/p"] = scanner.ProjectStatus{Info: info}
	m.specCursor = 1
	m.syncDocument()

	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	opened.width, opened.height = 100, 30
	opened.recalcLayout()
	opened.syncDocument()

	if len(opened.specProblems) < 2 {
		t.Fatalf("got %d reasons, want every one: %v",
			len(opened.specProblems), opened.specProblems)
	}
	frame := ansi.Strip(opened.renderFrame())
	for _, want := range []string{"delta header", "Purpose", "unparsable"} {
		if !strings.Contains(frame, want) {
			t.Errorf("the report does not mention %q:\n%s", want, frame)
		}
	}
	// The delta header is on line 3 of that string.
	if !strings.Contains(frame, "line 3") {
		t.Errorf("the report does not give the line:\n%s", frame)
	}
}

func TestTheReportOffersTheEditor(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "nvim")
	started := fakeEditor(t, nil)

	m, _, root, _ := onDiskStoreModel(t, threeTasks)
	info := m.projects[m.repoPaths[0]].Info
	info.SpecContents = map[string]string{"some-capability": "# not a spec\n"}
	m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
	m.detailTab = tabSpecs
	m.recalcLayout()
	m.syncDocument()

	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if opened.specStructured() {
		t.Fatal("the fixture should not structure")
	}
	opened.width = 240
	if !strings.Contains(ansi.Strip(opened.renderNavBar()), "E edit") {
		t.Errorf("the report should offer E:\n%s", ansi.Strip(opened.renderNavBar()))
	}

	press(opened, tea.KeyPressMsg{Code: 'E', Text: "E"})
	if len(*started) != 1 {
		t.Fatalf("%d processes started, want 1", len(*started))
	}
	c := (*started)[0]
	want := filepath.Join(root, "openspec", "specs", "some-capability", "spec.md")
	if got := c.Args[len(c.Args)-1]; got != want {
		t.Errorf("opened %q, want the spec that would not structure %q", got, want)
	}
}

func TestTheReportScrollsAndTheOutlineKeysAreInert(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 1
	m.syncDocument()
	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	opened.width, opened.height = 80, 20
	opened.recalcLayout()
	opened.syncDocument()

	// One panel, so the vertical keys scroll it rather than moving an outline
	// that is not there.
	if !opened.docActive() {
		t.Error("the report is one panel, so it owns the vertical keys")
	}
	if opened.splitTab() {
		t.Error("there is no second half for tab to reach")
	}
	if opened.listPage() != 1 {
		t.Errorf("listPage = %d, want no outline page", opened.listPage())
	}

	after := press(opened, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if after.specNode != 0 {
		t.Errorf("j moved an outline cursor to %d", after.specNode)
	}
	if after.docHasCursor() {
		t.Error("a report has no task cursor")
	}
}

func TestAFileRepairedWhileOpenStructuresItself(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 1
	m.syncDocument()
	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if opened.specStructured() {
		t.Fatal("the fixture should open as a report")
	}

	// The reader fixes it in their editor; the rescan brings it back.
	info := opened.projects["/p"].Info
	info.SpecContents = map[string]string{"demo": twoReqSpec, "unparsable": twoReqSpec}
	opened.projects["/p"] = scanner.ProjectStatus{Info: info}
	opened.reparseOpenSpec()

	if !opened.specStructured() {
		t.Errorf("the repaired file should structure without leaving the view: %v",
			opened.specProblems)
	}
	if len(opened.specProblems) != 0 {
		t.Errorf("the reasons should be gone, got %v", opened.specProblems)
	}
}

func TestEscFromAReportReturnsToTheSpecsTab(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 1
	m.syncDocument()
	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})

	back := press(opened, tea.KeyPressMsg{Code: tea.KeyEscape})
	if back.level != levelProject || back.detailTab != tabSpecs {
		t.Errorf("level %d, tab %d, want the specs tab", back.level, back.detailTab)
	}
	if back.specCursor != 1 {
		t.Errorf("the same spec must be selected, got %d", back.specCursor)
	}
	if len(back.specProblems) != 0 || back.specName != "" {
		t.Error("the report is dropped on the way out")
	}
	// And the markdown is right there, which is why the report needs no
	// fallback of its own.
	if !strings.Contains(ansi.Strip(back.renderFrame()), "just prose") {
		t.Error("the whole file should be readable as markdown on the tab")
	}
}

func TestTheReportFitsTheFrame(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {100, 30}, {200, 50}} {
		m := specDetailModel(t)
		m.specCursor = 1
		m.syncDocument()
		opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
		opened.width, opened.height = size.w, size.h
		opened.recalcLayout()
		opened.syncDocument()

		lines := strings.Split(opened.renderFrame(), "\n")
		if len(lines) != size.h {
			t.Errorf("%dx%d: %d rows, want %d", size.w, size.h, len(lines), size.h)
		}
		for _, l := range lines {
			if w := ansi.StringWidth(l); w > size.w {
				t.Errorf("%dx%d: a row is %d columns wide", size.w, size.h, w)
				break
			}
		}
	}
}

func TestEscReturnsToTheSpecsTab(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 0
	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	back := press(opened, tea.KeyPressMsg{Code: tea.KeyEscape})

	if back.level != levelProject || back.detailTab != tabSpecs {
		t.Errorf("level %d, tab %d, want the specs tab", back.level, back.detailTab)
	}
	if back.specCursor != 0 {
		t.Errorf("the same spec must be selected, got %d", back.specCursor)
	}
	if back.focus != focusListPane {
		t.Errorf("the spec list holds the keyboard, got focus %d", back.focus)
	}
	if len(back.specTree.nodes) != 0 {
		t.Error("the tree is dropped on the way out")
	}
}

func TestEnterAtTheSpecLevelDoesNothing(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	again := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if again.level != levelSpec {
		t.Errorf("level = %d, want it unchanged", again.level)
	}
	if again.specNode != m.specNode {
		t.Errorf("the cursor moved to %d", again.specNode)
	}
}

// --- 4.1 the split ---

func TestTheSplitAgreesWithTheDocumentRegion(t *testing.T) {
	for _, width := range []int{60, 80, 100, 140} {
		m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		m.width, m.height = width, 24
		m.recalcLayout()

		_, cardOuter := specDetailSplit(m.panelContentWidth())
		docW, _ := m.docRegion()
		if docW != cardOuter-boxChrome {
			t.Errorf("width %d: docRegion is %d, the card's box holds %d",
				width, docW, cardOuter-boxChrome)
		}
	}
}

func TestTheSplitLeavesBothHalvesUsable(t *testing.T) {
	for _, width := range []int{60, 80, 120} {
		outline, card := specDetailSplit(width - 4)
		if outline+card+1 != width-4 {
			t.Errorf("width %d: halves and gap are %d, want %d", width, outline+card+1, width-4)
		}
		if outline < 20 || card < 20 {
			t.Errorf("width %d: outline %d, card %d; both must stay usable", width, outline, card)
		}
	}
}

// --- 4.2 to 4.5 the outline ---

func TestALongLabelWrapsRatherThanBeingCut(t *testing.T) {
	tree, problems := parseSpec("demo", twoReqSpec)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	rows := outlineRows(tree, 30)

	// The long requirement title occupies more than one row and survives whole.
	var joined string
	for _, r := range rows {
		if tree.nodes[r.node].kind == nodeRequirement &&
			strings.Contains(tree.nodes[r.node].title, "long enough to wrap") {
			joined += strings.TrimSpace(r.text) + " "
		}
	}
	for _, word := range strings.Fields(tree.nodes[1].title) {
		if !strings.Contains(joined, word) {
			t.Errorf("the label lost %q; it must wrap rather than be cut:\n%s", word, joined)
		}
	}
	if strings.Contains(joined, "…") {
		t.Error("nothing may be elided")
	}
}

func TestEveryRowOfTheSelectedNodeIsHighlighted(t *testing.T) {
	tree, problems := parseSpec("demo", twoReqSpec)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	// Node 1 is the long requirement title, which wraps at this width.
	rows := outlineRows(tree, 30)
	first, last := rowRangeOfNode(rows, 1)
	if last <= first {
		t.Fatalf("node 1 occupies rows %d..%d; the test needs a wrapped node", first, last)
	}

	out := renderSpecOutline(tree, 1, 30, len(rows), true)
	const highlight = "\x1b[30;42m"
	n := strings.Count(out, highlight)
	if n != last-first+1 {
		t.Errorf("got %d highlighted rows, want %d", n, last-first+1)
	}
}

func TestTheOutlineScrollsAWholeNodeIntoView(t *testing.T) {
	tree, problems := parseSpec("demo", twoReqSpec)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	rows := outlineRows(tree, 30)

	// A pane too short for the tree, with the cursor on a node whose last row
	// is past the bottom.
	for node := range tree.nodes {
		const height = 4
		out := ansi.Strip(renderSpecOutline(tree, node, 30, height, true))
		lines := strings.Split(out, "\n")
		if len(lines) > height {
			t.Fatalf("node %d: %d lines, want at most %d", node, len(lines), height)
		}
		first, last := rowRangeOfNode(rows, node)
		want := strings.TrimSpace(rows[last].text)
		if last-first+1 > height {
			// Taller than the pane: the top wins.
			want = strings.TrimSpace(rows[first].text)
		}
		var seen bool
		for _, l := range lines {
			if strings.Contains(l, want) {
				seen = true
			}
		}
		if !seen {
			t.Errorf("node %d: %q is off screen:\n%s", node, want, out)
		}
	}
}

func TestTheCursorMovesOneNodePerKeystroke(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	total := len(m.specTree.nodes)
	if total < 4 {
		t.Fatalf("the demo spec has %d nodes; the test needs several", total)
	}

	for i := 1; i < total; i++ {
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
		if m.specNode != i {
			t.Fatalf("after %d presses the cursor is on %d", i, m.specNode)
		}
	}
	// And stops rather than running off.
	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.specNode != total-1 {
		t.Errorf("the cursor ran to %d", m.specNode)
	}
	for i := total - 2; i >= 0; i-- {
		m = press(m, tea.KeyPressMsg{Code: 'k', Text: "k"})
		if m.specNode != i {
			t.Fatalf("stepping back, the cursor is on %d, want %d", m.specNode, i)
		}
	}
}

// --- 4.6 the cursor survives a re-parse ---

func TestTheCursorSurvivesAnEditAboveIt(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	for i := 0; i < 3; i++ {
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	wantPath := m.specTree.nodes[m.specNode].path
	wantTitle := m.specTree.nodes[m.specNode].title

	// A scenario inserted above the cursor, which shifts every index below it.
	// Prose rewritten above the cursor would leave the indices alone, and a
	// cursor remembered by index would survive that by accident.
	edited := strings.Replace(twoReqSpec, "#### Scenario: The first one",
		"#### Scenario: An inserted one\n- **WHEN** x\n- **THEN** y\n\n"+
			"#### Scenario: The first one", 1)
	info := m.projects["/p"].Info
	info.SpecContents = map[string]string{"demo": edited, "unparsable": "just prose\n"}
	m.projects["/p"] = scanner.ProjectStatus{Info: info}
	m.reparseOpenSpec()

	if got := m.specTree.nodes[m.specNode].path; got != wantPath {
		t.Errorf("the cursor moved to %q, want %q", got, wantPath)
	}
	if got := m.specTree.nodes[m.specNode].title; got != wantTitle {
		t.Errorf("the cursor is on %q, want %q", got, wantTitle)
	}
}

func TestTheCursorClampsWhenItsNodeIsGone(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	for i := 0; i < len(m.specTree.nodes)-1; i++ {
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}

	// The last requirement and its scenario removed.
	cut := twoReqSpec[:strings.Index(twoReqSpec, "### Requirement: The second requirement")]
	info := m.projects["/p"].Info
	info.SpecContents = map[string]string{"demo": cut, "unparsable": "just prose\n"}
	m.projects["/p"] = scanner.ProjectStatus{Info: info}
	m.reparseOpenSpec()

	if m.specNode >= len(m.specTree.nodes) {
		t.Errorf("the cursor is on %d of %d nodes", m.specNode, len(m.specTree.nodes))
	}
	if _, ok := m.currentDoc(); !ok {
		t.Error("the card must still render")
	}
}

// --- 5.x the card ---

func TestARequirementsCardExcludesItsScenarios(t *testing.T) {
	tree, _ := mustParse(t, "demo", twoReqSpec)
	card := ansi.Strip(renderSpecCard(tree.nodes[1], 60))

	if !strings.Contains(card, "Its own prose") {
		t.Errorf("the requirement's own prose belongs on it:\n%s", card)
	}
	if strings.Contains(card, "The first one") || strings.Contains(card, "WHEN") {
		t.Errorf("its scenarios are nodes of their own:\n%s", card)
	}
}

func TestAClauseWrapsWithAHangingIndent(t *testing.T) {
	tree, _ := mustParse(t, "demo", twoReqSpec)
	// The long GIVEN clause.
	var scenario specNode
	for _, n := range tree.nodes {
		if strings.HasPrefix(n.title, "A scenario whose title") {
			scenario = n
		}
	}
	card := ansi.Strip(renderSpecCard(scenario, 60))
	lines := strings.Split(card, "\n")

	var inClause bool
	var rows int
	for _, l := range lines {
		if strings.TrimSpace(l) == "GIVEN" {
			inClause = true
			continue
		}
		if inClause {
			if strings.TrimSpace(l) == "" {
				break
			}
			rows++
			if !strings.HasPrefix(l, "    ") {
				t.Errorf("a clause row is not indented: %q", l)
			}
		}
	}
	if rows < 2 {
		t.Errorf("the long clause wrapped to %d rows, want several", rows)
	}
}

func TestTheCardStylesKeywordsAndCodeSpans(t *testing.T) {
	tree, _ := mustParse(t, "demo", twoReqSpec)
	card := renderSpecCard(tree.nodes[0], 60) // Purpose, which names a-thing in code

	if !strings.Contains(card, mdCodeStyle.Render("a-thing")) {
		t.Error("a backticked span must be styled and its marks dropped")
	}
	if strings.Contains(ansi.Strip(card), "`") {
		t.Error("the backticks must be gone")
	}

	scen := renderSpecCard(tree.nodes[2], 60)
	if !strings.Contains(scen, kwConditionStyle.Render("WHEN")) {
		t.Error("a clause keyword must be styled")
	}
}

func TestTheCardFitsItsPaneAtAnyWidth(t *testing.T) {
	tree, _ := mustParse(t, "demo", twoReqSpec)
	for _, width := range []int{70, 80, 120} {
		_, cardOuter := specDetailSplit(width - 4)
		inner := cardOuter - boxChrome
		for i, n := range tree.nodes {
			for _, l := range strings.Split(ansi.Strip(renderSpecCard(n, inner)), "\n") {
				if len([]rune(l)) > inner {
					t.Errorf("width %d, node %d: a row is %d columns in a %d pane: %q",
						width, i, len([]rune(l)), inner, l)
				}
			}
		}
	}
}

func TestMovingToAnotherNodeStartsAtTheTop(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	first := m.docKey
	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})

	if m.docKey == first {
		t.Error("another node is another document")
	}
	if m.docViewport.YOffset() != 0 {
		t.Errorf("it starts at the top, got offset %d", m.docViewport.YOffset())
	}
}

// --- 6.x the keys ---

func TestTabMovesBetweenTheOutlineAndTheCard(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.focus != focusListPane {
		t.Fatalf("focus %d, want the outline", m.focus)
	}
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusContentPane {
		t.Errorf("tab moves to the card, got focus %d", m.focus)
	}
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusListPane {
		t.Errorf("and back to the outline, got focus %d", m.focus)
	}
}

func TestTheCardReportsItsPositionOnlyWhileItHasTheKeyboard(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.docActive() {
		t.Error("the outline holds the keyboard, so the card does not own the axis")
	}
	onCard := press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if !onCard.docActive() {
		t.Error("it does once the card holds it")
	}
}

func TestPagingTheOutlineAndTheCard(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	m.width, m.height = 100, 12
	m.recalcLayout()

	// The outline pages by nodes while it holds the keyboard.
	paged := press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	if paged.specNode == m.specNode {
		t.Error("pgdown must move the outline")
	}
	if paged.docViewport.YOffset() != 0 {
		t.Error("and must not scroll the card")
	}

	// G reaches the last node, gg the first.
	end := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
	if end.specNode != len(m.specTree.nodes)-1 {
		t.Errorf("G went to %d, want the last node", end.specNode)
	}
	top := press(press(end, tea.KeyPressMsg{Code: 'g', Text: "g"}), tea.KeyPressMsg{Code: 'g', Text: "g"})
	if top.specNode != 0 {
		t.Errorf("gg went to %d, want the first", top.specNode)
	}
}

func TestListPageWithEveryNodeOneRowTall(t *testing.T) {
	// The page walks node heights, which with one-row nodes reduces to the row
	// count, the same thing every other list does.
	src := `## Purpose
Long enough a sentence to count as the purpose of a capability.

## Requirements

### Requirement: r
x

#### Scenario: a
- **WHEN** x
- **THEN** y

#### Scenario: b
- **WHEN** x
- **THEN** y
`
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	tree, problems := parseSpec("short", src)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	m.specTree = tree
	m.width, m.height = 120, 24
	m.recalcLayout()

	outlineOuter, _ := specDetailSplit(m.panelContentWidth())
	rows := m.mainPanelHeight() - 1 - boxRows
	if got := nodesInRows(tree, outlineOuter-boxChrome, 0, rows); got != len(tree.nodes) {
		t.Errorf("got %d nodes in %d rows, want all %d: every node is one row here",
			got, rows, len(tree.nodes))
	}
}

// --- 6.4 the nav bar ---

func TestTheNavBarNamesTheSpecLevelKeys(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	m.width = 200

	onOutline := ansi.Strip(m.renderNavBar())
	for _, want := range []string{"esc", "back to specs", "tab", "navigate"} {
		if !strings.Contains(onOutline, want) {
			t.Errorf("the outline's nav bar is missing %q:\n%s", want, onOutline)
		}
	}

	onCard := ansi.Strip(press(m, tea.KeyPressMsg{Code: tea.KeyTab}).renderNavBar())
	for _, want := range []string{"esc", "tab", "^f^b", "gg/G"} {
		if !strings.Contains(onCard, want) {
			t.Errorf("the card's nav bar is missing %q:\n%s", want, onCard)
		}
	}
}

// --- the frame ---

func TestTheSpecViewFitsTheFrame(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {92, 30}, {120, 50}} {
		m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		m.width, m.height = size.w, size.h
		m.recalcLayout()
		m.syncDocument()

		lines := strings.Split(m.renderFrame(), "\n")
		if len(lines) != size.h {
			t.Errorf("%dx%d: got %d rows, want %d", size.w, size.h, len(lines), size.h)
		}
		for _, l := range lines {
			if w := ansi.StringWidth(l); w > size.w {
				t.Errorf("%dx%d: a row is %d columns wide", size.w, size.h, w)
				break
			}
		}
	}
}

func TestTheViewNamesTheSpec(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	m.syncDocument()
	if !strings.Contains(ansi.Strip(m.renderFrame()), "demo") {
		t.Error("the view must never be anonymous")
	}
}

// --- 2.x inline markdown ---

func TestBackticksAreStyledAndDropped(t *testing.T) {
	got := renderInlineMarkdown("a `code span` in prose")
	if !strings.Contains(got, mdCodeStyle.Render("code span")) {
		t.Errorf("the span must be styled: %q", got)
	}
	if strings.Contains(ansi.Strip(got), "`") {
		t.Errorf("the marks must be gone: %q", ansi.Strip(got))
	}
}

func TestAnUnclosedBacktickIsLeftAlone(t *testing.T) {
	got := ansi.Strip(renderInlineMarkdown("an `unclosed span in prose"))
	if got != "an `unclosed span in prose" {
		t.Errorf("got %q, want the line unchanged", got)
	}
}

// TestTheTreeIsParsedOnceNotPerRender asserts the tree is held rather than
// re-derived. A hand-edited title survives several renders only if nothing on
// the draw path parses the file again.
func TestTheTreeIsParsedOnceNotPerRender(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	m.specTree.nodes[1].title = "SENTINEL-TITLE"

	for i := 0; i < 3; i++ {
		m.syncDocument()
		if !strings.Contains(ansi.Strip(m.renderFrame()), "SENTINEL-TITLE") {
			t.Fatalf("render %d lost the held tree; something re-parses on the draw path", i)
		}
	}
	if m.specTree.nodes[1].title != "SENTINEL-TITLE" {
		t.Error("the tree was replaced")
	}
}

func TestListPageWalksUnequalNodeHeights(t *testing.T) {
	tree, problems := parseSpec("demo", twoReqSpec)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	// Narrow enough that the two long titles wrap and the short ones do not.
	const width, rows = 28, 4
	all := outlineRows(tree, width)

	var heights []int
	for node := range tree.nodes {
		first, last := rowRangeOfNode(all, node)
		heights = append(heights, last-first+1)
	}
	var tall int
	for _, h := range heights {
		if h > 1 {
			tall++
		}
	}
	if tall == 0 {
		t.Fatalf("no node wraps at width %d; heights %v", width, heights)
	}

	// From the top, a page covers one pane of rows, so with a wrapped node in
	// the way it moves fewer nodes than there are rows.
	got := nodesInRows(tree, width, 0, rows)
	if got >= rows {
		t.Errorf("moved %d nodes in %d rows; a wrapped node must cost more than one",
			got, rows)
	}
	var used int
	for i := 0; i < got; i++ {
		used += heights[i]
	}
	if used > rows {
		t.Errorf("the page covered %d rows, more than the pane's %d", used, rows)
	}
	// And it always moves, even when the node under the cursor is taller than
	// the whole pane.
	if n := nodesInRows(tree, 12, 1, 1); n < 1 {
		t.Errorf("a page must always move at least one node, got %d", n)
	}
}

func TestTheSplitGivesTheCardRoomOnANarrowTerminal(t *testing.T) {
	// Below about 46 columns the outline's floor and the card's minimum cannot
	// both be met. The outline gives way rather than the card collapsing to
	// nothing, because a card of one column shows no spec at all.
	for _, width := range []int{30, 36, 44} {
		outline, card := specDetailSplit(width)
		if outline < 1 || card < 1 {
			t.Errorf("width %d: outline %d, card %d; neither may vanish", width, outline, card)
		}
		if outline+card+1 != width {
			t.Errorf("width %d: the halves and the gap come to %d", width, outline+card+1)
		}
		if card < outline {
			t.Errorf("width %d: card %d, outline %d; the card gets what is left",
				width, card, outline)
		}
	}
}

func TestReparseLeavesThingsAloneWhenThereIsNothingToDo(t *testing.T) {
	t.Run("not at the spec level", func(t *testing.T) {
		m := specDetailModel(t)
		m.reparseOpenSpec()
		if len(m.specTree.nodes) != 0 {
			t.Error("no tree is parsed below the spec level")
		}
	})

	t.Run("the spec is gone from the scan", func(t *testing.T) {
		m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		before := m.specTree
		info := m.projects["/p"].Info
		info.SpecContents = map[string]string{}
		m.projects["/p"] = scanner.ProjectStatus{Info: info}

		m.reparseOpenSpec()
		if len(m.specTree.nodes) != len(before.nodes) {
			t.Error("a spec that vanished from the scan keeps the tree it opened with")
		}
	})

	t.Run("edited into a shape that will not parse", func(t *testing.T) {
		m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		before := len(m.specTree.nodes)
		info := m.projects["/p"].Info
		info.SpecContents = map[string]string{"demo": "## Purpose\nhalf typed\n"}
		m.projects["/p"] = scanner.ProjectStatus{Info: info}

		m.reparseOpenSpec()
		if len(m.specTree.nodes) != before {
			t.Errorf("the view emptied mid-read: %d nodes, want %d",
				len(m.specTree.nodes), before)
		}
		if _, ok := m.currentDoc(); !ok {
			t.Error("the card must still render")
		}
	})
}

func TestOpeningASpecWithNothingSelectedIsInert(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 99 // past the end, as a filtered list can leave it

	after := m.openSelectedSpec()
	if after.level != levelProject {
		t.Errorf("level = %d, want it unchanged", after.level)
	}
	if after.statusMsg != "" {
		t.Errorf("status = %q, want nothing reported", after.statusMsg)
	}

	empty := newModel(&scanner.Config{}, true, "0.0.0")
	if got := empty.openSelectedSpec(); got.level != levelProject {
		t.Error("with no project at all, enter must do nothing")
	}
}

func TestTheSpecLevelObeysTheMinimumTerminalSize(t *testing.T) {
	// Below the documented minimum the view refuses rather than drawing a
	// squeezed outline, the same as every other level.
	for _, size := range []struct{ w, h int }{{59, 30}, {100, 19}, {20, 4}} {
		m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		m.width, m.height = size.w, size.h
		if got := m.View().Content; !strings.Contains(got, "Terminal too small") {
			t.Errorf("at %dx%d the spec level did not refuse: %q", size.w, size.h, got)
		}
	}
}

// --- the arrow keys at the spec level ---

func TestArrowKeysDoNothingInTheSpecView(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{
		{Code: tea.KeyRight}, {Code: tea.KeyLeft},
	} {
		m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		before := m

		after := press(m, key)

		if after.detailTab != before.detailTab {
			t.Errorf("%v moved the tab bar from %d to %d, and it is not on screen",
				key, before.detailTab, after.detailTab)
		}
		if after.focus != before.focus {
			t.Errorf("%v moved the keyboard from focus %d to %d",
				key, before.focus, after.focus)
		}
		if after.level != levelSpec {
			t.Errorf("%v left the level, now %d", key, after.level)
		}
		if after.specNode != before.specNode {
			t.Errorf("%v moved the outline cursor to %d", key, after.specNode)
		}
	}
}

// TestTheCardKeepsTheKeyboardAcrossAnArrow is the symptom a reader sees: the
// keyboard jumping back to the outline mid-read.
func TestTheCardKeepsTheKeyboardAcrossAnArrow(t *testing.T) {
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	onCard := press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if onCard.focus != focusContentPane {
		t.Fatal("the test needs the card holding the keyboard")
	}

	for _, key := range []tea.KeyPressMsg{{Code: tea.KeyRight}, {Code: tea.KeyLeft}} {
		after := press(onCard, key)
		if after.focus != focusContentPane {
			t.Errorf("%v took the keyboard off the card", key)
		}
		// The lit border follows the focus, so it must not have moved either.
		frame := ansi.Strip(after.renderFrame())
		if frame != ansi.Strip(onCard.renderFrame()) {
			t.Errorf("%v changed what is drawn", key)
		}
	}
}

func TestTheSpecsTabIsActiveAfterArrowsAndEsc(t *testing.T) {
	// This passed before the fix, because esc forces the tab on the way out,
	// which is what masked the defect. It must still pass.
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	for i := 0; i < 4; i++ {
		m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
		m = press(m, tea.KeyPressMsg{Code: tea.KeyLeft})
	}

	back := press(m, tea.KeyPressMsg{Code: tea.KeyEscape})
	if back.detailTab != tabSpecs {
		t.Errorf("tab = %d, want the specs tab", back.detailTab)
	}
	if back.level != levelProject {
		t.Errorf("level = %d, want the project view", back.level)
	}
}

func TestTheArrowsStillWorkWhereTheyBelong(t *testing.T) {
	t.Run("the project tab bar", func(t *testing.T) {
		m := specDetailModel(t)
		m.detailTab = tabChanges
		m.recalcLayout()
		m.syncDocument()

		right := press(m, tea.KeyPressMsg{Code: tea.KeyRight})
		if right.detailTab != tabSpecs {
			t.Errorf("right went to tab %d, want the next one", right.detailTab)
		}
		if left := press(right, tea.KeyPressMsg{Code: tea.KeyLeft}); left.detailTab != tabChanges {
			t.Errorf("left went to tab %d, want back", left.detailTab)
		}
	})

	t.Run("artifact sub-tabs in an open change", func(t *testing.T) {
		// This fixture carries proposal.md and tasks.md, so there is somewhere
		// for the keys to go.
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		m.level = levelChange
		m.changeArtifactTab = 0
		m.recalcLayout()
		m.syncDocument()
		r, ok := m.selectedRow()
		if !ok {
			t.Fatal("no change selected")
		}
		if n := len(r.artifactTabNames()); n < 2 {
			t.Fatalf("the fixture has %d artifact tab(s); the test needs two", n)
		}

		right := press(m, tea.KeyPressMsg{Code: tea.KeyRight})
		if right.changeArtifactTab != 1 {
			t.Errorf("right went to sub-tab %d, want 1", right.changeArtifactTab)
		}
		if left := press(right, tea.KeyPressMsg{Code: tea.KeyLeft}); left.changeArtifactTab != 0 {
			t.Errorf("left went to sub-tab %d, want 0", left.changeArtifactTab)
		}
	})
}

// mustParse is the two-value parse the detail tests want, since the shape of a
// spec is specparse_test.go's subject rather than theirs.
func mustParse(t *testing.T, name, content string) (specTree, []specProblem) {
	t.Helper()
	tree, problems := parseSpec(name, content)
	if len(problems) > 0 {
		t.Fatalf("%s: %v", name, problems)
	}
	return tree, nil
}

func TestASpanThatWrapsKeepsItsStyle(t *testing.T) {
	// Styling used to run after wrapping, so a backticked span straddling the
	// wrap became two halves with one backtick each: the marks stayed on screen
	// and the colour never arrived.
	n := specNode{kind: nodeScenario, title: "s", parts: []specPart{{
		kind:    partClause,
		keyword: "THEN",
		text: "the value SHALL be written to " +
			"`a/path/that/is/long/enough/to/straddle/the/wrap.yaml` and read back",
	}}}

	card := renderSpecCard(n, 46)
	if strings.Contains(ansi.Strip(card), "`") {
		t.Errorf("a backtick survived the wrap:\n%s", ansi.Strip(card))
	}
	if !strings.Contains(card, "\x1b[") {
		t.Error("the span was not styled at all")
	}
}

func TestTheReportStylesItsOwnSpans(t *testing.T) {
	problems := []specProblem{{line: 3, text: "`## ADDED Requirements` is a delta " +
		"header, and a main spec keeps its requirements under `## Requirements`."}}

	out := renderSpecReport("thing", problems, 50)
	if strings.Contains(ansi.Strip(out), "`") {
		t.Errorf("the report shows its backticks:\n%s", ansi.Strip(out))
	}
	if !strings.Contains(ansi.Strip(out), "line 3") {
		t.Errorf("the line is missing:\n%s", ansi.Strip(out))
	}
}

// --- 3.x prose on the card ---

func TestProseSitsWhereItWasWritten(t *testing.T) {
	tree, _ := mustParse(t, "mixed", fixture(t, "mixed-parts"))
	var n specNode
	for _, node := range tree.nodes {
		if node.title == "Clauses around a paragraph" {
			n = node
		}
	}

	card := ansi.Strip(renderSpecCard(n, 60))
	first := strings.Index(card, "the first thing happens")
	para := strings.Index(card, "Rationale")
	last := strings.Index(card, "the second thing happens")
	if first < 0 || para < 0 || last < 0 {
		t.Fatalf("something is missing from the card:\n%s", card)
	}
	if !(first < para && para < last) {
		t.Errorf("the paragraph is not between the clauses:\n%s", card)
	}
}

func TestALongParagraphWraps(t *testing.T) {
	long := strings.Repeat("a paragraph of prose that has to wrap somewhere. ", 5)
	n := specNode{kind: nodeScenario, title: "s", parts: []specPart{
		{kind: partProse, text: strings.TrimSpace(long)},
	}}
	if len(strings.TrimSpace(long)) < 229 {
		t.Fatalf("the test needs a paragraph over 229 characters, got %d", len(long))
	}

	card := ansi.Strip(renderSpecCard(n, 60))
	var rows int
	for _, l := range strings.Split(card, "\n") {
		if strings.Contains(l, "paragraph of prose") {
			rows++
		}
		if len([]rune(l)) > 60 {
			t.Errorf("a row is %d columns in a 60 pane: %q", len([]rune(l)), l)
		}
	}
	if rows < 3 {
		t.Errorf("the paragraph wrapped to %d rows, want several:\n%s", rows, card)
	}
}

func TestACardOfMixedPartsFitsItsPane(t *testing.T) {
	tree, _ := mustParse(t, "mixed", fixture(t, "mixed-parts"))
	for _, width := range []int{70, 80, 120} {
		_, cardOuter := specDetailSplit(width - 4)
		inner := cardOuter - boxChrome
		for i, n := range tree.nodes {
			for _, l := range strings.Split(ansi.Strip(renderSpecCard(n, inner)), "\n") {
				if len([]rune(l)) > inner {
					t.Errorf("width %d, node %d: a row is %d columns in a %d pane: %q",
						width, i, len([]rune(l)), inner, l)
				}
			}
		}
	}
}

func TestAScenarioInAnyShapeShowsSomething(t *testing.T) {
	// The whole point of the change, asserted through the card rather than the
	// parser: no shape in the corpus renders as a bare title.
	for _, name := range []string{
		"canonical", "bare-uppercase-clauses", "bold-title-clauses",
		"prose-scenario", "mixed-parts", "heading-without-prefix", "heading-in-fence",
	} {
		tree, _ := mustParse(t, name, fixture(t, name))
		for _, n := range tree.nodes {
			if n.kind != nodeScenario {
				continue
			}
			card := ansi.Strip(renderSpecCard(n, 56))
			body := strings.TrimSpace(strings.Replace(card, "Scenario: "+n.title, "", 1))
			if body == "" {
				t.Errorf("%s: scenario %q renders as a bare title", name, n.title)
			}
		}
	}
}

func TestTheReportsNavBarOffersOnlyWhatAReportCanDo(t *testing.T) {
	m := specDetailModel(t)
	m.specCursor = 1
	m.syncDocument()
	opened := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	opened.width = 240
	opened.recalcLayout()
	opened.syncDocument()

	bar := ansi.Strip(opened.renderNavBar())
	for _, want := range []string{"esc back to specs", "E edit", "scroll", "^f^b"} {
		if !strings.Contains(bar, want) {
			t.Errorf("the report's nav bar is missing %q:\n%s", want, bar)
		}
	}
	for _, unwanted := range []string{"tab focus", "navigate"} {
		if strings.Contains(bar, unwanted) {
			t.Errorf("the report's nav bar offers %q, which does nothing here:\n%s",
				unwanted, bar)
		}
	}
}
