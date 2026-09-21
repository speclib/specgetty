package ui

import (
	"regexp"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// Note for anyone wondering where the colour-forcing helper went: lipgloss v2
// has no global colour profile. Style.Render always emits, and downgrading
// happens at the output writer. Under v1 these tests needed a helper to stop
// being silently vacuous, because Render checked for a terminal and stripped
// styling when there was none. That trap is gone.

func lineWidths(s string) []int {
	lines := strings.Split(s, "\n")
	out := make([]int, len(lines))
	for i, l := range lines {
		out[i] = lipgloss.Width(l)
	}
	return out
}

func TestPickerLinesAllHaveTheSameWidth(t *testing.T) {
	// Note: this alone does NOT catch an overflowing row. lipgloss pads the
	// wrapped remainder, so every line still measures the same width. The
	// detector for that is TestPickerDoesNotOverflowItsBox, which counts lines.
	// This test guards ragged output from other causes.

	for _, cursor := range []int{0, 1} {
		m := makePickerModel()
		m.pickerOpen = true
		m.width, m.height = 100, 30
		m.pickerCursor = cursor

		widths := lineWidths(m.renderPicker())
		for i, w := range widths {
			if w != widths[0] {
				t.Errorf("cursor %d: line %d is %d columns, line 0 is %d; the highlighted row must line up with the rest",
					cursor, i, w, widths[0])
			}
		}
	}
}

func TestPickerDoesNotOverflowItsBox(t *testing.T) {

	m := makePickerModel()
	m.pickerOpen = true
	m.width, m.height = 100, 30

	boxWidth, boxHeight := m.pickerBox()
	lines := strings.Split(m.renderPicker(), "\n")

	if got := lipgloss.Width(lines[0]); got != boxWidth {
		t.Errorf("box rendered %d columns wide, want %d", got, boxWidth)
	}
	// An overflowing row is wrapped rather than rejected, so extra lines are
	// the symptom. lipgloss Style.Width covers content plus padding and
	// excludes the border, which is the trap this guards.
	if len(lines) != boxHeight {
		t.Errorf("box rendered %d lines, want %d: extra lines mean a row overflowed and wrapped",
			len(lines), boxHeight)
	}
}

func TestPickerWidthsHoldWithAMatchHintColumn(t *testing.T) {

	m := makePickerModel()
	m.pickerOpen = true
	m.width, m.height = 100, 30
	m.pickerInput.SetValue(":flake")
	m.pickerSync()

	widths := lineWidths(m.renderPicker())
	for i, w := range widths {
		if w != widths[0] {
			t.Errorf("line %d is %d columns, line 0 is %d, with a hint column present", i, w, widths[0])
		}
	}
}

func TestPickerWidthsHoldAtAwkwardTerminalSizes(t *testing.T) {

	for _, size := range []struct{ w, h int }{
		{60, 20}, {61, 21}, {80, 24}, {100, 30}, {201, 51},
	} {
		m := makePickerModel()
		m.pickerOpen = true
		m.width, m.height = size.w, size.h
		m.pickerCursor = 1

		widths := lineWidths(m.renderPicker())
		for i, w := range widths {
			if w != widths[0] {
				t.Errorf("at %dx%d: line %d is %d columns, line 0 is %d",
					size.w, size.h, i, w, widths[0])
				break
			}
		}
	}
}

func TestPanelTopBorderMatchesTheBoxWidth(t *testing.T) {

	// The top border is assembled by hand rather than by lipgloss, so it can
	// drift from the box it caps. It did: it was one column short.
	for _, width := range []int{40, 58, 80, 98, 120} {
		m := makeListModel()
		m.width, m.height = width+2, 30

		widths := lineWidths(m.renderPanel(viewDetail, width, 5, "content"))
		for i, w := range widths {
			if w != widths[0] {
				t.Errorf("width %d: line %d is %d columns, top border is %d",
					width, i, w, widths[0])
				break
			}
		}
	}
}

func TestPanelTopBorderHandlesALongTitle(t *testing.T) {

	// A title longer than the box must not make the border negative or ragged.
	m := makeListModel()
	m.width, m.height = 30, 20
	m.displayNames = []string{strings.Repeat("very-long-project-name", 3)}

	widths := lineWidths(m.renderPanel(viewDetail, 28, 4, "content"))
	for i, w := range widths {
		if w < widths[len(widths)-1] {
			t.Errorf("line %d is %d columns, narrower than the bottom border at %d",
				i, w, widths[len(widths)-1])
		}
	}
}

// --- the panel's content width ---

// The panel's content width used to be derived in three files, each spelling
// out m.width - 2 for itself, while the comments on docRegion and specsSplit
// both warned that these must agree. These tests hold the agreement rather
// than restating it in a fourth comment.

func TestPanelContentWidthIsWhatThePanelActuallyDraws(t *testing.T) {
	// Not m.width - 2 asserted against itself: the panel is rendered and the
	// region between its borders is measured. If renderPanel ever draws its
	// content at a different width than panelContentWidth reports, the
	// document wraps somewhere other than where the viewport thinks it does,
	// and the scroll percentage starts lying.
	for _, width := range []int{60, 72, 100, 137} {
		m := makeViewModel()
		m.width = width
		m.recalcLayout()

		marker := strings.Repeat("x", m.panelContentWidth())
		rendered := m.renderPanel(viewDetail, m.width, 3, marker)

		var found bool
		for _, line := range strings.Split(rendered, "\n") {
			if !strings.Contains(line, "x") {
				continue
			}
			found = true
			// The marker fills the content region exactly, so the whole line
			// is the two borders plus that region.
			if got := lipgloss.Width(line); got != width {
				t.Errorf("at width %d the content row measured %d", width, got)
			}
			if strings.Count(line, "x") != m.panelContentWidth() {
				t.Errorf("at width %d the panel drew %d content columns, panelContentWidth says %d",
					width, strings.Count(line, "x"), m.panelContentWidth())
			}
		}
		if !found {
			t.Errorf("at width %d the marker never reached the panel", width)
		}
	}
}

func TestPanelContentWidthNeverGoesBelowOne(t *testing.T) {
	// renderFrame refuses to draw under 60 columns, but recalcLayout runs
	// whenever a size arrives, and a viewport set to a negative width is a
	// panic waiting for the next resize.
	for _, width := range []int{1, 2, 3} {
		m := makeViewModel()
		m.width = width
		if got := m.panelContentWidth(); got < 1 {
			t.Errorf("at terminal width %d the content width came out %d", width, got)
		}
	}
}

func TestEveryConsumerOfTheContentWidthAgrees(t *testing.T) {
	// The drift this guards against is a caller working the width out for
	// itself again. Each case below is a place that used to.
	for _, width := range []int{60, 80, 120} {
		m := makeViewModel()
		m.width = width

		t.Run("log viewport", func(t *testing.T) {
			m.logVisible = true
			m.recalcLayout()
			if got := m.logViewport.Width(); got != m.panelContentWidth() {
				t.Errorf("log viewport is %d wide at terminal width %d, panel content is %d",
					got, width, m.panelContentWidth())
			}
			m.logVisible = false
		})

		t.Run("document in an open change", func(t *testing.T) {
			m.level = levelChange
			m.recalcLayout()
			if w, _ := m.docRegion(); w != m.contentBoxWidth() {
				t.Errorf("docRegion is %d wide at terminal width %d, the content box is %d",
					w, width, m.contentBoxWidth())
			}
			m.level = levelProject
		})

		t.Run("document in the specs split", func(t *testing.T) {
			m.detailTab = tabSpecs
			m.recalcLayout()
			// The specs tab splits the whole region into two bordered boxes, so
			// the document gets what is inside the right-hand one.
			_, contentOuter := specsSplit(m.panelContentWidth())
			want := contentOuter - boxChrome
			if w, _ := m.docRegion(); w != want {
				t.Errorf("docRegion is %d wide at terminal width %d, the right box's inside is %d",
					w, width, want)
			}
			m.detailTab = tabChanges
		})
	}
}

// --- the inset ---

// stripANSI is deliberately not lipgloss.Width: these tests care about which
// column a character sits in, and a width alone cannot tell a leading blank
// from a missing one.
func panelRows(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(ansi.Strip(l), "│") {
			out = append(out, ansi.Strip(l))
		}
	}
	return out
}

