package ui

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/specgetty/src/scanner"
)

func testInfo() scanner.ProjectInfo {
	return scanner.ProjectInfo{
		Changes: []scanner.ChangeInfo{
			{Name: "alpha", TasksTotal: 4, TasksDone: 1},
			{Name: "beta"},
		},
		ArchivedChanges: []scanner.ChangeInfo{
			{Name: "gamma"},
		},
	}
}

func TestBuildRowsModes(t *testing.T) {
	info := testInfo()

	tests := []struct {
		mode int
		want []string
	}{
		{modeOpen, []string{"alpha", "beta"}},
		{modeArchived, []string{"gamma"}},
		{modeBoth, []string{"alpha", "beta", "gamma"}},
	}
	for _, tt := range tests {
		got := names(buildRows(info, tt.mode))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("mode %d = %v, want %v", tt.mode, got, tt.want)
		}
	}
}

func TestBuildRowsTagsArchived(t *testing.T) {
	rows := buildRows(testInfo(), modeBoth)
	if rows[0].archived || rows[1].archived {
		t.Error("active changes must not be tagged archived")
	}
	if !rows[2].archived {
		t.Error("archived change must be tagged archived")
	}
}

func TestRowKeyDistinguishesArchivedFromActive(t *testing.T) {
	// An archived change keeps the name it had while active, so the name alone
	// is not unique in modeBoth.
	active := changeRow{ci: scanner.ChangeInfo{Name: "same"}}
	archived := changeRow{ci: scanner.ChangeInfo{Name: "same"}, archived: true}
	if active.key() == archived.key() {
		t.Errorf("keys collide: both are %q", active.key())
	}
}

func TestIndexOfKey(t *testing.T) {
	rows := buildRows(testInfo(), modeBoth)
	if got := indexOfKey(rows, "open/beta"); got != 1 {
		t.Errorf("indexOfKey(open/beta) = %d, want 1", got)
	}
	if got := indexOfKey(rows, "archived/gamma"); got != 2 {
		t.Errorf("indexOfKey(archived/gamma) = %d, want 2", got)
	}
	if got := indexOfKey(rows, "open/missing"); got != -1 {
		t.Errorf("indexOfKey(open/missing) = %d, want -1", got)
	}
}

func TestEmptyListMessageDiffersByMode(t *testing.T) {
	seen := map[string]bool{}
	for _, mode := range []int{modeOpen, modeArchived, modeBoth} {
		msg := emptyListMessage(mode)
		if msg == "" {
			t.Errorf("mode %d has no message", mode)
		}
		if seen[msg] {
			t.Errorf("mode %d repeats an earlier message %q", mode, msg)
		}
		seen[msg] = true
	}
}

func TestArtifactTabNames(t *testing.T) {
	r := changeRow{ci: scanner.ChangeInfo{
		ArtifactFiles: []string{"design.md", "proposal.md", "tasks.md"},
		SpecNames:     []string{"a"},
	}}
	want := []string{"design", "proposal", "tasks", "specs"}
	if got := r.artifactTabNames(); !reflect.DeepEqual(got, want) {
		t.Errorf("artifactTabNames = %v, want %v", got, want)
	}

	// No specs directory means no specs sub-tab.
	r.ci.SpecNames = nil
	if got := r.artifactTabNames(); len(got) != 3 {
		t.Errorf("artifactTabNames without specs = %v, want 3 entries", got)
	}
}

// --- model-level behaviour ---

func makeListModel() model {
	m := newModel(nil, true, "test")
	m.scanning = false // newModel starts scanning, which swallows every key
	m.width, m.height = 100, 40
	m.repoPaths = []string{"/p"}
	m.displayNames = []string{"p"}
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: testInfo()}}
	m.level = levelProject
	m.activeView = viewDetail
	m.detailTab = tabChanges
	return m
}

