package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// makeSpecsModel opens the specs tab on a project with the given specs.
func makeSpecsModel(t *testing.T, specs map[string]string) model {
	t.Helper()
	names := make([]string, 0, len(specs))
	for n := range specs {
		names = append(names, n)
	}
	// ParseProjectInfo sorts; mirror that so the cursor order is predictable.
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[j] < names[i] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	m := makeListModel()
	m.width, m.height = 100, 30
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		SpecCount:    len(names),
		SpecNames:    names,
		SpecContents: specs,
	}}}
	m.level = levelProject
	m.detailTab = tabSpecs
	m.recalcLayout()
	m.syncDocument()
	return m
}

func twoSpecs(t *testing.T) model {
	return makeSpecsModel(t, map[string]string{
		"alpha": numberedDoc(200),
		"beta":  numberedDoc(150),
	})
}

// --- 4.1 default focus ---

func TestSpecsFocusStartsOnTheList(t *testing.T) {
	m := twoSpecs(t)
	if m.specsFocus != specsFocusList {
		t.Errorf("focus is %d, want the list", m.specsFocus)
	}
	if m.docActive() {
		t.Error("the list holds the keyboard, so the document must not")
	}
}

func TestSpecsFocusResetsWhenLeavingTheTab(t *testing.T) {
	m := twoSpecs(t)
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.specsFocus != specsFocusContent {
		t.Fatal("expected the content to take the keyboard")
	}

	// Away to another tab and back.
	m = press(m, tea.KeyPressMsg{Code: '1', Text: "1"})
	m = press(m, tea.KeyPressMsg{Code: '2', Text: "2"})
	if m.specsFocus != specsFocusList {
		t.Errorf("focus is %d after returning, want the list", m.specsFocus)
	}

	// The arrow keys change tabs too, and must reset it the same way.
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m = press(m, tea.KeyPressMsg{Code: tea.KeyLeft})
	m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	if m.specsFocus != specsFocusList {
		t.Errorf("focus is %d after changing tab with the arrows, want the list", m.specsFocus)
	}
}

// --- 4.2 and 4.3 the tab cycle ---

func TestTabCyclesTheSpecsHalves(t *testing.T) {
	m := twoSpecs(t)
	m.logVisible = false

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.specsFocus != specsFocusContent {
		t.Fatalf("first tab gave focus %d, want the content", m.specsFocus)
	}
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.specsFocus != specsFocusList {
		t.Errorf("second tab gave focus %d, want the list", m.specsFocus)
	}
	if m.activeView != viewDetail {
		t.Errorf("activeView is %d, want viewDetail with the log closed", m.activeView)
	}
}

func TestTabCyclesThroughTheLogPanel(t *testing.T) {
	m := twoSpecs(t)
	m.logVisible = true
	m.recalcLayout()

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.specsFocus != specsFocusContent || m.activeView != viewDetail {
		t.Fatalf("first tab: focus=%d view=%d, want the content half", m.specsFocus, m.activeView)
	}

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.activeView != viewLog {
		t.Fatalf("second tab: view=%d, want the log panel", m.activeView)
	}

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.activeView != viewDetail || m.specsFocus != specsFocusList {
		t.Errorf("third tab: view=%d focus=%d, want back to the spec list", m.activeView, m.specsFocus)
	}
}

// --- 4.4 the keys follow the focus ---

func TestJAndKMoveTheSpecCursorWhenTheListIsFocused(t *testing.T) {
	m := twoSpecs(t)

	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.specCursor != 1 {
		t.Errorf("spec cursor is %d, want 1", m.specCursor)
	}
	if m.docViewport.YOffset() != 0 {
		t.Error("the document scrolled while the list held the keyboard")
	}
}

func TestJAndKScrollWhenTheContentIsFocused(t *testing.T) {
	m := twoSpecs(t)
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})

	before := m.specCursor
	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})

	if m.docViewport.YOffset() != 1 {
		t.Errorf("document offset is %d, want 1", m.docViewport.YOffset())
	}
	if m.specCursor != before {
		t.Errorf("the spec cursor moved to %d while the content held the keyboard", m.specCursor)
	}
}

