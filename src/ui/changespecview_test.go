package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// changeSpecModel opens a change whose deltas touch two capabilities, one of
// them modifying a requirement the project's specs still hold.
func changeSpecModel(t *testing.T, archived bool) model {
	t.Helper()
	m := newModel(&scanner.Config{}, true, "0.0.0")
	m.width, m.height = 100, 30
	m.repoPaths = []string{"/p"}
	m.displayNames = []string{"p"}
	m.detailTab = tabChanges
	m.focus = focusDetail
	m.fields = append([]string(nil), defaultFields...)

	ci := scanner.ChangeInfo{
		Name:             "a-change",
		DirName:          "a-change",
		ArtifactFiles:    []string{"proposal.md", "tasks.md"},
		ArtifactContents: map[string]string{"proposal.md": "## Why\nBecause.\n", "tasks.md": "- [ ] 1.1 Do it\n"},
		SpecNames:        []string{"a-capability", "b-capability"},
		SpecContents: map[string]string{
			"a-capability": modifyingDelta,
			"b-capability": threeOps,
		},
	}
	info := scanner.ProjectInfo{
		Root: "/p", Origin: "/p",
		SpecNames:    []string{"a-capability"},
		SpecContents: map[string]string{"a-capability": liveSpec},
	}
	if archived {
		ci.DirName = "2026-01-01-a-change"
		info.ArchivedChanges = []scanner.ChangeInfo{ci}
	} else {
		info.Changes = []scanner.ChangeInfo{ci}
	}
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: info}}
	m.recalcLayout()
	m.syncCursor()

	// Open the change and move to its specs sub-tab.
	m = press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.level != levelChange {
		t.Fatalf("the change did not open, level = %d", m.level)
	}
	for m.changeArtifactTab < m.changeSpecsTabIndex() {
		m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	}
	m.syncDocument()
	return m
}

func openedChangeSpecs(t *testing.T, archived bool) model {
	t.Helper()
	m := press(changeSpecModel(t, archived), tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.level != levelChangeSpec {
		t.Fatalf("the deltas did not open, level = %d (%q)", m.level, m.statusMsg)
	}
	return m
}

// --- 5.x the level ---

func TestEnterOnTheSpecsSubTabOpensTheDeltas(t *testing.T) {
	m := openedChangeSpecs(t, false)
	if m.specTree.name != "a-change" {
		t.Errorf("opened %q, want the change", m.specTree.name)
	}
	if m.focus != focusListPane {
		t.Errorf("the outline holds the keyboard on open, got focus %d", m.focus)
	}
	if m.specNode != 0 {
		t.Errorf("the cursor starts on the first node, got %d", m.specNode)
	}
	if m.cardView != viewDiff {
		t.Errorf("the card opens on the difference, got view %d", m.cardView)
	}
}

func TestEnterOnAnArtifactSubTabDoesNothing(t *testing.T) {
	m := changeSpecModel(t, false)
	for tab := 0; tab < m.changeSpecsTabIndex(); tab++ {
		m.changeArtifactTab = tab
		after := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
		if after.level != levelChange {
			t.Errorf("sub-tab %d: level = %d, want the change to stay open", tab, after.level)
		}
	}
}

func TestEscReturnsToTheSpecsSubTab(t *testing.T) {
	m := openedChangeSpecs(t, false)
	back := press(m, tea.KeyPressMsg{Code: tea.KeyEscape})
	if back.level != levelChange {
		t.Fatalf("level = %d, want levelChange", back.level)
	}
	if back.changeArtifactTab != back.changeSpecsTabIndex() {
		t.Errorf("returned to sub-tab %d, want the specs sub-tab %d",
			back.changeArtifactTab, back.changeSpecsTabIndex())
	}
	if len(back.specTree.nodes) != 0 {
		t.Error("the outline should be dropped on the way out")
	}
}

func TestEnterGoesNoDeeperThanTheDeltas(t *testing.T) {
	m := openedChangeSpecs(t, false)
	after := press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if after.level != levelChangeSpec {
		t.Errorf("level = %d, want it unchanged", after.level)
	}
}

// --- 5.3 and 5.4 the arrow keys ---

func nodeAt(m model, title string) int {
	for i, n := range m.specTree.nodes {
		if n.title == title {
			return i
		}
	}
	return -1
}

func TestArrowsMoveTheChooserAndStopAtBothEnds(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "Still this")
	m.rememberSpecNode()
	if n, _ := m.selectedNode(); !n.comparable() {
		t.Fatal("this node should have an original to compare against")
	}

	for _, want := range []int{viewOld, viewNew, viewNew} {
		m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
		if m.cardView != want {
			t.Fatalf("right: got view %d, want %d", m.cardView, want)
		}
	}
	for _, want := range []int{viewOld, viewDiff, viewDiff} {
		m = press(m, tea.KeyPressMsg{Code: tea.KeyLeft})
		if m.cardView != want {
			t.Fatalf("left: got view %d, want %d", m.cardView, want)
		}
	}
}