func TestCurrentRowsAppliesModeAndQuery(t *testing.T) {
	m := makeListModel()
	if got := len(m.currentRows()); got != 2 {
		t.Errorf("modeOpen gave %d rows, want 2", got)
	}

	m.listMode = modeBoth
	if got := len(m.currentRows()); got != 3 {
		t.Errorf("modeBoth gave %d rows, want 3", got)
	}

	m.searchInput.SetValue("gamma")
	got := names(m.currentRows())
	if !reflect.DeepEqual(got, []string{"gamma"}) {
		t.Errorf("filtered rows = %v, want [gamma]", got)
	}
}

func TestSyncCursorFollowsSurvivingSelection(t *testing.T) {
	m := makeListModel()
	m.listMode = modeBoth
	m.changeCursor = 2 // gamma
	m.rememberSelection()

	// Narrowing the list moves gamma to index 0; the cursor must follow it
	// rather than stay on index 2 or clamp to a different change.
	m.searchInput.SetValue("gamma")
	m.syncCursor()

	if m.changeCursor != 0 {
		t.Errorf("cursor = %d, want 0", m.changeCursor)
	}
	r, ok := m.selectedRow()
	if !ok || r.ci.Name != "gamma" {
		t.Errorf("selected %v, want gamma", r.ci.Name)
	}
}

func TestSyncCursorClampsWhenSelectionFilteredOut(t *testing.T) {
	m := makeListModel()
	m.changeCursor = 1 // beta
	m.rememberSelection()

	m.searchInput.SetValue("alpha")
	m.syncCursor()

	rows := m.currentRows()
	if m.changeCursor < 0 || m.changeCursor >= len(rows) {
		t.Fatalf("cursor %d out of range for %d rows", m.changeCursor, len(rows))
	}
	if rows[m.changeCursor].ci.Name != "alpha" {
		t.Errorf("selected %q, want alpha", rows[m.changeCursor].ci.Name)
	}
}

func TestSyncCursorOnEmptyList(t *testing.T) {
	m := makeListModel()
	m.changeCursor = 1
	m.searchInput.SetValue("nothing-matches-this")
	m.syncCursor()
	if m.changeCursor != 0 {
		t.Errorf("cursor = %d, want 0 on an empty list", m.changeCursor)
	}
	if _, ok := m.selectedRow(); ok {
		t.Error("selectedRow should report no row when the list is empty")
	}
}

func TestEnterDescendsAndEscAscends(t *testing.T) {
	m := makeListModel()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	um := updated.(model)
	if um.level != levelChange {
		t.Fatalf("after enter level = %d, want levelChange", um.level)
	}

	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyEsc})
	um = updated.(model)
	if um.level != levelProject {
		t.Errorf("after esc level = %d, want levelProject", um.level)
	}
}

func TestEscDoesNotRunPastTheEnds(t *testing.T) {
	m := makeListModel()
	m.level = levelProjects
	m.activeView = viewProjects
	m.fullScanDone = true

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if um := updated.(model); um.level != levelProjects {
		t.Errorf("esc at the top level moved to %d, want to stay at levelProjects", um.level)
	}
}

func TestRightArrowDoesNotSpillIntoTheTabBar(t *testing.T) {
	// The regression this change exists to fix: right at the last artifact
	// sub-tab used to jump to a different project tab and reset the position.
	m := makeListModel()
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Changes: []scanner.ChangeInfo{{
			Name:          "alpha",
			ArtifactFiles: []string{"design.md", "proposal.md"},
		}},
	}}}
	m.level = levelChange
	m.changeArtifactTab = 1 // the last sub-tab
	before := m.detailTab

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	um := updated.(model)

	if um.detailTab != before {
		t.Errorf("detailTab moved to %d, want it unchanged at %d", um.detailTab, before)
	}
	if um.changeArtifactTab != 1 {
		t.Errorf("changeArtifactTab = %d, want it to stay at 1", um.changeArtifactTab)
	}
}