func TestPanelContentStartsOneColumnInsideTheBorder(t *testing.T) {
	// The misalignment this change exists to remove: the project header used
	// to sit one column in and everything under it flush against the border.
	m := makeViewModel()
	m.syncDocument()

	rows := panelRows(m.renderFrame())
	if len(rows) < 6 {
		t.Fatalf("expected a drawn panel, got %d rows", len(rows))
	}
	for i, r := range rows {
		body := []rune(r)
		if len(body) < 3 {
			continue
		}
		if body[1] != ' ' {
			t.Errorf("row %d starts flush against the border: %q", i, r)
		}
		if body[len(body)-2] != ' ' {
			t.Errorf("row %d ends flush against the border: %q", i, r)
		}
	}
}

func TestTheHeaderAndTheRowsBelowItShareAColumn(t *testing.T) {
	// The header carries padding of its own. If it kept the horizontal half it
	// would sit two columns in while everything under it sat at one, which is
	// the one thing the inset cannot be allowed to reintroduce.
	m := makeViewModel()
	m.syncDocument()

	var project, box, table int
	for _, r := range panelRows(m.renderFrame()) {
		col := len([]rune(r)) - len([]rune(strings.TrimLeft(r, "│ ")))
		switch {
		case strings.Contains(r, "/p") && project == 0:
			project = col
		case strings.Contains(r, "╭") && box == 0:
			box = col
		case strings.Contains(r, "name") && strings.Contains(r, "tasks") && table == 0:
			table = col
		}
	}
	if project == 0 || box == 0 || table == 0 {
		t.Fatalf("did not find all three rows (header %d, box %d, table %d)", project, box, table)
	}
	// The header and the content border are siblings under the panel, so they
	// share a column. The table is inside that border, so it sits further in by
	// exactly the border and its inset.
	if project != box {
		t.Errorf("the header starts at column %d and the content border at %d; they must agree", project, box)
	}
	if table != box+2 {
		t.Errorf("the table starts at column %d, want %d: one for the border, one for its inset", table, box+2)
	}

	// The tab bar is measured separately and on the raw line. Its chips are
	// painted blocks that begin with a space of their own, so on stripped text
	// they look one column further in than they are drawn. What must hold is
	// that the chip's paint starts right after the gutter.
	for _, l := range strings.Split(m.renderFrame(), "\n") {
		if !strings.Contains(ansi.Strip(l), "changes") || !strings.Contains(ansi.Strip(l), "properties") {
			continue
		}
		if !strings.HasPrefix(l, "\x1b[32m│\x1b[m \x1b[") {
			t.Errorf("the tab bar's first chip does not begin one column inside the border: %q", l)
		}
		return
	}
	t.Fatal("no tab bar in the frame")
}

