package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// pageKeys are every key this capability governs, in both directions.
var pageKeys = []struct {
	name string
	key  tea.KeyPressMsg
	down bool
}{
	{"pgdown", tea.KeyPressMsg{Code: tea.KeyPgDown}, true},
	{"ctrl+f", tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl}, true},
	{"ctrl+d", tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl}, true},
	{"pgup", tea.KeyPressMsg{Code: tea.KeyPgUp}, false},
	{"ctrl+b", tea.KeyPressMsg{Code: 'b', Mod: tea.ModCtrl}, false},
	{"ctrl+u", tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}, false},
}

// manyChanges is a project with more changes than any pane can show.
func manyChanges(t *testing.T, n int) model {
	t.Helper()
	info := scanner.ProjectInfo{Root: "/p", Origin: "/p"}
	for i := 0; i < n; i++ {
		info.Changes = append(info.Changes,
			scanner.ChangeInfo{Name: string(rune('a'+i/26)) + string(rune('a'+i%26)), DirName: "d"})
	}
	m := model{width: 100, height: 24,
		repoPaths:    []string{"/p"},
		displayNames: []string{"p"},
		detailTab:    tabChanges,
		focus:        focusDetail,
		fields:       append([]string(nil), defaultFields...),
		projects:     scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: info}},
	}
	m.recalcLayout()
	return m
}

// --- 2.1 to 2.4 the change list ---

func TestChangeListPagesInBothDirections(t *testing.T) {
	m := manyChanges(t, 60)
	page := m.listPage()
	if page < 2 {
		t.Fatalf("a page is %d rows, too small to test with", page)
	}

	down := press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	if down.changeCursor != page {
		t.Errorf("pgdown moved to %d, want a page of %d", down.changeCursor, page)
	}
	up := press(down, tea.KeyPressMsg{Code: tea.KeyPgUp})
	if up.changeCursor != 0 {
		t.Errorf("pgup moved to %d, want back to the top", up.changeCursor)
	}
	half := press(m, tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	if half.changeCursor != page/2 {
		t.Errorf("ctrl+d moved to %d, want half a page of %d", half.changeCursor, page/2)
	}
}

func TestChangeListJumpsToBothEnds(t *testing.T) {
	m := manyChanges(t, 60)
	last := len(m.currentRows()) - 1

	end := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
	if end.changeCursor != last {
		t.Errorf("G moved to %d, want the last change at %d", end.changeCursor, last)
	}
	top := press(press(end, tea.KeyPressMsg{Code: 'g', Text: "g"}),
		tea.KeyPressMsg{Code: 'g', Text: "g"})
	if top.changeCursor != 0 {
		t.Errorf("gg moved to %d, want the first change", top.changeCursor)
	}
}

func TestChangeListStopsAtBothEndsFromAnywhere(t *testing.T) {
	m := manyChanges(t, 60)
	last := len(m.currentRows()) - 1

	for _, start := range []int{0, 1, last / 2, last - 1, last} {
		for _, k := range pageKeys {
			at := m
			at.changeCursor = start
			// Enough presses to run off either end several times over.
			for i := 0; i < 12; i++ {
				at = press(at, k.key)
			}
			if at.changeCursor < 0 || at.changeCursor > last {
				t.Fatalf("%s from %d ran out of range to %d", k.name, start, at.changeCursor)
			}
			want := 0
			if k.down {
				want = last
			}
			if at.changeCursor != want {
				t.Errorf("%s from %d settled at %d, want %d", k.name, start, at.changeCursor, want)
			}
		}
	}
}

func TestPagingNeverLandsOnAGroupHeader(t *testing.T) {
	// A page is counted in drawn lines, which include the group headers, so a
	// jump that crosses the boundary has to land on a change either way.
	m := manyChanges(t, 40)
	info := m.projects["/p"].Info
	info.ArchivedChanges = append(info.ArchivedChanges,
		scanner.ChangeInfo{Name: "arch-one"}, scanner.ChangeInfo{Name: "arch-two"})
	m.projects["/p"] = scanner.ProjectStatus{Info: info}

	at := m
	for i := 0; i < 20; i++ {
		at = press(at, tea.KeyPressMsg{Code: tea.KeyPgDown})
		if _, ok := at.selectedRow(); !ok {
			t.Fatalf("after %d pages the selection is not a change", i+1)
		}
	}
}

// --- 2.6 and 2.7 the selection and the filter ---

func TestPagingRemembersTheSelection(t *testing.T) {
	m := manyChanges(t, 60)
	paged := press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	want, ok := paged.selectedRow()
	if !ok {
		t.Fatal("expected a selection")
	}
	if paged.selectedKey != want.key() {
		t.Errorf("selectedKey = %q, want %q: a page move must be remembered like a row move",
			paged.selectedKey, want.key())
	}
}

func TestPagingMovesWithinTheFilter(t *testing.T) {
	m := manyChanges(t, 60)
	m.searchInput.SetValue("aa")
	m.syncCursor()
	n := len(m.currentRows())
	if n == 0 || n >= 60 {
		t.Fatalf("the filter left %d rows, want a narrowed list", n)
	}
	end := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
	if end.changeCursor != n-1 {
		t.Errorf("G moved to %d, want the last of the %d matches", end.changeCursor, n)
	}
}

// --- 1.3 a page is what the surface shows ---

func TestAPageGrowsWithTheTerminal(t *testing.T) {
	short := manyChanges(t, 60)
	short.height = 20
	short.recalcLayout()

	tall := manyChanges(t, 60)
	tall.height = 50
	tall.recalcLayout()

	if tall.listPage() <= short.listPage() {
		t.Errorf("a page is %d rows at 20 high and %d at 50; it must grow",
			short.listPage(), tall.listPage())
	}
	if press(tall, tea.KeyPressMsg{Code: tea.KeyPgDown}).changeCursor <=
		press(short, tea.KeyPressMsg{Code: tea.KeyPgDown}).changeCursor {
		t.Error("and the key must move further in the taller pane")
	}
}

func TestListPageIsOneOutsideTheProjectLevel(t *testing.T) {
	m := manyChanges(t, 10)
	m.level = levelChange
	if got := m.listPage(); got != 1 {
		t.Errorf("got %d, want 1 where there is no list", got)
	}
}

// --- 5.1 the whole rule, as one table ---

func TestTheKeysActOnWhateverHoldsTheKeyboard(t *testing.T) {
	type moved struct{ list, doc bool }

	cases := []struct {
		name  string
		build func(*testing.T) model
		want  moved
		of    func(model) int
	}{
		{"the changes table", func(t *testing.T) model { return manyChanges(t, 60) },
			moved{list: true}, func(m model) int { return m.changeCursor }},
		{"a spec list", func(t *testing.T) model { return pagingModel(t, tabSpecs) },
			moved{list: true}, func(m model) int { return m.specCursor }},
		{"a spec document", func(t *testing.T) model {
			return press(pagingModel(t, tabSpecs), tea.KeyPressMsg{Code: tea.KeyTab})
		}, moved{doc: true}, func(m model) int { return m.specCursor }},
		{"the properties rows", func(t *testing.T) model { return pagingModel(t, tabProperties) },
			moved{list: true}, func(m model) int { return m.propSection }},
		{"a properties document", func(t *testing.T) model {
			return press(pagingModel(t, tabProperties), tea.KeyPressMsg{Code: tea.KeyTab})
		}, moved{doc: true}, func(m model) int { return m.propSection }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.build(t)
			after := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})

			if gotList := tc.of(after) != tc.of(m); gotList != tc.want.list {
				t.Errorf("list moved = %v, want %v", gotList, tc.want.list)
			}
			if gotDoc := after.docViewport.YOffset() != m.docViewport.YOffset(); gotDoc != tc.want.doc {
				t.Errorf("document moved = %v, want %v", gotDoc, tc.want.doc)
			}
		})
	}
}