func TestSpecCursorBoundsAreUnchanged(t *testing.T) {
	m := twoSpecs(t)

	m = press(m, tea.KeyPressMsg{Code: 'k', Text: "k"})
	if m.specCursor != 0 {
		t.Errorf("k at the first spec moved to %d", m.specCursor)
	}
	for i := 0; i < 5; i++ {
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	if m.specCursor != 1 {
		t.Errorf("the cursor ran past the last spec to %d", m.specCursor)
	}
}

// --- 4.5 every row reachable ---

func TestEverySpecRowIsReachable(t *testing.T) {
	m := twoSpecs(t)
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})

	seen := map[string]bool{}
	for {
		for _, row := range strings.Split(m.docViewport.View(), "\n") {
			if trimmed := strings.TrimSpace(row); trimmed != "" {
				seen[trimmed] = true
			}
		}
		if m.docViewport.AtBottom() {
			break
		}
		m.docViewport.PageDown()
	}
	for i := 1; i <= 200; i++ {
		want := fmt.Sprintf("row-%03d", i)
		if !seen[want] {
			t.Fatalf("%s of the spec was never displayed", want)
		}
	}
}

// --- 4.6 changing spec resets the position ---

func TestSelectingADifferentSpecStartsAtTheTop(t *testing.T) {
	m := twoSpecs(t)
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m.docViewport.SetYOffset(80)

	// Back to the list, then down one spec.
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})

	if m.docViewport.YOffset() != 0 {
		t.Errorf("the newly selected spec opened at offset %d, want 0", m.docViewport.YOffset())
	}
}

// --- 4.7 the position indicator follows the focus ---

func TestPositionReportedOnlyWhileTheContentIsFocused(t *testing.T) {
	m := twoSpecs(t)
	if pct := m.docScrollPercent(); pct != -1 {
		t.Errorf("the list holds the keyboard but a position of %d%% was reported", pct)
	}

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if pct := m.docScrollPercent(); pct != 0 {
		t.Errorf("the content holds the keyboard, position is %d%%, want 0%%", pct)
	}

	title := strings.Split(m.renderPanel(viewDetail, m.width-2, m.mainPanelHeight(), "x"), "\n")[0]
	if !strings.Contains(title, "0%") {
		t.Errorf("the title should carry the position:\n%s", title)
	}
}

// --- 4.8 wrapping respects the content half ---

func TestSpecContentWrapsToItsOwnHalf(t *testing.T) {
	long := strings.Repeat("word ", 60)
	m := makeSpecsModel(t, map[string]string{"alpha": long})
	listWidth, contentWidth := specsSplit(m.panelContentWidth())

	if m.docViewport.Width() != contentWidth {
		t.Errorf("viewport width is %d, want the content half %d (list %d)",
			m.docViewport.Width(), contentWidth, listWidth)
	}
	for _, row := range strings.Split(m.docViewport.View(), "\n") {
		if got := ansi.StringWidth(row); got > contentWidth {
			t.Errorf("row is %d cells, wider than the content half %d: %q",
				got, contentWidth, row)
		}
	}

	// A line that fits the whole panel but not the content half must still
	// wrap, which is what proves the narrower width is being used.
	fits := strings.Repeat("x", contentWidth+10)
	m2 := makeSpecsModel(t, map[string]string{"alpha": fits})
	if rows := strings.Split(m2.docViewport.View(), "\n"); strings.TrimSpace(rows[1]) == "" {
		t.Error("a line wider than the content half did not wrap")
	}
}

// --- 4.9 a project with no specs ---

func TestSpecsTabWithNoSpecs(t *testing.T) {
	m := makeSpecsModel(t, map[string]string{})

	if m.docActive() {
		t.Error("there is no spec, so nothing should own the vertical axis")
	}

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.specsFocus != specsFocusList {
		t.Errorf("focus moved to %d with no specs to focus on", m.specsFocus)
	}

	out := m.renderSpecsTab(m.width-2, 10)
	if !strings.Contains(out, "No specs found") {
		t.Errorf("expected the empty message, got:\n%s", out)
	}
}

func TestSpecWithoutAFile(t *testing.T) {
	m := makeSpecsModel(t, map[string]string{"alpha": ""})
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})

	if !strings.Contains(m.docViewport.View(), "No spec.md found") {
		t.Errorf("expected the missing-file message, got:\n%s", m.docViewport.View())
	}
}

func TestSpecsListShowsWhichHalfHasTheKeyboard(t *testing.T) {
	m := twoSpecs(t)

	lit := m.renderSpecsTab(m.width-2, 10)
	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	dimmed := m.renderSpecsTab(m.width-2, 10)

	if lit == dimmed {
		t.Error("the spec list looks identical whichever half holds the keyboard")
	}
}