func TestArrowsDoNothingOnANodeWithNothingToCompare(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "Brand new") // added, so no original
	m.rememberSpecNode()

	after := press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	if after.cardView != viewDiff {
		t.Errorf("view moved to %d on a node with no original", after.cardView)
	}
	if after.cardViewRowShown() {
		t.Error("the chooser should not be drawn for a node with no original")
	}
}

func TestArrowsChangeNeitherTheTabBarNorTheSubTab(t *testing.T) {
	m := openedChangeSpecs(t, false)
	tab, sub, lvl := m.detailTab, m.changeArtifactTab, m.level

	for i := 0; i < 6; i++ {
		m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
		m = press(m, tea.KeyPressMsg{Code: tea.KeyLeft})
	}
	if m.detailTab != tab {
		t.Errorf("the project tab bar moved to %d", m.detailTab)
	}
	if m.changeArtifactTab != sub {
		t.Errorf("the change's sub-tab moved to %d", m.changeArtifactTab)
	}
	if m.level != lvl {
		t.Errorf("the level moved to %d", m.level)
	}
}

func TestMovingToAnotherNodeOpensOnTheDifference(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "Still this")
	m.rememberSpecNode()
	m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	if m.cardView != viewOld {
		t.Fatalf("setup: view = %d", m.cardView)
	}

	m.specNode = nodeAt(m, "Left alone")
	m.rememberSpecNode()
	if m.cardView != viewDiff {
		t.Errorf("a new node opens on the difference, got view %d", m.cardView)
	}
}

// --- 3.2 to 3.4 the outline marks ---

func outlineText(m model) string {
	outlineOuter, _ := specDetailSplit(m.panelContentWidth())
	return ansi.Strip(renderSpecOutline(m.specTree, m.selectedSpecNode(),
		outlineOuter-boxChrome, 60, true))
}

func TestTheOutlineMarksWhatTheChangeDoes(t *testing.T) {
	m := openedChangeSpecs(t, false)
	out := outlineText(m)

	for _, c := range []struct{ want, why string }{
		{"~ An old thing", "a modified requirement"},
		{"+ A new thing", "an added requirement"},
		{"- A gone thing", "a removed requirement"},
		{"+ Brand new", "a scenario the change introduces"},
		{"~ Left alone", "a scenario the change edits"},
	} {
		if !strings.Contains(out, c.want) {
			t.Errorf("the outline should mark %s as %q:\n%s", c.why, c.want, out)
		}
	}
	// An unchanged scenario carries no sigil, only the gutter its neighbours use.
	if !strings.Contains(out, "    Still this") {
		t.Errorf("an unchanged scenario keeps the gutter without a sigil:\n%s", out)
	}
}

func TestAMainSpecOutlineIsUnmarked(t *testing.T) {
	// The gutter belongs to a change's outline. A live spec is drawn exactly as
	// it was before there were marks at all.
	m := press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
	out := ansi.Strip(renderSpecOutline(m.specTree, 0, 40, 40, true))
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "+ ") ||
			strings.HasPrefix(strings.TrimSpace(line), "~ ") {
			t.Errorf("a main spec outline carries a mark: %q", line)
		}
	}
}

func TestAWrappedLabelCarriesItsMarkOnlyOnce(t *testing.T) {
	m := openedChangeSpecs(t, false)
	out := ansi.Strip(renderSpecOutline(m.specTree, 0, 18, 60, true))
	marks := strings.Count(out, "~ ") + strings.Count(out, "+ ") + strings.Count(out, "- ")
	reqs, scens := m.specTree.counts()
	if marks > reqs+scens {
		t.Errorf("a wrapped label repeats its mark: %d marks for %d nodes:\n%s",
			marks, reqs+scens, out)
	}
}

// --- 4.x the card ---

func cardText2(t *testing.T, m model) string {
	t.Helper()
	m.syncDocument()
	width, _ := m.docRegion()
	n, _ := m.selectedNode()
	return ansi.Strip(renderChangeCard(n, m.cardView, width))
}