// --- 5.4 the line keys are untouched ---

func TestLineKeysAreUnchangedOnEverySurface(t *testing.T) {
	m := manyChanges(t, 60)
	if press(m, tea.KeyPressMsg{Code: 'j', Text: "j"}).changeCursor != 1 {
		t.Error("j moves the change list one row")
	}
	s := pagingModel(t, tabSpecs)
	if press(s, tea.KeyPressMsg{Code: 'j', Text: "j"}).specCursor != 1 {
		t.Error("j moves the spec list one row")
	}
	p := pagingModel(t, tabProperties)
	if press(p, tea.KeyPressMsg{Code: 'j', Text: "j"}).propSection != 1 {
		t.Error("j moves the properties rows one row")
	}
}

// --- 4.x the picker ---

func pickerPagingModel(t *testing.T, n int) model {
	t.Helper()
	m := manyChanges(t, 2)
	m.pickerOpen = true
	m.pickerLoaded = true
	for i := 0; i < n; i++ {
		name := string(rune('a'+i/26)) + string(rune('a'+i%26))
		m.pickerAll = append(m.pickerAll, projectRow{path: "/" + name, display: name})
	}
	m.recalcLayout()
	return m
}

func TestPickerPagesAndJumps(t *testing.T) {
	m := pickerPagingModel(t, 40)
	page := m.pickerPage()
	if page < 2 {
		t.Fatalf("a page is %d rows, too small to test with", page)
	}

	down := press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	if down.pickerCursor != page {
		t.Errorf("pgdown moved to %d, want a page of %d", down.pickerCursor, page)
	}
	if press(down, tea.KeyPressMsg{Code: tea.KeyPgUp}).pickerCursor != 0 {
		t.Error("pgup did not come back")
	}
	half := press(m, tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	if half.pickerCursor != page/2 {
		t.Errorf("ctrl+d moved to %d, want %d", half.pickerCursor, page/2)
	}

	last := len(m.pickerVisibleRows()) - 1
	if got := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"}).pickerCursor; got != last {
		t.Errorf("G moved to %d, want %d: the jumps are unchanged", got, last)
	}
	if got := press(press(m, tea.KeyPressMsg{Code: 'G', Text: "G"}),
		tea.KeyPressMsg{Code: 'g', Text: "g"}).pickerCursor; got != 0 {
		t.Errorf("g moved to %d, want the first project", got)
	}
}

func TestPickerPagingStopsAtBothEnds(t *testing.T) {
	m := pickerPagingModel(t, 40)
	last := len(m.pickerVisibleRows()) - 1
	for _, k := range pageKeys {
		at := m
		for i := 0; i < 12; i++ {
			at = press(at, k.key)
		}
		want := 0
		if k.down {
			want = last
		}
		if at.pickerCursor != want {
			t.Errorf("%s settled at %d, want %d", k.name, at.pickerCursor, want)
		}
	}
}

func TestPickerPagingMovesWithinTheFilter(t *testing.T) {
	m := pickerPagingModel(t, 40)
	m.pickerInput.SetValue("aa")
	m.pickerSync()
	n := len(m.pickerVisibleRows())
	if n == 0 || n >= 40 {
		t.Fatalf("the filter left %d rows, want a narrowed list", n)
	}
	if got := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"}).pickerCursor; got != n-1 {
		t.Errorf("G moved to %d, want the last of the %d matches", got, n)
	}
}