func TestLeftArrowDoesNotSpillIntoTheTabBar(t *testing.T) {
	m := makeListModel()
	m.level = levelChange
	m.detailTab = tabSpecs
	m.changeArtifactTab = 0
	before := m.detailTab

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	um := updated.(model)

	if um.detailTab != before {
		t.Errorf("detailTab moved to %d, want it unchanged at %d", um.detailTab, before)
	}
}

func TestNumberKeysInertWhileChangeIsOpen(t *testing.T) {
	m := makeListModel()
	m.level = levelChange
	m.detailTab = tabChanges

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	if um := updated.(model); um.detailTab != tabChanges {
		t.Errorf("detailTab = %d, want it unchanged while a change is open", um.detailTab)
	}
}

func TestNumberKeysSwitchTabsAtTheProjectLevel(t *testing.T) {
	m := makeListModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	if um := updated.(model); um.detailTab != tabSpecs {
		t.Errorf("detailTab = %d, want tabSpecs (%d)", um.detailTab, tabSpecs)
	}
}

func TestFCyclesListMode(t *testing.T) {
	m := makeListModel()
	for _, want := range []int{modeArchived, modeBoth, modeOpen} {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		m = updated.(model)
		if m.listMode != want {
			t.Fatalf("listMode = %d, want %d", m.listMode, want)
		}
	}
}

func TestArchiveKeyIgnoredOnArchivedRow(t *testing.T) {
	m := makeListModel()
	m.listMode = modeArchived // only gamma, which is archived

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if um := updated.(model); um.archiveState != archiveIdle {
		t.Errorf("archiveState = %d, want archiveIdle: archiving an archived change is a no-op", um.archiveState)
	}
}

func TestExportKeyWorksOnArchivedRow(t *testing.T) {
	m := makeListModel()
	m.listMode = modeArchived

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	um := updated.(model)
	if um.exportState != exportConfirming {
		t.Fatalf("exportState = %d, want exportConfirming", um.exportState)
	}
	if !um.exportIsArchived {
		t.Error("exportIsArchived = false, want true for an archived row")
	}
}

func TestSearchPromptCapturesKeys(t *testing.T) {
	m := makeListModel()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = updated.(model)
	if !m.searchFocused {
		t.Fatal("/ should focus the search prompt")
	}

	// 'a' is the archive action, but inside the prompt it is just a character.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)
	if m.archiveState != archiveIdle {
		t.Error("a while typing must not trigger archive")
	}
	if m.searchInput.Value() != "a" {
		t.Errorf("query = %q, want a", m.searchInput.Value())
	}
}

func TestEscapeClearsTheFilter(t *testing.T) {
	m := makeListModel()
	m.searchFocused = true
	m.searchInput.Focus()
	m.searchInput.SetValue("alpha")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	um := updated.(model)
	if um.searchFocused {
		t.Error("esc should unfocus the prompt")
	}
	if um.searchInput.Value() != "" {
		t.Errorf("query = %q, want cleared", um.searchInput.Value())
	}
	if len(um.currentRows()) != 2 {
		t.Errorf("full list not restored: %v", names(um.currentRows()))
	}
}

func TestEnterFromPromptOpensTheChange(t *testing.T) {
	m := makeListModel()
	m.searchFocused = true
	m.searchInput.Focus()
	m.searchInput.SetValue("beta")
	m.syncCursor()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	um := updated.(model)
	if um.level != levelChange {
		t.Fatalf("level = %d, want levelChange", um.level)
	}
	// The filter stays applied to the list behind the open change.
	if um.searchInput.Value() != "beta" {
		t.Errorf("query = %q, want it preserved", um.searchInput.Value())
	}
	r, ok := um.selectedRow()
	if !ok || r.ci.Name != "beta" {
		t.Errorf("opened %v, want beta", r.ci.Name)
	}
}

func TestArrowsNavigateWhileTyping(t *testing.T) {
	m := makeListModel()
	m.searchFocused = true
	m.searchInput.Focus()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	um := updated.(model)
	if um.changeCursor != 1 {
		t.Errorf("cursor = %d, want 1", um.changeCursor)
	}
	if um.searchInput.Value() != "" {
		t.Errorf("query = %q, want the arrow not to be typed", um.searchInput.Value())
	}
}

