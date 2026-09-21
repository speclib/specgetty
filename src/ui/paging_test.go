package ui

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// pagingModel opens this project, which has documents long enough to scroll on
// both split tabs.
func pagingModel(t *testing.T, tab int) model {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	key, st, err := scanner.ScanResolved(root)
	if err != nil || key == "" {
		t.Skip("not this project")
	}
	m := model{width: 100, height: 24,
		repoPaths:    []string{key},
		projects:     scanner.ProjectMap{key: st},
		displayNames: []string{"specgetty"},
		detailTab:    tabChanges,
		focus:        focusDetail,
		fields:       append([]string(nil), defaultFields...),
	}
	m.recalcLayout()
	return press(m, tea.KeyPressMsg{Code: rune('1' + tab), Text: string(rune('1' + tab))})
}

// --- paging follows the keyboard ---

func TestPagingTheContentHalfScrollsTheDocument(t *testing.T) {
	for _, tc := range []struct {
		name string
		tab  int
	}{
		{"properties", tabProperties},
		{"specs", tabSpecs},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := press(pagingModel(t, tc.tab), tea.KeyPressMsg{Code: tea.KeyTab})
			if m.focus != focusContentPane {
				t.Fatalf("expected the content to hold the keyboard, got focus %d", m.focus)
			}
			if m.docViewport.TotalLineCount() <= m.docViewport.Height() {
				t.Skip("the document fits; nothing to scroll")
			}

			for _, k := range []tea.KeyPressMsg{
				{Code: tea.KeyPgDown},
				{Code: 'f', Mod: tea.ModCtrl},
				{Code: 'd', Mod: tea.ModCtrl},
			} {
				moved := press(m, k)
				if moved.docViewport.YOffset() == 0 {
					t.Errorf("%v did not scroll the document", k)
				}
			}

			end := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
			if end.docViewport.YOffset() == 0 {
				t.Error("G did not reach the end")
			}
			top := press(press(end, tea.KeyPressMsg{Code: 'g', Text: "g"}),
				tea.KeyPressMsg{Code: 'g', Text: "g"})
			if top.docViewport.YOffset() != 0 {
				t.Errorf("gg left the document at %d", top.docViewport.YOffset())
			}

			atEnd := end.docViewport.YOffset()
			if up := press(end, tea.KeyPressMsg{Code: tea.KeyPgUp}); up.docViewport.YOffset() >= atEnd {
				t.Error("pgup did not scroll back")
			}
			if halfUp := press(end, tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}); halfUp.docViewport.YOffset() >= atEnd {
				t.Error("ctrl+u did not scroll back")
			}
		})
	}
}

func TestPagingTheListHalfMovesTheListNotTheDocument(t *testing.T) {
	for _, tc := range []struct {
		name string
		tab  int
		of   func(model) int
	}{
		{"properties", tabProperties, func(m model) int { return m.propSection }},
		{"specs", tabSpecs, func(m model) int { return m.specCursor }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := pagingModel(t, tc.tab)
			if m.focus != focusListPane {
				t.Fatalf("expected the list to hold the keyboard, got focus %d", m.focus)
			}

			end := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
			if tc.of(end) == tc.of(m) {
				t.Error("G must move the list selection")
			}
			if end.docViewport.YOffset() != 0 {
				t.Errorf("and must not scroll the document, offset %d", end.docViewport.YOffset())
			}

			back := press(press(end, tea.KeyPressMsg{Code: 'g', Text: "g"}),
				tea.KeyPressMsg{Code: 'g', Text: "g"})
			if tc.of(back) != 0 {
				t.Errorf("gg must reach the first row, got %d", tc.of(back))
			}

			paged := press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
			if paged.docViewport.YOffset() != 0 {
				t.Error("pgdown must not scroll the document while the list holds the keyboard")
			}
		})
	}
}

// --- 2.3 the line keys are untouched ---

func TestLineKeysStillWalkTheListRatherThanScroll(t *testing.T) {
	m := pagingModel(t, tabProperties)
	before := m.propSection

	down := press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if down.propSection == before {
		t.Error("j must move the selected row while the list holds the keyboard")
	}
	if down.docViewport.YOffset() != 0 {
		t.Errorf("j must not scroll the document, offset %d", down.docViewport.YOffset())
	}

	s := pagingModel(t, tabSpecs)
	sdown := press(s, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if sdown.specCursor != 1 {
		t.Errorf("j must move the spec cursor, got %d", sdown.specCursor)
	}
	if sdown.docViewport.YOffset() != 0 {
		t.Error("j must not scroll the spec content")
	}
}

// --- 2.5 an open change is unaffected ---

func TestAnOpenChangeStillPages(t *testing.T) {
	m := pagingModel(t, tabChanges)
	m.changeCursor = 1
	m.rememberSelection()
	m = press(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.level != levelChange {
		t.Fatal("did not descend into a change")
	}
	if m.docViewport.TotalLineCount() <= m.docViewport.Height() {
		t.Skip("the artifact fits")
	}
	paged := press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	if paged.docViewport.YOffset() == 0 {
		t.Error("an open change was never a split and must still page")
	}
}

// --- 2.6 the reading position is unchanged ---

func TestThePositionIsStillReportedOnlyWhileTheContentIsFocused(t *testing.T) {
	m := pagingModel(t, tabSpecs)
	if pct := m.docScrollPercent(); pct != -1 {
		t.Errorf("the list holds the keyboard but a position of %d%% was reported", pct)
	}
	focused := press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if pct := focused.docScrollPercent(); pct < 0 {
		t.Error("the content holding the keyboard does report one")
	}
}

func TestTheDocumentOwnsTheKeysOnlyWhileItHoldsTheKeyboard(t *testing.T) {
	m := pagingModel(t, tabProperties)
	if m.docActive() {
		t.Error("the list holds the keyboard, so the document does not own the keys")
	}
	onContent := press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if !onContent.docActive() {
		t.Error("it does once the content holds the keyboard")
	}
}