func TestSelectedRowKeepsTheGutterOnBothSides(t *testing.T) {
	// The picker already draws its selected row this way. This is the change
	// list catching up, not a new look.
	m := makeViewModel()
	m.syncDocument()

	for _, l := range strings.Split(m.renderFrame(), "\n") {
		if !strings.Contains(l, "\x1b[30;42m") {
			continue
		}
		plain := []rune(ansi.Strip(l))
		if plain[1] != ' ' || plain[len(plain)-2] != ' ' {
			t.Errorf("the highlighted row touches a border: %q", ansi.Strip(l))
		}
		return
	}
	t.Fatal("no highlighted row in the frame")
}

func TestNavBarIndentsItsTextAndStillPaintsEveryColumn(t *testing.T) {
	// A strip whose job is to mark the bottom edge of the screen must not have
	// holes in it. The fill between the hints and the version used to be bare
	// spaces, which left one.
	m := makeViewModel()
	bar := m.renderNavBar()

	plain := []rune(ansi.Strip(bar))
	if plain[0] != ' ' || plain[len(plain)-1] != ' ' {
		t.Errorf("the nav bar text is not inset: %q", string(plain))
	}
	// Every run of spaces inside the bar must carry the background. A bare run
	// is a hole.
	for _, seg := range strings.Split(bar, "\x1b[m") {
		if strings.TrimSpace(seg) == "" && strings.Contains(seg, "  ") && !strings.Contains(seg, "48;5;236") {
			t.Errorf("an unpainted run of %d spaces in the nav bar", len(seg))
		}
	}
}