func TestACapabilityCardNamesWhatTheChangeDoesToIt(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = 0
	m.rememberSpecNode()
	if n, _ := m.selectedNode(); n.kind != nodeCapability {
		t.Fatal("the outline should open on a capability")
	}
	out := cardText2(t, m)
	if !strings.Contains(out, "a-capability") || !strings.Contains(out, "modifies") {
		t.Errorf("the capability card should name it and what happens to it:\n%s", out)
	}
}

func TestARemovalsCardKeepsItsReasonAndMigration(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "A gone thing")
	m.rememberSpecNode()
	out := cardText2(t, m)
	for _, want := range []string{"Reason", "replaced by the new thing", "Migration"} {
		if !strings.Contains(out, want) {
			t.Errorf("the removal's card loses %q:\n%s", want, out)
		}
	}
}

func TestTheDifferenceMarksOnlyWhatChanged(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "Left alone")
	m.rememberSpecNode()
	m.syncDocument()
	width, _ := m.docRegion()
	n, _ := m.selectedNode()

	raw := renderChangeCard(n, viewDiff, width)
	plain := ansi.Strip(raw)

	// Interleaved at word level, which is the point: the clause reads as one
	// sentence with the dropped word beside the word that replaced it, rather
	// than as two paragraphs the reader has to compare by eye.
	if !strings.Contains(plain, "nothing something SHALL happen after all") {
		t.Errorf("the difference should interleave both sides:\n%s", plain)
	}
	// The words neither side touched keep their own styling, and the words
	// that differ carry the marks.
	if strings.Count(raw, "\x1b") == 0 {
		t.Error("the difference should style what changed")
	}
	if !strings.Contains(raw, removedWordStyle.Render("nothing")) {
		t.Error("the dropped word should be struck through")
	}
	if !strings.Contains(raw, addedWordStyle.Render("something")) {
		t.Error("the word that replaced it should be marked as added")
	}
	if ansi.Strip(renderChangeCard(n, viewOld, width)) == plain {
		t.Error("the original should differ from the difference")
	}
}

func TestTheOldAndNewViewsShowEachSideAlone(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "Left alone")
	m.rememberSpecNode()
	m.syncDocument()
	width, _ := m.docRegion()
	n, _ := m.selectedNode()

	old := ansi.Strip(renderChangeCard(n, viewOld, width))
	nw := ansi.Strip(renderChangeCard(n, viewNew, width))
	if !strings.Contains(old, "nothing SHALL happen") || strings.Contains(old, "after all") {
		t.Errorf("the old view should be the original alone:\n%s", old)
	}
	if !strings.Contains(nw, "after all") || strings.Contains(nw, "nothing SHALL happen") {
		t.Errorf("the new view should be what the change proposes alone:\n%s", nw)
	}
}

func TestAnArchivedChangeOffersNoChooser(t *testing.T) {
	m := openedChangeSpecs(t, true)
	for i := range m.specTree.nodes {
		m.specNode = i
		m.rememberSpecNode()
		if m.cardViewRowShown() {
			t.Fatalf("node %q offers a chooser for an archived change",
				m.specTree.nodes[i].title)
		}
	}
}

// --- 6.x opening the file ---

func TestEOpensTheDeltaTheNodeBelongsTo(t *testing.T) {
	m := openedChangeSpecs(t, false)
	for _, c := range []struct{ node, want string }{
		{"An old thing", "a-capability"},
		{"A new thing", "b-capability"},
	} {
		at := nodeAt(m, c.node)
		if at < 0 {
			t.Fatalf("node %q is missing", c.node)
		}
		m.specNode = at
		m.rememberSpecNode()
		m.syncDocument()
		want := "/openspec/changes/a-change/specs/" + c.want + "/spec.md"
		if !strings.HasSuffix(m.docPath, want) {
			t.Errorf("node %q: E would open %q, want it to end %q", c.node, m.docPath, want)
		}
	}
}

func TestTheSpecsSubTabStillOffersNoFile(t *testing.T) {
	m := changeSpecModel(t, false)
	m.syncDocument()
	if m.docPath != "" {
		t.Errorf("the specs sub-tab shows several files and should offer none, got %q", m.docPath)
	}
}

// --- 5.5 the level's arm in every switch that branches on one ---

