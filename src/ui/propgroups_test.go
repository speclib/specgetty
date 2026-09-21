package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// propertiesList renders the list half alone, stripped, for reading.
func propertiesList(m model) []string {
	m.recalcLayout()
	m.syncDocument()
	width := m.panelContentWidth()
	listOuter, _ := m.propertiesSplit(width)
	out := ansi.Strip(m.renderPropertiesTab(width, m.mainPanelHeight()-5))

	var rows []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < listOuter {
			continue
		}
		// The list is the left box; everything past its border belongs to the
		// content beside it.
		rows = append(rows, strings.TrimRight(line[:listOuter], " "))
	}
	return rows
}

// --- 2.3 to 2.5 the drawn shape ---

func TestTheListIsDrawnAsGroupsWithHeaders(t *testing.T) {
	rows := propertiesList(twoSchemaModel())
	joined := strings.Join(rows, "\n")

	for _, want := range []string{"PROJECT", "  config", "  store", "SCHEMAS",
		"  spec-driven", "  tinychange"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the list should draw %q:\n%s", want, joined)
		}
	}

	// The order, and a blank line between the groups.
	var shape []string
	for _, r := range rows {
		s := strings.TrimSpace(strings.TrimLeft(r, "│ "))
		if s == "" || strings.ContainsAny(s, "╭╰─") {
			continue // the box's own borders
		}
		shape = append(shape, s)
	}
	want := "PROJECT,config,store,SCHEMAS,spec-driven,tinychange"
	if got := strings.Join(shape, ","); !strings.HasPrefix(got, want) {
		t.Errorf("got %q, want it to start %q", got, want)
	}
}

func TestAHeaderIsNotIndentedAndARowIs(t *testing.T) {
	for _, r := range propertiesList(twoSchemaModel()) {
		inner := strings.TrimPrefix(r, "│ ")
		if strings.HasPrefix(inner, "PROJECT") || strings.HasPrefix(inner, "SCHEMAS") {
			return // a header sits flush against the border
		}
	}
	t.Error("no header was drawn flush; rows are indented under one")
}

func TestABlankLineSeparatesTheGroups(t *testing.T) {
	lines := propertyLines(propSections(twoSchemaModel().projects["/work/specgetty"].Info))
	var blanks int
	for i, l := range lines {
		if l.section == lineOwner && l.header == "" {
			blanks++
			if i+1 >= len(lines) || lines[i+1].header == "" {
				t.Error("the blank line should sit directly above a group header")
			}
		}
	}
	if blanks != 1 {
		t.Errorf("got %d blank lines, want one between the two groups", blanks)
	}
	if lines[0].header != groupProject {
		t.Errorf("the first line is %+v, want the first header with no blank above it", lines[0])
	}
}

func TestAGroupKeepsItsHeaderWhenItHoldsNothing(t *testing.T) {
	// A project with no recorded schema. "No schemas" is an answer, and an
	// absent header cannot be told from one that was filtered away.
	m := twoSchemaModel()
	info := m.projects["/work/specgetty"].Info
	info.SchemaUsage = nil
	m.projects["/work/specgetty"] = scanner.ProjectStatus{Info: info}

	joined := strings.Join(propertiesList(m), "\n")
	if !strings.Contains(joined, "SCHEMAS") {
		t.Errorf("the schemas header should still be drawn:\n%s", joined)
	}
	if !strings.Contains(joined, "  config") {
		t.Errorf("the project group should still hold its rows:\n%s", joined)
	}
}

// --- 3.x the cursor ---

func pressN(m model, key rune, n int) model {
	for i := 0; i < n; i++ {
		m = press(m, tea.KeyPressMsg{Code: key, Text: string(key)})
	}
	return m
}

func TestTheCursorStepsFromOneGroupIntoTheNext(t *testing.T) {
	m := twoSchemaModel()
	m.recalcLayout()
	sections := m.currentSections()
	if sections[1].group != groupProject || sections[2].group != groupSchemas {
		t.Fatalf("setup: %+v", sections)
	}

	m.propSection = 1 // the last row of PROJECT
	after := pressN(m, 'j', 1)
	if after.propSection != 2 {
		t.Errorf("got row %d, want the first row of the next group", after.propSection)
	}
	back := pressN(after, 'k', 1)
	if back.propSection != 1 {
		t.Errorf("got row %d, want the last row of the previous group", back.propSection)
	}
}