// --- column separation ---

func TestTableColumnsAreSeparatedByTwoBlankColumns(t *testing.T) {
	// A value that exactly fills its column used to sit one space from the
	// value beside it, which reads as one run-on field rather than two.
	defs := []fieldDef[changeRow]{
		{header: "name", width: 0, value: func(r changeRow) string { return r.ci.Name }},
		{header: "tasks", width: 7, value: func(r changeRow) string { return "6/6" }},
	}
	// The name is built to fill the flexible column exactly, because that is
	// the case the wider gap exists for. A shorter name leaves its own padding
	// behind and would make this pass for the wrong reason.
	_, widths := layoutFields(defs, 50)
	rows := []filtered[changeRow]{{row: changeRow{ci: scanner.ChangeInfo{
		Name: strings.Repeat("n", widths[0])}}}}

	body := ansi.Strip(strings.Split(renderTable(rows, defs, 0, 50, 2), "\n")[1])

	i := strings.LastIndex(body, "n")
	if i < 0 {
		t.Fatalf("the name never reached the row: %q", body)
	}
	if got := body[i+1 : i+4]; got != "  6" {
		t.Errorf("want two blank columns between a full name and the next value, got %q in %q", got, body)
	}
}

func TestTableRowsStillMeasureTheTableWidth(t *testing.T) {
	// What this does NOT catch, checked by breaking it: the gap count and the
	// separator disagreeing. fitCell normalises the joined row to the table
	// width in both directions, padding a short row and truncating a long one,
	// so the measurement comes out right either way. The disagreement is
	// caught by TestTableColumnsAreSeparatedByTwoBlankColumns in one direction
	// and TestWiderGapsShiftWhenColumnsAreDropped in the other. This test
	// pins the normalisation itself, which is what keeps a row from ever
	// running past its panel.
	defs := []fieldDef[changeRow]{
		{header: "name", width: 0, value: func(r changeRow) string { return r.ci.Name }},
		{header: "tasks", width: 7, value: func(r changeRow) string { return "6/6" }},
		{header: "specs", width: 5, value: func(r changeRow) string { return "1" }},
	}
	rows := []filtered[changeRow]{{row: changeRow{ci: scanner.ChangeInfo{Name: "alpha"}}}}

	for _, width := range []int{40, 56, 80, 120} {
		for i, line := range strings.Split(renderTable(rows, defs, -1, width, 2), "\n") {
			if got := ansi.StringWidth(line); got != width {
				t.Errorf("at table width %d, row %d measured %d", width, i, got)
			}
		}
	}
}

func TestWiderGapsShiftWhenColumnsAreDropped(t *testing.T) {
	// Recorded rather than discovered on a laptop: the gaps come out of the
	// flexible column's budget, so a column starts being dropped at a slightly
	// wider terminal than it used to.
	defs := []fieldDef[changeRow]{
		{header: "name", width: 0, value: func(r changeRow) string { return r.ci.Name }},
		{header: "tasks", width: 7, value: func(r changeRow) string { return "6/6" }},
		{header: "specs", width: 5, value: func(r changeRow) string { return "1" }},
	}
	// Measured, not predicted. With one-column gaps the three columns survived
	// down to 26; with two they survive to 28, and the two-column form to 21
	// rather than 19. Two columns of gap, two columns of threshold.
	for _, tc := range []struct {
		width, want int
	}{
		{60, 3},
		{28, 3}, // the narrowest that still fits all three
		{27, 2}, // specs is dropped here
		{21, 2}, // the narrowest that still fits two
		{20, 1}, // tasks goes too
	} {
		kept, _ := layoutFields(defs, tc.width)
		if len(kept) != tc.want {
			names := make([]string, len(kept))
			for i, d := range kept {
				names[i] = d.header
			}
			t.Errorf("at width %d, %d columns survive (%v), want %d", tc.width, len(kept), names, tc.want)
		}
	}
}