func TestProjectSwitchClearsTheFilter(t *testing.T) {
	m := makeListModel()
	m.level = levelProjects
	m.activeView = viewProjects
	m.repoPaths = []string{"/p", "/q"}
	m.displayNames = []string{"p", "q"}
	m.projects["/q"] = scanner.ProjectStatus{Info: testInfo()}
	m.searchInput.SetValue("alpha")
	m.listMode = modeBoth

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	um := updated.(model)

	if um.searchInput.Value() != "" {
		t.Errorf("query = %q, want cleared on project switch", um.searchInput.Value())
	}
	if um.listMode != modeOpen {
		t.Errorf("listMode = %d, want modeOpen on project switch", um.listMode)
	}
}

func TestRescanPreservesTheFilter(t *testing.T) {
	// The watcher fires a rescan on every file save. Losing the filter there
	// would be maddening, so scanMsg must not reset it.
	m := makeListModel()
	m.searchInput.SetValue("alpha")
	m.rememberSelection()

	updated, _ := m.Update(scanMsg{projects: scanner.ProjectMap{
		"/p": scanner.ProjectStatus{Info: testInfo()},
	}})
	um := updated.(model)

	if um.searchInput.Value() != "alpha" {
		t.Errorf("query = %q, want it preserved across a rescan", um.searchInput.Value())
	}
}

func TestRenderChangesTabShowsNoMatchMessage(t *testing.T) {
	m := makeListModel()
	m.searchInput.SetValue("zzzz")
	got := m.renderChangesTab(80, 10)
	if !strings.Contains(got, "zzzz") {
		t.Errorf("no-match message should echo the query, got:\n%s", got)
	}
	if !strings.Contains(got, "No changes match") {
		t.Errorf("expected a no-match message, got:\n%s", got)
	}
}

func TestRenderChangesTabShowsModeSpecificEmptyMessage(t *testing.T) {
	m := makeListModel()
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{}}}
	m.listMode = modeArchived
	got := m.renderChangesTab(80, 10)
	if !strings.Contains(got, "No archived changes") {
		t.Errorf("expected the archived empty message, got:\n%s", got)
	}
}

func TestRenderChangeTableShowsHeadersAndRows(t *testing.T) {
	rows := buildRows(testInfo(), modeOpen)
	got := renderChangeTable(rows, defaultFields, 0, 80, 10)
	for _, want := range []string{"name", "tasks", "specs", "alpha", "beta", "1/4"} {
		if !strings.Contains(got, want) {
			t.Errorf("table missing %q, got:\n%s", want, got)
		}
	}
}

func TestRenderChangeTableShowsMatchHint(t *testing.T) {
	rows := []changeRow{{
		ci:           scanner.ChangeInfo{Name: "alpha"},
		matchedFiles: []string{"proposal", "tasks"},
	}}
	got := renderChangeTable(rows, defaultFields, -1, 80, 10)
	if !strings.Contains(got, "proposal, tasks") {
		t.Errorf("table should explain a body match, got:\n%s", got)
	}
}

func TestRenderChangeDetailShowsNameAndTabs(t *testing.T) {
	r := changeRow{ci: scanner.ChangeInfo{
		Name:             "alpha",
		ArtifactFiles:    []string{"proposal.md", "tasks.md"},
		ArtifactContents: map[string]string{"proposal.md": "why this", "tasks.md": "- [ ] a"},
		TasksTotal:       1,
	}}
	got := renderChangeDetail(r, 0, 80, 20)
	for _, want := range []string{"alpha", "open", "proposal", "tasks", "why this"} {
		if !strings.Contains(got, want) {
			t.Errorf("detail missing %q, got:\n%s", want, got)
		}
	}

	r.archived = true
	if got := renderChangeDetail(r, 0, 80, 20); !strings.Contains(got, "archived") {
		t.Errorf("archived change should say so, got:\n%s", got)
	}
}