func TestNoVerticalKeyLandsTheCursorOnAHeader(t *testing.T) {
	// The cursor indexes sections, and a header is not one, so this can only
	// fail if some key started indexing drawn lines instead.
	m := twoSchemaModel()
	m.recalcLayout()
	total := len(m.currentSections())

	for _, key := range []tea.KeyPressMsg{
		{Code: 'j', Text: "j"}, {Code: 'k', Text: "k"},
		{Code: tea.KeyDown}, {Code: tea.KeyUp},
		{Code: tea.KeyPgDown}, {Code: tea.KeyPgUp},
		{Code: 'd', Mod: tea.ModCtrl}, {Code: 'u', Mod: tea.ModCtrl},
		{Code: 'G', Text: "G"},
	} {
		for start := 0; start < total; start++ {
			m.propSection = start
			got := press(m, key).propSection
			if got < 0 || got >= total {
				t.Fatalf("%v from row %d left the cursor at %d, outside the rows",
					key, start, got)
			}
		}
	}
}

func TestGGAndGLandOnRowsNotHeaders(t *testing.T) {
	m := twoSchemaModel()
	m.recalcLayout()
	total := len(m.currentSections())

	end := press(m, tea.KeyPressMsg{Code: 'G', Text: "G"})
	if end.propSection != total-1 {
		t.Errorf("G landed on row %d, want the last row %d", end.propSection, total-1)
	}
	top := pressN(pressN(end, 'g', 1), 'g', 1)
	if top.propSection != 0 {
		t.Errorf("gg landed on row %d, want the first", top.propSection)
	}
}

func TestTheListScrollsByDrawnLines(t *testing.T) {
	// A pane too short for the whole list must still bring the selected row
	// into view, counting the headers and the spacer as the lines they take.
	m := twoSchemaModel()
	m.height = 16
	m.recalcLayout()
	m.propSection = len(m.currentSections()) - 1
	m.syncDocument()

	joined := strings.Join(propertiesList(m), "\n")
	if !strings.Contains(joined, "tinychange") {
		t.Errorf("the selected row should be on screen:\n%s", joined)
	}
}

func TestThePropertiesPaneNeverOpensOnABlankLine(t *testing.T) {
	m := twoSchemaModel()
	for _, height := range []int{12, 14, 16, 20, 40} {
		m.height = height
		m.recalcLayout()
		for row := 0; row < len(m.currentSections()); row++ {
			m.propSection = row
			m.syncDocument()
			rows := propertiesList(m)
			if len(rows) < 2 {
				continue
			}
			// rows[0] is the box's top border; rows[1] is its first line.
			if first := strings.TrimLeft(rows[1], "│ "); strings.TrimSpace(first) == "" {
				t.Errorf("height %d row %d: the pane opens on a blank line:\n%s",
					height, row, strings.Join(rows, "\n"))
			}
		}
	}
}

// --- 4.x the width ---

func TestTheListTakesItsFloorAndKeepsItsCap(t *testing.T) {
	m := twoSchemaModel()
	for _, c := range []struct{ term, want int }{
		{60, 18},  // the one-third cap gives way first
		{80, 20},  // the floor
		{100, 20}, // still the floor
		{140, 20},
		{200, 20},
	} {
		m.width = c.term
		m.recalcLayout()
		listOuter, _ := m.propertiesSplit(m.panelContentWidth())
		if listOuter != c.want {
			t.Errorf("%d columns: list is %d wide, want %d", c.term, listOuter, c.want)
		}
	}
}

func TestEveryExtraColumnGoesToTheContent(t *testing.T) {
	m := twoSchemaModel()
	m.width = 100
	m.recalcLayout()
	list100, content100 := m.propertiesSplit(m.panelContentWidth())
	m.width = 200
	m.recalcLayout()
	list200, content200 := m.propertiesSplit(m.panelContentWidth())

	if list100 != list200 {
		t.Errorf("the list moved from %d to %d columns; the divider should stay put",
			list100, list200)
	}
	if content200-content100 != 100 {
		t.Errorf("content grew by %d of the 100 extra columns", content200-content100)
	}
}

func TestNeitherHalfIsWiderThanThePanel(t *testing.T) {
	m := twoSchemaModel()
	m.width, m.height = 60, 20
	m.recalcLayout()
	width := m.panelContentWidth()
	listOuter, contentOuter := m.propertiesSplit(width)
	if listOuter+contentOuter+1 > width {
		t.Errorf("list %d and content %d exceed the panel's %d",
			listOuter, contentOuter, width)
	}
	if listOuter < 1 || contentOuter < 1 {
		t.Errorf("a half collapsed: list %d, content %d", listOuter, contentOuter)
	}
}

