package ui

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

func testProjectRows() []projectRow {
	return []projectRow{
		{
			path: "/home/p/specgetty", display: "specgetty",
			files: []scanner.FileEntry{{Path: "specs/zoom-mode/spec.md"}},
			info: scanner.ProjectInfo{
				SpecCount:    3,
				SpecContents: map[string]string{"zoom-mode": "the user SHALL enter zoom mode"},
				Changes:      []scanner.ChangeInfo{{Name: "alpha"}},
				TasksTotal:   10, TasksDone: 4,
			},
		},
		{
			path: "/home/p/mipnix", display: "mipnix",
			files: []scanner.FileEntry{{Path: "specs/flake/spec.md"}},
			info: scanner.ProjectInfo{
				SpecCount:       7,
				SpecContents:    map[string]string{"flake": "nix flake inputs"},
				ArchivedChanges: []scanner.ChangeInfo{{Name: "old"}},
			},
		},
	}
}

func makePickerModel() model {
	m := makeListModel()
	m.pickerAll = testProjectRows()
	m.pickerLoaded = true
	return m
}

func projectNames(rows []filtered[projectRow]) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.row.display
	}
	return out
}

func TestPickerOpensFromBothLevels(t *testing.T) {
	for _, level := range []int{levelProject, levelChange} {
		m := makePickerModel()
		m.level = level

		updated, _ := m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
		um := updated.(model)
		if !um.pickerOpen {
			t.Errorf("at level %d, p did not open the picker", level)
		}
		if um.level != level {
			t.Errorf("opening the picker changed the level from %d to %d", level, um.level)
		}
	}
}

func TestPickerDismissLeavesCurrentProjectAlone(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true
	before := m.repoPaths[0]

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	um := updated.(model)
	if um.pickerOpen {
		t.Error("esc should close the picker")
	}
	if um.repoPaths[0] != before {
		t.Errorf("current project changed to %q on dismiss, want %q", um.repoPaths[0], before)
	}
}

func TestPickerSelectionSwitchesProject(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true
	m.pickerCursor = 1 // mipnix

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)

	if um.pickerOpen {
		t.Error("selecting should close the picker")
	}
	if len(um.repoPaths) != 1 || um.repoPaths[0] != "/home/p/mipnix" {
		t.Errorf("repoPaths = %v, want the selected project", um.repoPaths)
	}
	if um.displayNames[0] != "mipnix" {
		t.Errorf("displayNames = %v, want mipnix", um.displayNames)
	}
}

func TestPickerSelectionFromAnOpenChangeLandsAtTheProjectView(t *testing.T) {
	// The change that was open does not exist in the newly selected project,
	// so staying at the change level would be meaningless.
	m := makePickerModel()
	m.level = levelChange
	m.pickerOpen = true
	m.pickerCursor = 1

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)

	if um.level != levelProject {
		t.Errorf("level = %d, want levelProject after switching project", um.level)
	}
	if um.detailTab != tabChanges {
		t.Errorf("detailTab = %d, want tabChanges", um.detailTab)
	}
}

func TestPickerSelectionResetsChangeListState(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true
	m.pickerCursor = 1
	m.searchInput.SetValue("alpha")

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)

	if um.searchInput.Value() != "" {
		t.Errorf("change filter = %q, want cleared when the project changes", um.searchInput.Value())
	}
	if um.changeCursor != 0 {
		t.Errorf("changeCursor = %d, want the top of the new project's list", um.changeCursor)
	}
}

func TestPickerKeysOutrankTheChangeList(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true

	// 'a' archives in the change list; while the picker is open it must not.
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	if um := updated.(model); um.archiveState != archiveIdle {
		t.Error("a reached the change list while the picker was open")
	}

	// 'd' discards in the change list; while the picker is open it must not.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'd', Text: "d"})
	if um := updated.(model); um.discardState != discardIdle {
		t.Error("d reached the change list while the picker was open")
	}
}

func TestConfirmModalOutranksThePicker(t *testing.T) {
	m := makePickerModel()
	m.archiveState = archiveConfirming
	m.archiveChangeName = "alpha"

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	um := updated.(model)
	if um.pickerOpen {
		t.Error("p opened the picker while a confirmation was waiting for an answer")
	}
}