func TestTheLevelHasItsArmEverywhere(t *testing.T) {
	m := openedChangeSpecs(t, false)

	if !m.splitTab() {
		t.Error("splitTab: the view is an outline beside a card")
	}
	if w, h := m.docRegion(); w < 10 || h < 1 {
		t.Errorf("docRegion: got %dx%d", w, h)
	}
	if m.listPage() < 1 {
		t.Errorf("listPage: got %d", m.listPage())
	}

	// The vertical keys belong to the card only while the card holds them.
	m.focus = focusListPane
	if m.docActive() {
		t.Error("docActive: the outline holds the keyboard")
	}
	m.focus = focusContentPane
	if !m.docActive() {
		t.Error("docActive: the card holds the keyboard")
	}

	// The list keys move the outline only while the outline holds them.
	m.focus = focusListPane
	before := m.specNode
	m.moveListCursor(1)
	if m.specNode == before {
		t.Error("moveListCursor: the outline did not move")
	}
	m.gotoListEnd(true)
	if m.specNode != len(m.specTree.nodes)-1 {
		t.Errorf("gotoListEnd: got node %d, want the last", m.specNode)
	}
	m.gotoListEnd(false)
	if m.specNode != 0 {
		t.Errorf("gotoListEnd: got node %d, want the first", m.specNode)
	}

	nav := ansi.Strip(m.renderNavBar())
	if !strings.Contains(nav, "back to change") {
		t.Errorf("the nav bar should offer the way out:\n%s", nav)
	}
}

func TestTheChooserIsAdvertisedOnlyWhereItActs(t *testing.T) {
	m := openedChangeSpecs(t, false)

	m.specNode = nodeAt(m, "Still this")
	m.rememberSpecNode()
	if !strings.Contains(ansi.Strip(m.renderNavBar()), "diff/old/new") {
		t.Error("a node with an original should advertise the chooser")
	}

	m.specNode = nodeAt(m, "Brand new")
	m.rememberSpecNode()
	if strings.Contains(ansi.Strip(m.renderNavBar()), "diff/old/new") {
		t.Error("a node with no original should not advertise it")
	}
}

// --- 5.6 the cursor survives a rebuild ---

func TestTheCursorStaysOnItsNodeAcrossARebuild(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.specNode = nodeAt(m, "Brand new")
	m.rememberSpecNode()
	was := m.specNodePath

	// A requirement added above the cursor, which is what an edit during a
	// read looks like.
	info := m.projects["/p"].Info
	ci := info.Changes[0]
	ci.SpecContents = map[string]string{
		"a-capability": "## ADDED Requirements\n\n### Requirement: Inserted first\n" +
			"It SHALL.\n\n#### Scenario: S\n- **WHEN** a\n- **THEN** b\n\n" + modifyingDelta,
		"b-capability": threeOps,
	}
	info.Changes = []scanner.ChangeInfo{ci}
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: info}}

	m.reparseOpenChangeSpecs()
	if m.specNodePath != was {
		t.Errorf("the cursor moved to %q, want it to stay on %q", m.specNodePath, was)
	}
	if n, _ := m.selectedNode(); n.title != "Brand new" {
		t.Errorf("the cursor landed on %q", n.title)
	}
}

func TestADeltaEditedIntoNonsenseKeepsTheOutlineOpen(t *testing.T) {
	m := openedChangeSpecs(t, false)
	nodes := len(m.specTree.nodes)

	info := m.projects["/p"].Info
	ci := info.Changes[0]
	ci.SpecContents = map[string]string{"a-capability": "nothing here\n", "b-capability": threeOps}
	info.Changes = []scanner.ChangeInfo{ci}
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: info}}

	m.reparseOpenChangeSpecs()
	if len(m.specTree.nodes) != nodes {
		t.Error("the outline should be kept rather than emptied under a reader")
	}
}

// --- 7.1 the view behaves as spec-detail-view requires ---

func TestTheFrameIsExactAtEveryWidth(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {100, 30}, {140, 50}} {
		m := openedChangeSpecs(t, false)
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

func TestTheCardScrollsAndReportsItsPosition(t *testing.T) {
	m := openedChangeSpecs(t, false)
	m.height = 12
	m.recalcLayout()
	m.focus = focusContentPane
	m.specNode = nodeAt(m, "A gone thing")
	m.rememberSpecNode()
	m.syncDocument()

	if !m.docActive() {
		t.Fatal("the card holds the keyboard")
	}
	m.focus = focusListPane
	if m.docScrollPercent() != -1 {
		t.Error("no position is reported while the outline holds the keyboard")
	}
}