func TestThePropertiesFrameFitsAtEverySize(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {100, 30}, {200, 50}} {
		m := twoSchemaModel()
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

// TestTheHighlightIsOnTheSelectedRow is what makes the cursor's indexing
// observable. Without it, a renderer that highlighted drawn lines rather than
// sections would pass every other test in this file while highlighting the
// wrong row, or none.
func TestTheHighlightIsOnTheSelectedRow(t *testing.T) {
	m := twoSchemaModel()
	m.recalcLayout()
	sections := m.currentSections()

	for row := range sections {
		m.propSection = row
		m.syncDocument()
		width := m.panelContentWidth()
		listOuter, _ := m.propertiesSplit(width)
		out := m.renderPropertiesTab(width, m.mainPanelHeight()-5)

		// The exact string the renderer emits for a selected row. Compared
		// against that rather than against "something on this line is styled":
		// the content half beside the list is styled too, and matching it
		// would let a renderer highlighting the wrong row pass.
		want := selectedStyle.Width(listOuter - boxChrome).Render("  " + sections[row].label)
		if !strings.Contains(out, want) {
			t.Errorf("row %d (%q) is selected but is not the row drawn as selected:\n%s",
				row, sections[row].label, ansi.Strip(out))
		}
	}
}

func TestNoHeaderIsEverDrawnAsSelected(t *testing.T) {
	m := twoSchemaModel()
	m.recalcLayout()
	for row := range m.currentSections() {
		m.propSection = row
		m.syncDocument()
		width := m.panelContentWidth()
		out := m.renderPropertiesTab(width, m.mainPanelHeight()-5)

		for _, line := range strings.Split(out, "\n") {
			plain := strings.TrimSpace(ansi.Strip(line))
			for _, g := range propGroups {
				// A header carries the section-header style and never the
				// selection one, whatever the cursor is on.
				if strings.HasPrefix(strings.TrimLeft(plain, "│ "), g) &&
					strings.Contains(line, selectedStyle.Width(1).Render("")) {
					t.Errorf("row %d: the header %q is drawn as selected", row, g)
				}
			}
		}
	}
}

// TestThePropertiesTabRendersForEveryProjectOnThisMachine is task 5.3. The
// fixtures in this file are shapes this author wrote; real projects carry
// schema names, store paths and configurations nobody here chose.
//
// Skipped unless SPECGETTY_CORPUS names a directory of OpenSpec projects.
func TestThePropertiesTabRendersForEveryProjectOnThisMachine(t *testing.T) {
	root := os.Getenv("SPECGETTY_CORPUS")
	if root == "" {
		t.Skip("set SPECGETTY_CORPUS to a directory of OpenSpec projects to run this")
	}

	var seen, drawn int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() || filepath.Base(path) != "openspec" {
			return nil
		}
		project := filepath.Dir(path)
		_, st, err := scanner.ScanResolved(project)
		if err != nil || st.Info.Root == "" {
			return filepath.SkipDir
		}
		seen++

		m := twoSchemaModel()
		m.repoPaths = []string{project}
		m.displayNames = []string{filepath.Base(project)}
		m.projects = scanner.ProjectMap{project: st}
		m.cursor = 0

		for _, size := range []struct{ w, h int }{{60, 20}, {100, 30}, {200, 50}} {
			m.width, m.height = size.w, size.h
			m.recalcLayout()
			for row := range m.currentSections() {
				m.propSection = row
				m.syncDocument()
				width := m.panelContentWidth()
				listOuter, _ := m.propertiesSplit(width)
				out := m.renderPropertiesTab(width, m.mainPanelHeight()-5)
				for _, line := range strings.Split(out, "\n") {
					if w := ansi.StringWidth(line); w > width {
						t.Errorf("%s at %dx%d row %d: a line is %d columns, pane is %d",
							project, size.w, size.h, row, w, width)
						return filepath.SkipDir
					}
				}
				if listOuter < 1 {
					t.Errorf("%s: the list collapsed", project)
				}
				drawn++
			}
		}
		return filepath.SkipDir
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d projects, %d rows drawn at three sizes each", seen, drawn)
	if seen == 0 {
		t.Skip("no OpenSpec projects found under SPECGETTY_CORPUS")
	}
}
