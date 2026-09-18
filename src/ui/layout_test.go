package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
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
			if w, _ := m.docRegion(); w != m.panelContentWidth() {
				t.Errorf("docRegion is %d wide at terminal width %d, panel content is %d",
					w, width, m.panelContentWidth())
			}
			m.level = levelProject
		})

		t.Run("document in the specs split", func(t *testing.T) {
			m.detailTab = tabSpecs
			m.recalcLayout()
			_, want := specsSplit(m.panelContentWidth())
			if w, _ := m.docRegion(); w != want {
				t.Errorf("docRegion is %d wide at terminal width %d, the split's content half is %d",
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

	var project, table int
	for _, r := range panelRows(m.renderFrame()) {
		col := len([]rune(r)) - len([]rune(strings.TrimLeft(r, "│ ")))
		switch {
		case strings.Contains(r, "/p") && project == 0:
			project = col
		case strings.Contains(r, "name") && strings.Contains(r, "tasks") && table == 0:
			table = col
		}
	}
	if project == 0 || table == 0 {
		t.Fatalf("did not find both rows (header %d, table %d)", project, table)
	}
	if project != table {
		t.Errorf("the header starts at column %d and the table at %d; they must agree", project, table)
	}

	// The tab bar is measured separately and on the raw line. Its chips are
	// painted blocks that begin with a space of their own, so on stripped text
	// they look one column further in than they are drawn. What must hold is
	// that the chip's paint starts right after the gutter.
	for _, l := range strings.Split(m.renderFrame(), "\n") {
		if !strings.Contains(ansi.Strip(l), "changes") || !strings.Contains(ansi.Strip(l), "config") {
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