// --- the content border ---

// borderRow reports the index of the content border's top row, and -1 when
// there is none. The panel's own top border starts the string, so a corner
// anywhere else is the content box.
func borderRow(frame string) int {
	for i, l := range strings.Split(frame, "\n") {
		p := ansi.Strip(l)
		if strings.Contains(p, "╭") && !strings.HasPrefix(p, "╭") {
			return i
		}
	}
	return -1
}

func rowOf(frame, needle string) int {
	for i, l := range strings.Split(frame, "\n") {
		if strings.Contains(ansi.Strip(l), needle) {
			return i
		}
	}
	return -1
}

func TestTheContentBorderSitsDirectlyUnderTheTabBar(t *testing.T) {
	// No blank row between them: the chips have to read as tabs belonging to
	// the box, which is the whole point of drawing it.
	m := makeViewModel()
	m.syncDocument()
	frame := m.renderFrame()

	tabs := rowOf(frame, "changes")
	box := borderRow(frame)
	if tabs < 0 || box < 0 {
		t.Fatalf("frame is missing the tab bar (%d) or the border (%d)", tabs, box)
	}
	if box != tabs+1 {
		t.Errorf("the border is on row %d and the tab bar on row %d; want them adjacent", box, tabs)
	}
}

func TestTheSearchPromptIsInsideTheBorder(t *testing.T) {
	// It reports how many rows the box is showing, so it belongs to the box.
	// The config tab's filename names the box instead, and goes above it; that
	// pairing is what openspec/specs/panel-layout calls chrome and content.
	m := makeViewModel()
	m.searchInput.SetValue("alpha")
	m.searchFocused = true
	m.syncDocument()
	frame := m.renderFrame()

	box := borderRow(frame)
	prompt := rowOf(frame, "shown")
	if box < 0 || prompt < 0 {
		t.Fatalf("frame is missing the border (%d) or the prompt (%d)", box, prompt)
	}
	if prompt < box {
		t.Errorf("the search prompt is on row %d, above the border on row %d", prompt, box)
	}
}

func TestAnOpenChangeBoxesItsArtifactAndNotItsHeader(t *testing.T) {
	m := viewWithTasks(t, 4, 4)
	m.level = levelChange
	m.syncDocument()
	frame := m.renderFrame()

	name := rowOf(frame, "alpha")
	box := borderRow(frame)
	if name < 0 || box < 0 {
		t.Fatalf("frame is missing the change name (%d) or the border (%d)", name, box)
	}
	if name > box {
		t.Errorf("the change name is on row %d, inside the border on row %d", name, box)
	}
}

func TestTheChangeTableKeepsItsColumnsInsideTheBorder(t *testing.T) {
	// The border takes four columns off the table. At the narrowest terminal
	// specgetty draws, that must not cost a column that used to fit.
	m := makeViewModel()
	m.width, m.height = 60, 20
	m.recalcLayout()
	m.syncDocument()

	frame := m.renderFrame()
	header := ""
	for _, l := range strings.Split(frame, "\n") {
		if p := ansi.Strip(l); strings.Contains(p, "name") && strings.Contains(p, "tasks") {
			header = p
			break
		}
	}
	if header == "" {
		t.Fatal("no table header in the frame at 60x20")
	}
	for _, col := range []string{"name", "tasks", "specs"} {
		if !strings.Contains(header, col) {
			t.Errorf("the %q column was dropped at 60 columns: %q", col, header)
		}
	}
}