func TestPickerFilterGrammar(t *testing.T) {
	m := makePickerModel()

	t.Run("fuzzy on name", func(t *testing.T) {
		m.pickerInput.SetValue("spg")
		if got := projectNames(m.pickerVisibleRows()); !reflect.DeepEqual(got, []string{"specgetty"}) {
			t.Errorf("got %v, want [specgetty]", got)
		}
	})

	t.Run("literal on name", func(t *testing.T) {
		m.pickerInput.SetValue("'nix")
		if got := projectNames(m.pickerVisibleRows()); !reflect.DeepEqual(got, []string{"mipnix"}) {
			t.Errorf("got %v, want [mipnix]", got)
		}
	})

	t.Run("contents", func(t *testing.T) {
		m.pickerInput.SetValue(":flake inputs")
		if got := projectNames(m.pickerVisibleRows()); !reflect.DeepEqual(got, []string{"mipnix"}) {
			t.Errorf("got %v, want [mipnix]", got)
		}
	})

	t.Run("file paths", func(t *testing.T) {
		m.pickerInput.SetValue(":zoom-mode/spec.md")
		rows := m.pickerVisibleRows()
		if got := projectNames(rows); !reflect.DeepEqual(got, []string{"specgetty"}) {
			t.Errorf("got %v, want [specgetty]", got)
		}
		found := false
		for _, label := range rows[0].matched {
			if label == "paths" {
				found = true
			}
		}
		if !found {
			t.Errorf("matched labels %v should say the hit was in paths", rows[0].matched)
		}
	})

	t.Run("smart case", func(t *testing.T) {
		m.pickerInput.SetValue("MIPNIX")
		if got := len(m.pickerVisibleRows()); got != 0 {
			t.Errorf("an uppercase query matched %d rows, want 0", got)
		}
	})
}

func TestPickerCursorFollowsTheProjectAcrossAReFilter(t *testing.T) {
	m := makePickerModel()
	m.pickerCursor = 1 // mipnix
	m.pickerRemember()

	m.pickerInput.SetValue("nix")
	m.pickerSync()

	if got, _ := m.pickerSelected(); got.display != "mipnix" {
		t.Errorf("cursor landed on %q, want mipnix", got.display)
	}
}

func TestPickerCursorClampsWhenTheProjectIsFilteredOut(t *testing.T) {
	m := makePickerModel()
	m.pickerCursor = 1
	m.pickerRemember()

	m.pickerInput.SetValue("specgetty")
	m.pickerSync()

	rows := m.pickerVisibleRows()
	if m.pickerCursor < 0 || m.pickerCursor >= len(rows) {
		t.Fatalf("cursor %d out of range for %d rows", m.pickerCursor, len(rows))
	}
	if rows[m.pickerCursor].row.display != "specgetty" {
		t.Errorf("cursor landed on %q, want specgetty", rows[m.pickerCursor].row.display)
	}
}

func TestPickerSearchPromptCapturesKeys(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true

	updated, _ := m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
	m = updated.(model)
	if !m.pickerFocused {
		t.Fatal("/ should focus the picker filter")
	}

	// 'r' refreshes the picker; inside the prompt it is just a character.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	m = updated.(model)
	if m.pickerLoading {
		t.Error("r while typing must not start a refresh")
	}
	if m.pickerInput.Value() != "r" {
		t.Errorf("query = %q, want r", m.pickerInput.Value())
	}
}

func TestPickerEnterFromPromptSelects(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true
	m.pickerFocused = true
	m.pickerInput.SetValue("nix")
	m.pickerSync()

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)
	if um.pickerOpen {
		t.Error("enter should close the picker")
	}
	if um.repoPaths[0] != "/home/p/mipnix" {
		t.Errorf("opened %q, want mipnix", um.repoPaths[0])
	}
}

func TestPickerRefreshStartsLoading(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	if um := updated.(model); !um.pickerLoading {
		t.Error("r should start a refresh")
	}
	if cmd == nil {
		t.Error("r should issue a command to do the walk")
	}
}

func TestPickerLoadedMessagePopulatesRows(t *testing.T) {
	m := makeListModel()
	m.pickerLoading = true

	updated, _ := m.Update(pickerLoadedMsg{rows: testProjectRows()})
	um := updated.(model)
	if um.pickerLoading {
		t.Error("loading should finish")
	}
	if len(um.pickerAll) != 2 {
		t.Errorf("got %d rows, want 2", len(um.pickerAll))
	}
	if !um.pickerLoaded {
		t.Error("pickerLoaded should be set so a later open does not rescan")
	}
}