func TestEveryViewFitsTheTerminalWithTheBorderDrawn(t *testing.T) {
	for _, sz := range [][2]int{{100, 30}, {72, 24}, {60, 20}} {
		for _, tab := range []int{tabChanges, tabSpecs, tabProperties} {
			m := makeViewModel()
			m.width, m.height = sz[0], sz[1]
			m.detailTab = tab
			m.focus = m.defaultFocus()
			m.recalcLayout()
			m.syncDocument()

			lines := strings.Split(m.renderFrame(), "\n")
			if len(lines) != m.height {
				t.Errorf("tab %d at %dx%d is %d lines, want %d", tab, sz[0], sz[1], len(lines), m.height)
			}
			for i, l := range lines {
				if w := lipgloss.Width(l); w != m.width {
					t.Errorf("tab %d at %dx%d: line %d is %d columns, want %d", tab, sz[0], sz[1], i, w, m.width)
					break
				}
			}
		}
	}
}

// --- the lit border follows the keyboard ---

// litContentBorders reports, left to right and top to bottom, whether each
// content box inside the panel is drawn in the active colour.
//
// Rows are filtered to those the panel's own side border starts, which keeps
// the log panel out: it is a sibling of the panel, not a box inside it. The
// panel's top border is built by hand and colours only its title, so it is not
// readable this way and is checked through viewFocused instead.
var contentBoxTop = regexp.MustCompile(`\x1b\[([0-9;]+)m╭`)

func litContentBorders(frame string) []bool {
	var out []bool
	for _, l := range strings.Split(frame, "\n") {
		if !strings.HasPrefix(ansi.Strip(l), "│") {
			continue
		}
		for _, m := range contentBoxTop.FindAllStringSubmatch(l, -1) {
			out = append(out, m[1] == "32")
		}
	}
	return out
}

func specsFrame(t *testing.T, focus int, logOpen bool) string {
	t.Helper()
	m := makeSpecsModel(t, map[string]string{"alpha": numberedDoc(80), "beta": numberedDoc(60)})
	m.logVisible = logOpen
	m.focus = focus
	m.recalcLayout()
	m.syncDocument()
	return m.renderFrame()
}

func TestTheLitBorderIsTheOneHoldingTheKeyboard(t *testing.T) {
	for _, tc := range []struct {
		name      string
		focus     int
		wantInner []bool // list, content
	}{
		{"the spec list has it", focusListPane, []bool{true, false}},
		{"the spec content has it", focusContentPane, []bool{false, true}},
		{"the log panel has it", focusLog, []bool{false, false}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			frame := specsFrame(t, tc.focus, true)
			inner := litContentBorders(frame)
			if len(inner) != len(tc.wantInner) {
				t.Fatalf("found %d content borders, want %d", len(inner), len(tc.wantInner))
			}
			for i := range inner {
				if inner[i] != tc.wantInner[i] {
					t.Errorf("content border %d lit = %v, want %v", i, inner[i], tc.wantInner[i])
				}
			}
		})
	}
}

func TestTheTwoSpecsBordersAreNeverLitTogether(t *testing.T) {
	for _, focus := range []int{focusListPane, focusContentPane, focusLog} {
		inner := litContentBorders(specsFrame(t, focus, true))
		n := 0
		for _, l := range inner {
			if l {
				n++
			}
		}
		if n > 1 {
			t.Errorf("focus %d lit %d of the two specs borders; at most one holds the keyboard", focus, n)
		}
	}
}

func TestThePanelBorderStillMeansWhatItAlwaysMeant(t *testing.T) {
	// The one thing this change must not move. The panel is lit whenever the
	// keyboard is anywhere in the detail area, which is what it did before the
	// content had borders of its own.
	for _, tc := range []struct {
		focus int
		want  bool
	}{
		{focusDetail, true},
		{focusListPane, true},
		{focusContentPane, true},
		{focusLog, false},
	} {
		m := makeViewModel()
		m.focus = tc.focus
		if got := m.viewFocused(viewDetail); got != tc.want {
			t.Errorf("with focus %d the detail panel reads lit=%v, want %v", tc.focus, got, tc.want)
		}
		if got := m.viewFocused(viewLog); got == tc.want && tc.focus != focusLog {
			t.Errorf("with focus %d the log panel reads lit=%v, want the opposite", tc.focus, got)
		}
	}
}