func TestPickerLoadErrorIsReported(t *testing.T) {
	m := makeListModel()
	m.pickerLoading = true

	updated, _ := m.Update(pickerLoadedMsg{err: errTest})
	um := updated.(model)
	if um.pickerErr == "" {
		t.Error("a scan failure should be reported, not swallowed")
	}
	if um.pickerLoading {
		t.Error("loading should finish even on failure")
	}
}

func TestBuildProjectRowsDisambiguatesNames(t *testing.T) {
	projects := scanner.ProjectMap{
		"/a/specgetty": {},
		"/b/specgetty": {},
		"/c/other":     {},
	}
	rows := buildProjectRows(projects)
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}
	// Sorted by path, so /a and /b come first and both need disambiguating.
	if rows[0].display != "specgetty (a)" || rows[1].display != "specgetty (b)" {
		t.Errorf("got %q and %q, want parent-disambiguated names",
			rows[0].display, rows[1].display)
	}
	if rows[2].display != "other" {
		t.Errorf("got %q, want a bare basename when it is unique", rows[2].display)
	}
}

func TestRenderPickerShowsProjectsAndStats(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true

	got := m.renderPicker()
	for _, want := range []string{"Projects", "specgetty", "mipnix", "4/10"} {
		if !strings.Contains(got, want) {
			t.Errorf("picker missing %q, got:\n%s", want, got)
		}
	}
}

func TestRenderPickerNoMatchEchoesQuery(t *testing.T) {
	m := makePickerModel()
	m.pickerOpen = true
	m.pickerInput.SetValue("zzzz")

	got := m.renderPicker()
	if !strings.Contains(got, "zzzz") || !strings.Contains(got, "No projects match") {
		t.Errorf("expected a no-match message echoing the query, got:\n%s", got)
	}
}

func TestRenderPickerEmptyNamesTheRefreshKey(t *testing.T) {
	m := makeListModel()
	m.pickerOpen = true
	m.pickerLoaded = true

	got := m.renderPicker()
	if !strings.Contains(got, "No OpenSpec projects found") || !strings.Contains(got, "r") {
		t.Errorf("empty picker should name the refresh key, got:\n%s", got)
	}
}

func TestEmptyProjectViewNamesThePickerKey(t *testing.T) {
	m := makeListModel()
	m.repoPaths = nil
	m.displayNames = nil

	got := m.renderDetailPanel(80, 20)
	if !strings.Contains(got, "No project selected") || !strings.Contains(got, "p") {
		t.Errorf("empty state should name the picker key, got:\n%s", got)
	}
}

func TestStartupPromptOpensThePicker(t *testing.T) {
	m := makeListModel()
	m.askOpenPicker = true

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	um := updated.(model)
	if um.askOpenPicker {
		t.Error("answering should dismiss the prompt")
	}
	if !um.pickerOpen {
		t.Error("y should open the picker")
	}
}

func TestStartupPromptDeclineLeavesAUsableView(t *testing.T) {
	m := makeListModel()
	m.askOpenPicker = true

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	um := updated.(model)
	if um.askOpenPicker {
		t.Error("answering should dismiss the prompt")
	}
	if um.pickerOpen {
		t.Error("n should not open the picker")
	}
	if cmd != nil {
		t.Error("declining should not quit the application")
	}
}

func TestGenericTableRendersBothRowTypes(t *testing.T) {
	changes := renderTable(wrap(buildGroups(testInfo())[0].rows),
		changeFieldDefs(defaultFields), 0, 80, 6)
	if !strings.Contains(changes, "alpha") {
		t.Errorf("change table missing its rows, got:\n%s", changes)
	}

	projects := renderTable(wrap(testProjectRows()), projectFields, 0, 80, 6)
	if !strings.Contains(projects, "specgetty") {
		t.Errorf("project table missing its rows, got:\n%s", projects)
	}
	if !strings.Contains(projects, "project") {
		t.Errorf("project table missing its header, got:\n%s", projects)
	}
}

var errTest = errors.New("scan failed")
